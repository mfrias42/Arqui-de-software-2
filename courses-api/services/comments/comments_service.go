package comments

import (
	"context"
	"fmt"
	"time"

	commentsDAO "courses-api/DAO/comments"
	coursesDAO "courses-api/DAO/courses"
	commentsDomain "courses-api/domain/comments"
)

type CommentsRepository interface {
	CreateComment(ctx context.Context, comment commentsDAO.Comment) (commentsDAO.Comment, error)
	GetCommentsByCourseID(ctx context.Context, courseID int64) ([]commentsDAO.Comment, error)
	DeleteCommentsByCourseID(ctx context.Context, courseID int64) error
}

type CoursesRepository interface {
	GetCourseByID(ctx context.Context, id int64) (coursesDAO.Course, error)
	UpdateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error)
}

type Service struct {
	commentsRepository CommentsRepository
	coursesRepository  CoursesRepository
}

func NewService(commentsRepo CommentsRepository, coursesRepo CoursesRepository) Service {
	return Service{
		commentsRepository: commentsRepo,
		coursesRepository:  coursesRepo,
	}
}

func (s Service) CreateComment(ctx context.Context, courseID int64, req commentsDomain.CreateCommentRequest) (commentsDomain.CommentResponse, error) {
	// Verificar si el curso existe
	_, err := s.coursesRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		return commentsDomain.CommentResponse{}, fmt.Errorf("el curso con ID %d no existe: %v", courseID, err)
	}

	comment := commentsDAO.Comment{
		CourseID:  courseID,
		UserID:    req.UserID,
		Content:   req.Content,
		Rating:    req.Rating,
		CreatedAt: time.Now().Unix(),
	}

	createdComment, err := s.commentsRepository.CreateComment(ctx, comment)
	if err != nil {
		return commentsDomain.CommentResponse{}, fmt.Errorf("error al crear el comentario: %v", err)
	}

	// Actualizar el rating del curso
	err = s.updateCourseRating(ctx, courseID)
	if err != nil {
		return commentsDomain.CommentResponse{}, fmt.Errorf("error al actualizar el rating del curso: %v", err)
	}

	return commentsDomain.CommentResponse{
		ID:        createdComment.ID,
		CourseID:  createdComment.CourseID,
		UserID:    createdComment.UserID,
		Content:   createdComment.Content,
		Rating:    createdComment.Rating,
		CreatedAt: createdComment.CreatedAt,
	}, nil
}

func (s Service) GetCommentsByCourseID(ctx context.Context, courseID int64) ([]commentsDomain.CommentResponse, error) {
	// Verificar si el curso existe
	_, err := s.coursesRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("el curso con ID %d no existe: %v", courseID, err)
	}

	commentsDB, err := s.commentsRepository.GetCommentsByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los comentarios: %v", err)
	}

	var commentsResponse []commentsDomain.CommentResponse
	for _, comment := range commentsDB {
		commentsResponse = append(commentsResponse, commentsDomain.CommentResponse{
			ID:        comment.ID,
			CourseID:  comment.CourseID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			Rating:    comment.Rating,
			CreatedAt: comment.CreatedAt,
		})
	}

	return commentsResponse, nil
}

func (s Service) updateCourseRating(ctx context.Context, courseID int64) error {
	comments, err := s.commentsRepository.GetCommentsByCourseID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get comments: %v", err)
	}

	var totalRating float64
	for _, c := range comments {
		totalRating += float64(c.Rating)
	}
	newAverageRating := totalRating / float64(len(comments))

	course, err := s.coursesRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("error al obtener el curso: %v", err)
	}

	course.Rating = newAverageRating

	_, err = s.coursesRepository.UpdateCourse(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to update course rating: %v", err)
	}

	return nil
}

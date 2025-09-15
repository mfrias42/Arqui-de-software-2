package files

import (
	"context"
	coursesDAO "courses-api/DAO/courses"
	filesDAO "courses-api/DAO/files"
	"courses-api/domain/files"
	"fmt"
)

// Interface del repositorio de archivos
type Repository interface {
	CreateFile(ctx context.Context, file filesDAO.File) (filesDAO.File, error)
	GetFilesByCourseID(ctx context.Context, courseID int64) ([]filesDAO.File, error)
}

// Interface del repositorio de cursos
type CourseRepository interface {
	GetCourseByID(ctx context.Context, id int64) (coursesDAO.Course, error)
}

// Servicio de archivos
type Service struct {
	repository       Repository
	courseRepository CourseRepository
}

// Constructor del servicio
func NewService(repo Repository, courseRepo CourseRepository) Service {
	return Service{
		repository:       repo,
		courseRepository: courseRepo,
	}
}

// Crear archivo
func (s Service) CreateFile(ctx context.Context, req files.CreateFileRequest) (files.FileResponse, error) {
	// Verificar si el curso existe
	_, err := s.courseRepository.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return files.FileResponse{}, fmt.Errorf("el curso no existe: %v", err)
	}

	file := filesDAO.File{
		Name:     req.Name,
		Content:  []byte(req.Content),
		UserID:   req.UserID,
		CourseID: req.CourseID,
	}

	createdFile, err := s.repository.CreateFile(ctx, file)
	if err != nil {
		return files.FileResponse{}, fmt.Errorf("error al crear el archivo: %v", err)
	}

	return files.FileResponse{
		ID:       createdFile.ID,
		Name:     createdFile.Name,
		Content:  createdFile.Content,
		UserID:   createdFile.UserID,
		CourseID: createdFile.CourseID,
	}, nil
}

// Obtener archivos por ID de curso
func (s Service) GetFilesByCourseID(ctx context.Context, courseID int64) ([]files.FileResponse, error) {
	filesData, err := s.repository.GetFilesByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %v", err)
	}

	var response []files.FileResponse
	for _, f := range filesData {
		response = append(response, files.FileResponse{
			ID:       f.ID,
			Name:     f.Name,
			Content:  f.Content,
			UserID:   f.UserID,
			CourseID: f.CourseID,
		})
	}
	return response, nil
}

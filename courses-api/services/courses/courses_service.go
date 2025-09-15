package courses

import (
	"context"
	coursesDAO "courses-api/DAO/courses"
	"courses-api/clients"
	"courses-api/domain/courses"
	"fmt"
	"sync"
)

// Repository interface para las operaciones de curso
type Repository interface {
	CreateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error)
	GetCourses(ctx context.Context) ([]coursesDAO.Course, error)
	GetCourseByID(ctx context.Context, id int64) (coursesDAO.Course, error)
	UpdateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error)
	DeleteCourse(ctx context.Context, id int64) error
}

// CommentsRepository interface para las operaciones de comentarios
type CommentsRepository interface {
	DeleteCommentsByCourseID(ctx context.Context, courseID int64) error
}

type Queue interface {
	Publish(courseNew courses.CursosNew) error
}

// FilesRepository interface para las operaciones de archivos
type FilesRepository interface {
	DeleteFilesByCourseID(ctx context.Context, courseID int64) error
}

// Service estructura para el servicio de cursos
type Service struct {
	repository         Repository
	commentsRepository CommentsRepository
	filesRepository    FilesRepository
	eventsQueue        Queue
	httpClient         *clients.HTTPClient
}

// NewService constructor para el servicio de cursos
func NewService(repository Repository, commentsRepository CommentsRepository, filesRepository FilesRepository, eventsQueue Queue, httpClient *clients.HTTPClient) Service {
	return Service{
		repository:         repository,
		commentsRepository: commentsRepository,
		filesRepository:    filesRepository,
		eventsQueue:        eventsQueue,
		httpClient:         httpClient,
	}
}

func (s Service) CreateCourse(ctx context.Context, req courses.CreateCourseRequest) (courses.CourseResponse, error) {
	course := coursesDAO.Course{
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Duration:     req.Duration,
		InstructorID: req.InstructorID,
		ImageBase64:  req.ImageBase64, // Guardar la imagen subida por el usuario
		Capacity:     req.Capacity,
		Rating:       0,                // Inicialmente, el rating es 0
		Available:    req.Capacity > 0, // Asignar disponibilidad
	}

	createdCourse, err := s.repository.CreateCourse(ctx, course)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("failed to create course: %v", err)
	}

	go func() {
		if err := s.eventsQueue.Publish(courses.CursosNew{
			Operation: "POST",
			ID:        createdCourse.ID,
		}); err != nil {
			fmt.Printf("Error al publicar nuevo curso: %v", err)
		}
	}()

	return courses.CourseResponse{
		ID:           createdCourse.ID,
		Name:         createdCourse.Name,
		Description:  createdCourse.Description,
		Category:     createdCourse.Category,
		Duration:     createdCourse.Duration,
		InstructorID: createdCourse.InstructorID,
		ImageBase64:  createdCourse.ImageBase64,
		Capacity:     createdCourse.Capacity,
		Rating:       createdCourse.Rating,
		Available:    createdCourse.Available,
	}, nil
}

func (s Service) GetCourses(ctx context.Context) ([]courses.CourseResponse, error) {
	coursesDAO, err := s.repository.GetCourses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	var coursesResponse []courses.CourseResponse
	for _, course := range coursesDAO {
		coursesResponse = append(coursesResponse, courses.CourseResponse{
			ID:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			Category:     course.Category,
			Duration:     course.Duration,
			InstructorID: course.InstructorID,
			ImageBase64:  course.ImageBase64,
			Capacity:     course.Capacity,
			Rating:       course.Rating,
			Available:    course.Available,
		})
	}

	return coursesResponse, nil
}

func (s Service) GetCourseByID(ctx context.Context, id int64) (courses.CourseResponse, error) {
	course, err := s.repository.GetCourseByID(ctx, id)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("failed to get course: %v", err)
	}

	return courses.CourseResponse{
		ID:           course.ID,
		Name:         course.Name,
		Description:  course.Description,
		Category:     course.Category,
		Duration:     course.Duration,
		InstructorID: course.InstructorID,
		ImageBase64:  course.ImageBase64,
		Capacity:     course.Capacity,
		Rating:       course.Rating,
		Available:    course.Available,
	}, nil
}

func (s Service) UpdateCourse(ctx context.Context, id int64, req courses.UpdateCourseRequest) (courses.CourseResponse, error) {
	var auxiliar int
	var borrarcurso bool
	borrarcurso = false
	course, err := s.repository.GetCourseByID(ctx, id)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("course not found: %v", err)
	}

	// Verificar si la capacidad ha cambiado

	// Actualizar los campos del curso
	if req.Name != "" {
		course.Name = req.Name
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Category != "" {
		course.Category = req.Category
	}
	if req.Duration != "" {
		course.Duration = req.Duration
	}
	if req.InstructorID != 0 {
		course.InstructorID = req.InstructorID
	}
	if req.Capacity != 0 {
		if req.Capacity < course.Capacity {
			// Verificar inscripciones actuales
			inscriptions, err := s.httpClient.GetInscriptionsByCourse(uint(id))
			if err != nil {
				return courses.CourseResponse{}, fmt.Errorf("error al verificar inscripciones: %v", err)
			}
			if len(inscriptions) > req.Capacity {
				return courses.CourseResponse{}, fmt.Errorf("no se puede reducir la capacidad a un número menor que las inscripciones actuales")
			}

		}
		if req.Capacity > course.Capacity {
			if !course.Available {
				req.Available = false
				course.Available = true
			}
		}

		auxiliar = course.Capacity
		course.Capacity = req.Capacity
	}
	if req.ImageBase64 != "" {
		course.ImageBase64 = req.ImageBase64
	}

	updatedCourse, err := s.repository.UpdateCourse(ctx, course)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("failed to update course: %v", err)
	}

	if req.Capacity != 0 {
		if req.Capacity < auxiliar {
			inscriptions, err := s.httpClient.GetInscriptionsByCourse(uint(id))
			if err != nil {
				return courses.CourseResponse{}, fmt.Errorf("error al verificar inscripciones: %v", err)
			}
			if len(inscriptions) == req.Capacity {
				if err := s.UpdateCourseAvailability(ctx, course.ID); err != nil {
					return courses.CourseResponse{}, fmt.Errorf("error al actualizar disponibilidad: %v", err)
				}
				borrarcurso = true
			}

		}
		if req.Capacity > auxiliar {
			if !req.Available {
				go func() {
					if err := s.eventsQueue.Publish(courses.CursosNew{
						Operation: "POST",
						ID:        course.ID,
					}); err != nil {
						fmt.Printf("Error al publicar nuevo curso: %v", err)
					}
				}()

			}

		}
	}
	// Actualizar la disponibilidad del curso

	if course.Available {
		if !borrarcurso {
			go func() {
				if err := s.eventsQueue.Publish(courses.CursosNew{
					Operation: "UPDATE",
					ID:        course.ID,
				}); err != nil {
					fmt.Printf("Error al publicar nuevo curso: %v", err)
				}
			}()
		}
	}
	return courses.CourseResponse{
		ID:           updatedCourse.ID,
		Name:         updatedCourse.Name,
		Description:  updatedCourse.Description,
		Category:     updatedCourse.Category,
		Duration:     updatedCourse.Duration,
		InstructorID: updatedCourse.InstructorID,
		ImageBase64:  updatedCourse.ImageBase64,
		Capacity:     updatedCourse.Capacity,
		Rating:       updatedCourse.Rating,
		Available:    updatedCourse.Available,
	}, nil
}

func (s Service) DeleteCourse(ctx context.Context, id int64) error {
	// Verificar si hay inscripciones para este curso
	inscriptions, err := s.httpClient.GetInscriptionsByCourse(uint(id))
	if err != nil {
		return fmt.Errorf("error al verificar inscripciones: %v", err)
	}

	if len(inscriptions) > 0 {
		return fmt.Errorf("no se puede eliminar el curso porque tiene %d inscripciones activas", len(inscriptions))
	}

	// Eliminar los comentarios asociados al curso
	err = s.commentsRepository.DeleteCommentsByCourseID(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar los comentarios del curso: %v", err)
	}

	// Eliminar los archivos asociados al curso
	err = s.filesRepository.DeleteFilesByCourseID(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar los archivos del curso: %v", err)
	}

	// Eliminar el curso
	err = s.repository.DeleteCourse(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar el curso: %v", err)
	}

	go func() {
		if err := s.eventsQueue.Publish(courses.CursosNew{
			Operation: "DELETE",
			ID:        id,
		}); err != nil {
			fmt.Printf("Error al publicar eliminación de curso: %v", err)
		}
	}()

	return nil
}

// Agregar este método para actualizar el rating del curso
func (s Service) UpdateCourseRating(ctx context.Context, courseID int64, newRating float64) error {
	course, err := s.repository.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %v", err)
	}

	course.Rating = newRating
	_, err = s.repository.UpdateCourse(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to update course rating: %v", err)
	}

	return nil
}

func (s Service) UpdateCourseAvailability(ctx context.Context, courseID int64) error {
	// Obtener el curso por ID
	course, err := s.repository.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("error al obtener el curso: %v", err)
	}

	// Log para mostrar el estado del curso
	fmt.Printf("Estado del curso antes de la actualización: %+v\n", course)

	// Obtener las inscripciones actuales para el curso
	inscriptions, err := s.httpClient.GetInscriptionsByCourse(uint(courseID))
	if err != nil {
		return fmt.Errorf("error al obtener inscripciones: %v", err)
	}

	// Log para mostrar la cantidad de inscripciones
	fmt.Printf("Cantidad de inscripciones para el curso ID %d: %d\n", courseID, len(inscriptions))

	// Verificar si la cantidad de inscripciones es igual o mayor que la capacidad
	if len(inscriptions) >= course.Capacity {
		course.Available = false
		fmt.Printf("El curso ID %d ya no está disponible. Inscripciones: %d, Capacidad: %d\n", courseID, len(inscriptions), course.Capacity)

		// Publicar un mensaje en RabbitMQ para indicar que el curso ya no está disponible
		if err := s.eventsQueue.Publish(courses.CursosNew{
			Operation: "DELETE",
			ID:        course.ID,
		}); err != nil {
			fmt.Printf("Error al publicar eliminación de curso en RabbitMQ: %v", err)
		}
	} else {
		course.Available = true
		fmt.Printf("El curso ID %d está disponible. Inscripciones: %d, Capacidad: %d\n", courseID, len(inscriptions), course.Capacity)
	}

	// Actualizar el curso en la base de datos
	_, err = s.repository.UpdateCourse(ctx, course)
	if err != nil {
		return fmt.Errorf("error al actualizar la disponibilidad del curso: %v", err)
	}

	return nil
}

func (s Service) CourseAvailability(ctx context.Context) ([]courses.CourseResponse, error) {
	coursesDAO, err := s.repository.GetCourses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	var coursesResponse []courses.CourseResponse
	for _, course := range coursesDAO {
		coursesResponse = append(coursesResponse, courses.CourseResponse{
			ID:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			Category:     course.Category,
			Duration:     course.Duration,
			InstructorID: course.InstructorID,
			ImageBase64:  course.ImageBase64,
			Capacity:     course.Capacity,
			Rating:       course.Rating,
			Available:    course.Available,
		})
	}

	coursesChannel := make(chan courses.CourseResponse, len(coursesResponse))
	var wg sync.WaitGroup

	for _, course := range coursesResponse {
		wg.Add(1)
		go func(course courses.CourseResponse) {
			defer wg.Done()
			if course.Available {
				coursesChannel <- course
			}
		}(course)
	}

	wg.Wait()
	close(coursesChannel)

	var results []courses.CourseResponse
	for course := range coursesChannel {
		results = append(results, course)
	}

	return results, nil
}

package search

import (
	"context"
	"fmt"
	"log"
	dao "search-api/dao/courses"                            // Alias para los tipos de DAO
	domain "search-api/domain/courses"                      // Alias para los tipos de dominio
	httpRepo "search-api/repositories/courses/courses_http" // Importar el paquete HTTP
	"strconv"
)

// Repository define las operaciones necesarias en el índice de SolR
type Repository interface {
	Index(ctx context.Context, course dao.Course) (string, error)
	Update(ctx context.Context, course dao.Course) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int, offset int) ([]dao.Course, error)
}

// Service representa el servicio de búsqueda
type Service struct {
	repository Repository
	httpClient httpRepo.HTTP // Cliente HTTP para interactuar con la API de Cursos
}

// NewService crea una nueva instancia del servicio de búsqueda
func NewService(repository Repository, httpClient httpRepo.HTTP) Service {
	return Service{
		repository: repository,
		httpClient: httpClient,
	}
}

// HandleCourseUpdate procesa las actualizaciones de cursos recibidas desde RabbitMQ
func (service Service) HandleCourseUpdate(courseNew domain.CourseUpdate) {
	ctx := context.Background()

	// Agregar log para ver el mensaje recibido
	log.Printf("Mensaje recibido para procesar: %+v", courseNew)

	// Convertir ID a string
	courseIDStr := strconv.FormatInt(courseNew.ID, 10) // Convertir ID a string

	switch courseNew.Operation {
	case "POST":

		// Llamar a GetCourseByID y almacenar el resultado en 'curso'
		courseUpdate, err := service.httpClient.GetCourseByID(ctx, courseIDStr) // Usar courseIDStr
		curso := dao.Course{
			ID:          courseUpdate.ID,
			Name:        courseUpdate.Name,
			Category:    courseUpdate.Category,
			Description: courseUpdate.Description,
		}
		if err != nil {
			log.Printf("Error al obtener el curso (ID: %s): %v", courseIDStr, err) // Cambiar ID a string
			return                                                                 // Salir de la función si hay un error
		}
		log.Printf("Curso obtenido, listo para procesar la operación: %d", curso.ID) // ID sigue en int64
		log.Printf("Procesando operación POST para el curso: %d", curso.ID)
		// Indexar el nuevo curso en SolR
		if _, err := service.repository.Index(ctx, curso); err != nil {
			log.Printf("Error al indexar el curso (%d): %v", curso.ID, err)
		} else {
			log.Printf("Curso indexado exitosamente: %d", curso.ID)
		}

	case "UPDATE":

		courseUpdate, err := service.httpClient.GetCourseByID(ctx, courseIDStr) // Usar courseIDStr
		curso := dao.Course{
			ID:          courseUpdate.ID,
			Name:        courseUpdate.Name,
			Category:    courseUpdate.Category,
			Description: courseUpdate.Description,
		}
		if err != nil {
			log.Printf("Error al obtener el curso (ID: %s): %v", courseIDStr, err) // Cambiar ID a string
			return                                                                 // Salir de la función si hay un error
		}
		log.Printf("Curso obtenido, listo para procesar la operación: %d", curso.ID) // ID sigue en int64
		log.Printf("Procesando operación UPDATE para el curso: %d", curso.ID)
		// Actualizar el curso existente en SolR
		if err := service.repository.Update(ctx, curso); err != nil {
			log.Printf("Error al actualizar el curso (%d): %v", courseUpdate.ID, err)
		} else {
			log.Printf("Curso actualizado exitosamente: %d", courseUpdate.ID)
		}

	case "DELETE":
		log.Printf("Procesando operación DELETE para el curso: %d", courseNew.ID)
		// Eliminar el curso del índice de SolR
		if err := service.repository.Delete(ctx, courseIDStr); err != nil { // Usar courseIDStr
			log.Printf("Error al eliminar el curso (%d): %v", courseNew.ID, err)
		} else {
			log.Printf("Curso eliminado exitosamente: %d", courseNew.ID)
		}

	default:
		log.Printf("Operación desconocida: %s", courseNew.Operation)
	}
}

// Search busca cursos en SolR según el término de búsqueda, límite y desplazamiento
func (service Service) Search(ctx context.Context, query string, limit int, offset int) ([]domain.CourseUpdate, error) {
	daoResults, err := service.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error en la búsqueda de cursos: %w", err)
	}

	// Convertir de dao.Course a domain.CourseUpdate
	var results []domain.CourseUpdate
	for _, daoCourse := range daoResults {
		results = append(results, domain.CourseUpdate{
			Operation:   "SEARCH", // El tipo de operación es 'SEARCH' para los resultados de la búsqueda
			ID:          daoCourse.ID,
			Name:        daoCourse.Name,
			Category:    daoCourse.Category,
			Description: daoCourse.Description,
		})
	}

	return results, nil
}

package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	daoCourses "search-api/dao/courses"       // Alias para evitar conflictos
	domainCourses "search-api/domain/courses" // Alias para evitar conflictos
)

type HTTPConfig struct {
	Host string
	Port string
}

type HTTP struct {
	baseURL func(courseID string) string
}

// NewHTTP crea una nueva conexión a la API de cursos
func NewHTTP(config HTTPConfig) HTTP {
	return HTTP{
		baseURL: func(courseID string) string {
			return fmt.Sprintf("http://%s:%s/courses/%s", config.Host, config.Port, courseID)
		},
	}
}

// GetCourseByID obtiene los detalles de un curso usando su ID
func (repository HTTP) GetCourseByID(ctx context.Context, id string) (domainCourses.CourseUpdate, error) {
	resp, err := http.Get(repository.baseURL(id))
	if err != nil {
		return domainCourses.CourseUpdate{}, fmt.Errorf("error fetching course (%s): %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domainCourses.CourseUpdate{}, fmt.Errorf("failed to fetch course (%s): received status code %d", id, resp.StatusCode)
	}

	// Lee el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domainCourses.CourseUpdate{}, fmt.Errorf("error reading response body for course (%s): %w", id, err)
	}

	// Deserializa los datos del curso en la estructura CourseUpdate
	var course domainCourses.CourseUpdate // Usando CourseUpdate del dominio
	if err := json.Unmarshal(body, &course); err != nil {
		return domainCourses.CourseUpdate{}, fmt.Errorf("error unmarshaling course data (%s): %w", id, err)
	}

	return course, nil
}

// GetCourses obtiene todos los cursos
func (repository HTTP) GetCoursesAvailability(ctx context.Context) ([]daoCourses.Course, error) {
	log.Println("Iniciando la obtención de todos los cursos Disponibles...") // Log de inicio

	url := repository.baseURL("availability")      // Solo llama a baseURL sin concatenar http://
	log.Printf("URL de la API de cursos: %s", url) // Log de la URL

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error al hacer la solicitud HTTP: %v", err) // Log del error
		return nil, fmt.Errorf("error fetching courses: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Código de estado de la respuesta: %d", resp.StatusCode) // Log del código de estado

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error al obtener cursos: código de estado recibido %d", resp.StatusCode) // Log del error
		return nil, fmt.Errorf("failed to fetch courses: received status code %d", resp.StatusCode)
	}

	var courses []daoCourses.Course // Usando Course del DAO
	if err := json.NewDecoder(resp.Body).Decode(&courses); err != nil {
		log.Printf("Error al deserializar los datos de los cursos: %v", err) // Log del error
		return nil, fmt.Errorf("error unmarshaling courses data: %w", err)
	}

	log.Printf("Cursos obtenidos: %+v", courses) // Log de los cursos obtenidos
	return courses, nil
}

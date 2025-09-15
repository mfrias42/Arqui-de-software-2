package search

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"search-api/domain/courses"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Service define la interfaz del servicio de búsqueda
type Service interface {
	Search(ctx context.Context, query string, offset int, limit int) ([]courses.CourseUpdate, error)
}

// Controller representa el controlador de búsqueda
type Controller struct {
	service Service
}

// NewController crea una nueva instancia del controlador de búsqueda
func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

// Search maneja las solicitudes GET en el endpoint /search
func (controller Controller) Search(c *gin.Context) {
	// Parsear el parámetro de búsqueda "query" de la URL
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'q' es obligatorio"})
		return
	}

	// Parsear el parámetro "offset" de la URL
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0 // Valor por defecto si no se proporciona
	}

	// Parsear el parámetro "limit" de la URL
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Valor por defecto si no se proporciona o es inválido
	}

	// Llamar al servicio de búsqueda
	results, err := controller.service.Search(c.Request.Context(), query, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error en la búsqueda: %v", err)})
		return
	}

	// Log para ver los resultados de la búsqueda
	log.Printf("Resultados de la búsqueda para la consulta '%s': %+v", query, results)

	// Enviar los resultados como respuesta JSON
	c.JSON(http.StatusOK, results)
}

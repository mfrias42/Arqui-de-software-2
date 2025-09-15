package courses

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	daoCourses "search-api/dao/courses"
	"strconv" // Asegúrate de que esta línea esté presente

	"github.com/stevenferrer/solr-go"
)

type SolrConfig struct {
	Host       string // Solr host
	Port       string // Solr port
	Collection string // Solr collection name
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) Solr {
	baseURL := fmt.Sprintf("http://%s:%s", config.Host, config.Port)
	client := solr.NewJSONClient(baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
	}
}

// Index adds a new course document to the Solr collection
func (searchEngine Solr) Index(ctx context.Context, course daoCourses.Course) (string, error) {
	// Prepare the document for SolR
	doc := map[string]interface{}{
		"id":          course.ID,
		"name":        course.Name,
		"category":    course.Category,
		"description": course.Description,
	}

	// Log the course ID being indexed
	log.Printf("Indexando curso con ID: %d", course.ID)
	log.Printf("Documento a indexar: %+v", doc)

	// Prepare the index request
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Index the document in SolR
	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling course document: %w", err)
	}

	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		log.Printf("Error al indexar en Solr: %v", err)
		return "", fmt.Errorf("error indexing course: %w", err)
	}
	if resp.Error != nil {
		log.Printf("Respuesta de Solr después de indexar: %+v", resp)
		return "", fmt.Errorf("failed to index course: %v", resp.Error)
	}

	log.Printf("Respuesta de Solr después de indexar: %+v", resp) // Log de respuesta

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error committing changes to SolR: %w", err)
	}

	return fmt.Sprintf("%d", course.ID), nil // Convert CourseID to string
}

// Update modifies an existing course document in the Solr collection
func (searchEngine Solr) Update(ctx context.Context, course daoCourses.Course) error {
	doc := map[string]interface{}{
		"id":          course.ID,
		"name":        course.Name,
		"category":    course.Category,
		"description": course.Description,
	}

	log.Printf("Actualizando curso con ID: %d", course.ID)
	log.Printf("Documento a actualizar: %+v", doc)

	updateRequest := map[string]interface{}{
		"add": []interface{}{doc},
	}

	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error marshaling course document: %w", err)
	}

	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating course: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to update course: %v", resp.Error)
	}

	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to SolR: %w", err)
	}

	return nil
}

// Delete removes a course document from the Solr collection
func (searchEngine Solr) Delete(ctx context.Context, id string) error {
	docToDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	body, err := json.Marshal(docToDelete)
	if err != nil {
		return fmt.Errorf("error marshaling course document: %w", err)
	}

	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error deleting course: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to delete course: %v", resp.Error)
	}

	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to SolR: %w", err)
	}

	return nil
}

// Search searches for courses in the Solr collection based on a query
func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]daoCourses.Course, error) {
	if query == "" {
		return nil, fmt.Errorf("la consulta no puede estar vacía")
	}

	solrQuery := fmt.Sprintf("q=(name:*%s* OR description:*%s*)&rows=%d&start=%d&wt=json", query, query, limit, offset)

	log.Printf("Consulta a Solr: %s", solrQuery)

	resp, err := searchEngine.Client.Query(ctx, searchEngine.Collection, solr.NewQuery(solrQuery))
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la consulta de búsqueda: %w", err)
	}
	log.Printf("Respuesta de Solr: %+v", resp) // Log de la respuesta de Solr
	if resp.Error != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta de búsqueda: %v, consulta: %s", resp.Error, solrQuery)
	}

	var coursesList []daoCourses.Course
	for _, doc := range resp.Response.Documents {
		course := daoCourses.Course{
			ID:          getIntField(doc, "id"),
			Name:        getStringField(doc, "name"),
			Category:    getStringField(doc, "category"),
			Description: getStringField(doc, "description"),
		}
		coursesList = append(coursesList, course)
	}

	return coursesList, nil
}

// Helper function to safely get string fields from the document
func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get int64 fields from the document
func getIntField(doc map[string]interface{}, field string) int64 {
	log.Printf("Document fields: %+v", doc)    // Log para ver todos los campos del documento
	log.Printf("Buscando el campo: %s", field) // Log para ver qué campo se está buscando

	if val, ok := doc[field].(float64); ok {
		log.Printf("Campo '%s' encontrado como float64 con valor: %f", field, val)
		return int64(val)
	}
	if val, ok := doc[field].(int64); ok {
		log.Printf("Campo '%s' encontrado como int64 con valor: %d", field, val)
		return val
	}

	if val, ok := doc[field].(string); ok {
		log.Printf("Campo 'id' encontrado como string con valor: %s", val)
		if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			return intVal // Convertir a int64
		}
		log.Printf("Error al convertir el campo 'id' de string a int64: ")
	}
	log.Printf("Campo '%s' no encontrado o no es un float64 ni int64", field) // Ahora incluye el nombre del campo
	return 0
}

// IndexAllCourses indexa todos los cursos en Solr
func (searchEngine Solr) IndexAllCourses(ctx context.Context, courses []daoCourses.Course) error {
	for _, course := range courses {

		courseUpdate := daoCourses.Course{
			ID:          course.ID, // Asignar el ID convertido
			Name:        course.Name,
			Category:    course.Category,
			Description: course.Description,
		}
		if _, err := searchEngine.Index(ctx, courseUpdate); err != nil {
			return fmt.Errorf("error indexing course %s: %w", course.Name, err)
		}
	}
	return nil
}

package courses

// CourseUpdate representa una actualización de curso enviada a través de RabbitMQ
type CourseUpdate struct {
	Operation   string `json:"operation"`   // Tipo de operación: "CREATE", "UPDATE", "DELETE"
	ID          int64  `json:"id"`          // Identificador único del curso
	Name        string `json:"name"`        // Nombre del curso (para "CREATE" o "UPDATE")
	Category    string `json:"category"`    // Categoría del curso (para "CREATE" o "UPDATE")
	Description string `json:"description"` // Descripción del curso (para "CREATE" o "UPDATE")
	// Añadir más campos si es necesario
}

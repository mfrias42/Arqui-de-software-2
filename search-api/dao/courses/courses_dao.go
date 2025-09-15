package courses

// Course representa la estructura de un curso en la base de datos
type Course struct {
	ID          int64  `json:"id"`          // Cambiar a int64
	Name        string `json:"name"`        // Nombre del curso
	Category    string `json:"category"`    // Categoría del curso
	Description string `json:"description"` // Descripción del curso
}

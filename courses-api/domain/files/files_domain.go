package files

// CreateFileRequest representa la solicitud para cargar un archivo
type CreateFileRequest struct {
	Name     string `json:"name" binding:"required"`
	Content  string `json:"content" binding:"required"` // Base64 string
	UserID   int64  `json:"userId" binding:"required"`
	CourseID int64  `json:"-"` // Este campo se llenará con el valor de la URL
}

// FileResponse representa la respuesta al obtener un archivo
type FileResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Content  []byte `json:"content"`
	UserID   int64  `json:"userId"`
	CourseID int64  `json:"courseId"`
}

package courses

type CreateCourseRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Category     string `json:"category" binding:"required"`
	Duration     string `json:"duration" binding:"required"`
	InstructorID int64  `json:"instructor_id" binding:"required"`
	Capacity     int    `json:"capacity" binding:"required"`
	Available    bool   `json:"available"`
	ImageBase64  string `json:"imageBase64,omitempty"`
}

type UpdateCourseRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Duration     string  `json:"duration"`
	InstructorID int64   `json:"instructor_id"`
	Capacity     int     `json:"capacity"`
	Available    bool    `json:"available"`
	Rating       float64 `json:"rating"`
	ImageBase64  string  `json:"imageBase64,omitempty"`
}

type CourseResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Duration     string  `json:"duration"`
	InstructorID int64   `json:"instructor_id"`
	Capacity     int     `json:"capacity"`
	Available    bool    `json:"available"`
	Rating       float64 `json:"rating"`
	ImageBase64  string  `json:"imageBase64,omitempty"`
}
type CursosNew struct {
	Operation string `json:"operation"`
	ID        int64  `json:"id"`
}

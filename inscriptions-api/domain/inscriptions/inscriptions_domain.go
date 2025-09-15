package domain

type Inscription struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id"`
}

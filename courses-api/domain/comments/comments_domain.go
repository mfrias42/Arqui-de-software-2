package comments

type CreateCommentRequest struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
}

type CommentResponse struct {
	ID        int64  `json:"id"`
	CourseID  int64  `json:"course_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	Rating    int    `json:"rating"`
	CreatedAt int64  `json:"created_at"`
}

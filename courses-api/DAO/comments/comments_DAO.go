package comments

type Comment struct {
	ID        int64  `bson:"id"`
	CourseID  int64  `bson:"course_id"`
	UserID    int64  `bson:"user_id"`
	Content   string `bson:"content"`
	Rating    int    `bson:"rating"`
	CreatedAt int64  `bson:"created_at"`
}

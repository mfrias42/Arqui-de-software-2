package files

type File struct {
	ID       int64  `bson:"id"`
	Name     string `bson:"name"`
	Content  []byte `bson:"content"`
	UserID   int64  `bson:"user_id"`
	CourseID int64  `bson:"course_id"`
}

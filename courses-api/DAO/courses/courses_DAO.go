package courses

type Course struct {
	ID           int64   `bson:"id"`
	Name         string  `bson:"name"`
	Description  string  `bson:"description"`
	Category     string  `bson:"category"`
	Duration     string  `bson:"duration"`
	InstructorID int64   `bson:"instructor_id"`
	ImageBase64  string  `bson:"image_base64"`
	Capacity     int     `bson:"capacity"`
	Available    bool    `bson:"available"`
	Rating       float64 `bson:"rating"`
}

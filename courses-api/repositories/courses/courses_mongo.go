package repositories

import (
	"context"
	coursesDAO "courses-api/DAO/courses"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuración para MongoDB
type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

// Estructura del repositorio Mongo
type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

// Constante para la conexión
const (
	connectionURI = "mongodb://%s:%s"
)

// Nueva instancia de Mongo
func NewMongo(config MongoConfig) Mongo {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	ctx := context.Background()
	uri := fmt.Sprintf(connectionURI, config.Host, config.Port)
	cfg := options.Client().ApplyURI(uri).SetAuth(credentials)

	client, err := mongo.Connect(ctx, cfg)
	if err != nil {
		log.Panicf("error connecting to mongo DB: %v", err)
	}

	return Mongo{
		client:     client,
		database:   config.Database,
		collection: config.Collection,
	}
}

// Configuración del mutex y contador en memoria
var (
	counter   int64
	counterMu sync.Mutex
)

// Inicializa el contador en función del último ID en la colección
func InitializeCounter(mongoClient *mongo.Client, dbName, collectionName string) {
	time.Sleep(20 * time.Second)
	collection := mongoClient.Database(dbName).Collection(collectionName)
	var lastCourse coursesDAO.Course

	// Buscar el curso con el ID más alto
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&lastCourse)
	if err != nil {
		log.Printf("No se encontró un curso existente: %v. Iniciando el contador en 0.", err)
		counter = 0
	} else {
		counter = lastCourse.ID
		log.Printf("Contador inicializado en: %d", counter)
	}
}

// Obtener el próximo ID de manera segura
func getNextID() int64 {
	counterMu.Lock()
	defer counterMu.Unlock()
	counter++
	return counter
}

// Crear curso con el ID generado en memoria
func (m Mongo) CreateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error) {
	course.ID = getNextID()
	course.Rating = 0 // Inicializar el rating en 0

	collection := m.client.Database(m.database).Collection(m.collection)
	_, err := collection.InsertOne(ctx, course)
	if err != nil {
		return coursesDAO.Course{}, fmt.Errorf("failed to insert course: %v", err)
	}
	return course, nil
}

func (m Mongo) GetCourses(ctx context.Context) ([]coursesDAO.Course, error) {
	var courses []coursesDAO.Course
	collection := m.client.Database(m.database).Collection(m.collection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find courses: %v", err)
	}
	if err := cursor.All(ctx, &courses); err != nil {
		return nil, fmt.Errorf("failed to decode courses: %v", err)
	}
	return courses, nil
}

func (m Mongo) GetCourseByID(ctx context.Context, id int64) (coursesDAO.Course, error) {
	var course coursesDAO.Course
	collection := m.client.Database(m.database).Collection(m.collection)
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&course)
	if err != nil {
		return coursesDAO.Course{}, fmt.Errorf("failed to find course: %v", err)
	}
	return course, nil
}

func (m Mongo) UpdateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error) {
	collection := m.client.Database(m.database).Collection(m.collection)
	filter := bson.M{"id": course.ID}
	update := bson.M{"$set": course}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return coursesDAO.Course{}, fmt.Errorf("failed to update course: %v", err)
	}
	return course, nil
}

func (m Mongo) DeleteCourse(ctx context.Context, id int64) error {
	collection := m.client.Database(m.database).Collection(m.collection)
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete course: %v", err)
	}
	return nil
}

func (m Mongo) UpdateCourseRating(ctx context.Context, courseID int64, newRating float64) error {
	collection := m.client.Database(m.database).Collection(m.collection)
	filter := bson.M{"id": courseID}
	update := bson.M{"$set": bson.M{"rating": newRating}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update course rating: %v", err)
	}
	return nil
}

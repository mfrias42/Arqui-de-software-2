package files

import (
	"context"
	"fmt"
	"log"
	"sync"

	filesDAO "courses-api/DAO/files"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuración del mutex y contador en memoria
var (
	fileCounter   int64
	fileCounterMu sync.Mutex
)

// Repositorio MongoDB para archivos
type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

// Constructor del repositorio Mongo
func NewMongo(client *mongo.Client, db, collection string) Mongo {
	return Mongo{
		client:     client,
		database:   db,
		collection: collection,
	}
}

// Inicializa el contador en función del último ID en la colección
func InitializeFileCounter(mongoClient *mongo.Client, dbName, collectionName string) {
	collection := mongoClient.Database(dbName).Collection(collectionName)
	var lastFile filesDAO.File

	// Buscar el archivo con el ID más alto
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&lastFile)
	if err != nil {
		log.Printf("No se encontró un archivo existente: %v. Iniciando el contador en 0.", err)
		fileCounter = 0
	} else {
		fileCounter = lastFile.ID
		log.Printf("Contador de archivos inicializado en: %d", fileCounter)
	}
}

// Obtener el próximo ID de manera segura
func getNextFileID() int64 {
	fileCounterMu.Lock()
	defer fileCounterMu.Unlock()
	fileCounter++
	return fileCounter
}

// Crear archivo
func (m Mongo) CreateFile(ctx context.Context, file filesDAO.File) (filesDAO.File, error) {
	file.ID = getNextFileID()

	_, err := m.client.Database(m.database).Collection(m.collection).InsertOne(ctx, file)
	if err != nil {
		return filesDAO.File{}, fmt.Errorf("failed to insert file: %v", err)
	}
	return file, nil
}

// Obtener archivos por ID de curso
func (m Mongo) GetFilesByCourseID(ctx context.Context, courseID int64) ([]filesDAO.File, error) {
	var filesData []filesDAO.File
	cursor, err := m.client.Database(m.database).Collection(m.collection).Find(ctx, bson.M{"course_id": courseID})
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %v", err)
	}
	if err := cursor.All(ctx, &filesData); err != nil {
		return nil, fmt.Errorf("failed to decode files: %v", err)
	}
	return filesData, nil
}

func (m Mongo) DeleteFilesByCourseID(ctx context.Context, courseID int64) error {
	_, err := m.client.Database(m.database).Collection(m.collection).DeleteMany(ctx, bson.M{"course_id": courseID})
	if err != nil {
		return fmt.Errorf("failed to delete files for course: %v", err)
	}
	return nil
}

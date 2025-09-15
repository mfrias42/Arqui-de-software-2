package repositories

import (
	"context"
	"fmt"
	"log"
	"sync"

	commentsDAO "courses-api/DAO/comments"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentsMongo struct {
	client     *mongo.Client
	database   string
	collection string
}

var (
	commentCounter   int64
	commentCounterMu sync.Mutex
)

func NewCommentsMongo(client *mongo.Client, database, collection string) *CommentsMongo {
	return &CommentsMongo{
		client:     client,
		database:   database,
		collection: collection,
	}
}

// InitializeCommentCounter inicializa el contador de comentarios
func InitializeCommentCounter(client *mongo.Client, dbName, collectionName string) {
	collection := client.Database(dbName).Collection(collectionName)
	var lastComment commentsDAO.Comment

	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&lastComment)
	if err != nil {
		log.Printf("No se encontró un comentario existente: %v. Iniciando el contador en 0.", err)
		commentCounter = 0
	} else {
		commentCounter = lastComment.ID
		log.Printf("Contador de comentarios inicializado en: %d", commentCounter)
	}
}

// getNextCommentID obtiene el siguiente ID de comentario de manera segura
func getNextCommentID() int64 {
	commentCounterMu.Lock()
	defer commentCounterMu.Unlock()
	commentCounter++
	return commentCounter
}

func (m *CommentsMongo) CreateComment(ctx context.Context, comment commentsDAO.Comment) (commentsDAO.Comment, error) {
	comment.ID = getNextCommentID()

	collection := m.client.Database(m.database).Collection(m.collection)
	_, err := collection.InsertOne(ctx, comment)
	if err != nil {
		return commentsDAO.Comment{}, fmt.Errorf("failed to insert comment: %v", err)
	}
	return comment, nil
}

func (m *CommentsMongo) GetCommentsByCourseID(ctx context.Context, courseID int64) ([]commentsDAO.Comment, error) {
	collection := m.client.Database(m.database).Collection(m.collection)

	filter := bson.M{"course_id": courseID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find comments: %v", err)
	}
	defer cursor.Close(ctx)

	var comments []commentsDAO.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %v", err)
	}

	return comments, nil
}

func (m *CommentsMongo) DeleteCommentsByCourseID(ctx context.Context, courseID int64) error {
	collection := m.client.Database(m.database).Collection(m.collection)

	filter := bson.M{"course_id": courseID}
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete comments for course: %v", err)
	}

	return nil
}

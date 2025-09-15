package main

import (
	"context"
	"log"
	"os"
	"time"

	"courses-api/clients"
	"courses-api/clients/rabbit"
	commentsController "courses-api/controllers/comments"
	coursesController "courses-api/controllers/courses"
	filesController "courses-api/controllers/files"
	commentsRepositories "courses-api/repositories/comments"
	coursesRepositories "courses-api/repositories/courses"
	filesRepositories "courses-api/repositories/files"
	coursesRouter "courses-api/router/courses"
	commentsServices "courses-api/services/comments"
	coursesServices "courses-api/services/courses"
	filesServices "courses-api/services/files"

	"github.com/gin-contrib/cors"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	maxRetries = 5
	baseDelay  = 5 * time.Second
)

func connectWithRetry(connect func() error) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = connect()
		if err == nil {
			return nil
		}
		delay := time.Duration(attempt) * baseDelay
		log.Printf("Attempt %d failed: %v. Retrying in %v...", attempt, err, delay)
		time.Sleep(delay)
	}
	return err
}

func main() {
	// Configuración del cliente MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Crear cliente de MongoDB con reintento
	var client *mongo.Client
	err := connectWithRetry(func() error {
		var err error
		client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return client.Connect(ctx)
	})
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB after retries: %v", err)
	}

	// Inicializar contadores
	coursesRepositories.InitializeCounter(client, "courses-api", "courses")
	commentsRepositories.InitializeCommentCounter(client, "courses-api", "comments")
	filesRepositories.InitializeFileCounter(client, "courses-api", "files")

	// Configurar RabbitMQ
	rabbitURI := os.Getenv("RABBITMQ_URI")
	if rabbitURI == "" {
		rabbitURI = "amqp://guest:guest@localhost:5672/"
	}
	rabbitConfig := rabbit.RabbitConfig{
		URI:       rabbitURI,
		QueueName: "courses_queue",
	}

	// Conectar a RabbitMQ con reintento
	var rabbitConn *amqp.Connection
	err = connectWithRetry(func() error {
		var err error
		rabbitConn, err = amqp.Dial(rabbitConfig.URI)
		return err
	})
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after retries: %v", err)
	}
	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer rabbitChannel.Close()

	queue, err := rabbitChannel.QueueDeclare(
		rabbitConfig.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	rabbitQueue := rabbit.Rabbit{
		Connection: rabbitConn,
		Channel:    rabbitChannel,
		Queue:      queue,
	}

	// Crear instancias del repositorio
	courseRepo := coursesRepositories.NewMongo(coursesRepositories.MongoConfig{
		Host:       "mongodb",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "courses-api",
		Collection: "courses",
	})
	commentRepo := commentsRepositories.NewCommentsMongo(client, "courses-api", "comments")
	fileRepo := filesRepositories.NewMongo(client, "courses-api", "files")

	// Crear el cliente HTTP para la API de inscripciones
	inscriptionsAPIURL := os.Getenv("INSCRIPTIONS_API_URL")
	if inscriptionsAPIURL == "" {
		inscriptionsAPIURL = "http://nginx:8085"
	}
	httpClient := clients.NewHTTPClient(inscriptionsAPIURL)

	// Crear el servicio de cursos
	courseService := coursesServices.NewService(
		courseRepo,
		commentRepo,
		fileRepo,
		&rabbitQueue,
		httpClient,
	)

	// Crear el controlador de cursos
	courseController := coursesController.NewController(courseService)

	// Crear instancias para comentarios
	commentService := commentsServices.NewService(commentRepo, courseRepo)
	commentController := commentsController.NewController(commentService)

	// Crear instancias para archivos
	fileService := filesServices.NewService(fileRepo, courseRepo)
	fileController := filesController.NewController(fileService)

	// Leer la clave JWT desde la variable de entorno
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ThisIsAnExampleJWTKey!"
	}

	// Configurar las rutas
	router := coursesRouter.SetupRouter(courseController, commentController, fileController, jwtSecret)

	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"}, // Cambia esto si tu frontend está en otro puerto
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Leer el puerto desde las variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("Usando puerto por defecto:", port)
	}

	// Iniciar el servidor
	if err := router.Run(":" + port); err != nil {
		log.Fatal("No se pudo iniciar el servidor:", err)
	}
}

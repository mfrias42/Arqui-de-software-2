package main

import (
	"context"
	"log"
	"search-api/clients/queues"
	searchController "search-api/controllers/search"
	httpRepo "search-api/repositories/courses/courses_http"
	solrRepo "search-api/repositories/courses/courses_solr"
	searchService "search-api/services/search"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Esperando 2 minutos para que los servicios dependientes inicien...")
	time.Sleep(2 * time.Minute)
	log.Println("Iniciando la aplicación...")

	// Configuración de SolR
	solrClient := solrRepo.NewSolr(solrRepo.SolrConfig{
		Host:       "solr",    // SolR host
		Port:       "8983",    // SolR port
		Collection: "courses", // Nombre de la colección en SolR
	})

	// Configuración del cliente HTTP para la API de Cursos
	coursesClient := httpRepo.NewHTTP(httpRepo.HTTPConfig{
		Host: "courses-api",
		Port: "8080",
	})

	// Obtener todos los cursos
	coursesList, err := coursesClient.GetCoursesAvailability(context.Background())
	if err != nil {
		log.Fatalf("Error al obtener cursos: %v", err)
	}

	// Indexar todos los cursos en Solr
	if err := solrClient.IndexAllCourses(context.Background(), coursesList); err != nil {
		log.Fatalf("Error al indexar cursos en Solr: %v", err)
	}

	// Inicialización del servicio de búsqueda
	// Configuración de Solr y cliente HTTP
	searchSvc := searchService.NewService(solrClient, coursesClient)

	// Inicialización del controlador de búsqueda
	searchCtrl := searchController.NewController(searchSvc)

	// Lanzar el consumidor de RabbitMQ
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "courses_queue",
	})

	if err := eventsQueue.StartConsumer(searchSvc.HandleCourseUpdate); err != nil {
		log.Fatalf("Error al ejecutar el consumidor: %v", err)
	}

	// Configuración del router con Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type"},
	}))

	router.GET("/search", searchCtrl.Search)

	// Ejecutar la API en el puerto 8082
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Error al ejecutar la aplicación: %v", err)
	}
}

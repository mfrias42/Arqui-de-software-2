package main

import (
	"log"
	"os"
	"time"
	healthControllers "users-api/controllers/health"
	controllers "users-api/controllers/users"
	"users-api/internal/tokenizers"
	repositories "users-api/repositories/users"
	helathServices "users-api/services/health"
	services "users-api/services/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	log.Println("Esperando 2 minutos para que los servicios dependientes inicien...")
	time.Sleep(2 * time.Minute)
	log.Println("Iniciando la aplicación...")

	// MySQL
	mySQLRepo := repositories.NewMySQL(
		repositories.MySQLConfig{
			Host:     "mysql",
			Port:     "3306",
			Database: "users",
			Username: "root",
			Password: "root",
		},
	)

	// Cache
	cacheRepo := repositories.NewCache(repositories.CacheConfig{
		TTL: 1 * time.Minute,
	})

	// Memcached
	memcachedRepo := repositories.NewMemcached(repositories.MemcachedConfig{
		Host: "memcached",
		Port: "11211",
	})

	// Tokenizer
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ThisIsAnExampleJWTKey!"
	}
	jwtTokenizer := tokenizers.NewTokenizer(
		tokenizers.JWTConfig{
			Key:      jwtSecret,
			Duration: 1 * time.Hour,
		},
	)

	// Services
	service := services.NewService(mySQLRepo, cacheRepo, memcachedRepo, jwtTokenizer)

	// Handlers
	controller := controllers.NewController(service)
	helathService, err := helathServices.NewService() // Añadimos el manejo del error
	if err != nil {
		log.Fatalf("Error creating health service: %v", err)
	}
	log.Println("Servicio de health creado exitosamente")
	healthController := healthControllers.NewController(helathService)

	// Create router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.Default())

	// URL mappings
	router.GET("/users", controller.GetAll)
	router.GET("/users/:id", controller.GetByID)
	router.POST("/users", controller.Create)
	router.PUT("/users/:id", controller.Update)
	router.POST("/login", controller.Login)
	router.GET("/health", healthController.HealthCheck)

	// Run application
	if err := router.Run(":8083"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}

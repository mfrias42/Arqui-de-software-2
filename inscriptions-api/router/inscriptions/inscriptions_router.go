package router

import (
	controller "inscriptions-api/controllers/inscriptions"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// MapRoutes mapea las rutas del controlador de inscripciones.
func MapRoutes(r *gin.Engine, ctrl *controller.Controller) {
	// Configuración de CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8085"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas
	r.POST("/inscriptions", ctrl.CreateInscription)
	r.GET("/inscriptions", ctrl.GetInscriptions)
	r.GET("/users/:userID/inscriptions", ctrl.GetInscriptionsByUser)
	r.GET("/courses/:courseID/inscriptions", ctrl.GetInscriptionsByCourse)
}

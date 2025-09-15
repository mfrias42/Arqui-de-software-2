package health

import (
	"net/http"
	domain "users-api/domain/health"

	"log"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CheckServices() domain.HealthResponse
}

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) HealthCheck(ctx *gin.Context) {
	log.Println("Iniciando health check...")
	result := c.service.CheckServices()
	log.Printf("Health check completado. Encontrados %d servicios", len(result.Services))
	ctx.JSON(http.StatusOK, result)
}

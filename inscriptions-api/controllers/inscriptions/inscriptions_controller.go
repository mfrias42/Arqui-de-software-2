package controller

import (
	"context"
	"fmt"
	domain "inscriptions-api/domain/inscriptions"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateInscription(ctx context.Context, userID, courseID uint) (*domain.Inscription, error)
	GetInscriptions(ctx context.Context) ([]domain.Inscription, error)
	GetInscriptionsByUser(ctx context.Context, userID uint) ([]domain.Inscription, error)
	GetInscriptionsByCourse(ctx context.Context, courseID uint) ([]domain.Inscription, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) CreateInscription(c *gin.Context) {
	var req struct {
		UserID   uint `json:"user_id" binding:"required"`
		CourseID uint `json:"course_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid format: %s", err.Error())})
		return
	}

	inscription, err := ctrl.service.CreateInscription(c.Request.Context(), req.UserID, req.CourseID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user does not exist" || err.Error() == "course does not exist" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inscription)
}

func (ctrl *Controller) GetInscriptions(c *gin.Context) {
	inscriptions, err := ctrl.service.GetInscriptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, inscriptions)
}

func (ctrl *Controller) GetInscriptionsByUser(c *gin.Context) {
	userIDParam := strings.TrimSpace(c.Param("userID"))
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid user ID: %s", userIDParam)})
		return
	}

	inscriptions, err := ctrl.service.GetInscriptionsByUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inscriptions)
}

func (ctrl *Controller) GetInscriptionsByCourse(c *gin.Context) {
	courseIDParam := c.Param("courseID")
	courseID, err := strconv.ParseUint(courseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid course ID: %s", courseIDParam)})
		return
	}

	inscriptions, err := ctrl.service.GetInscriptionsByCourse(c.Request.Context(), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inscriptions)
}

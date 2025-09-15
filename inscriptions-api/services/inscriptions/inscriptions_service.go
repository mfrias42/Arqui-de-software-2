package service

import (
	"context"
	"errors"
	"fmt"
	dao "inscriptions-api/DAOs/inscriptions"
	"inscriptions-api/clients"
	domain "inscriptions-api/domain/inscriptions"
)

type Repository interface {
	CreateInscription(ctx context.Context, userID, courseID uint) (*dao.InscriptionModel, error)
	GetInscriptions(ctx context.Context) ([]dao.InscriptionModel, error)
	GetInscriptionsByUser(ctx context.Context, userID uint) ([]dao.InscriptionModel, error)
	GetInscriptionsByCourse(ctx context.Context, courseID uint) ([]dao.InscriptionModel, error)
}

type Service struct {
	repository Repository
	httpClient *clients.HTTPClient
}

func NewService(repository Repository, httpClient *clients.HTTPClient) *Service {
	return &Service{repository: repository, httpClient: httpClient}
}

func (s *Service) CreateInscription(ctx context.Context, userID, courseID uint) (*domain.Inscription, error) {
	if err := s.httpClient.CheckUserExists(userID); err != nil {
		return nil, fmt.Errorf("failed to verify user: %v", err)
	}

	course, err := s.httpClient.GetCourseDetails(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify course: %v", err)
	}

	if !course.Available {
		return nil, errors.New("course is not available for enrollment")
	}

	inscriptionModel, err := s.repository.CreateInscription(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to create inscription: %v", err)
	}

	if err := s.httpClient.UpdateCourseAvailability(int64(inscriptionModel.CourseID)); err != nil {
		return nil, fmt.Errorf("failed to update course availability: %v", err)
	}

	return &domain.Inscription{
		ID:       inscriptionModel.ID,
		UserID:   inscriptionModel.UserID,
		CourseID: inscriptionModel.CourseID,
	}, nil
}

func (s *Service) GetInscriptions(ctx context.Context) ([]domain.Inscription, error) {
	models, err := s.repository.GetInscriptions(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapModelsToDomain(models), nil
}

func (s *Service) GetInscriptionsByUser(ctx context.Context, userID uint) ([]domain.Inscription, error) {
	models, err := s.repository.GetInscriptionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mapModelsToDomain(models), nil
}

func (s *Service) GetInscriptionsByCourse(ctx context.Context, courseID uint) ([]domain.Inscription, error) {
	if err := s.httpClient.CheckCourseExists(courseID); err != nil {
		return nil, fmt.Errorf("failed to verify course: %v", err)
	}

	models, err := s.repository.GetInscriptionsByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	return s.mapModelsToDomain(models), nil
}

func (s *Service) mapModelsToDomain(models []dao.InscriptionModel) []domain.Inscription {
	inscriptions := make([]domain.Inscription, len(models))
	for i, model := range models {
		inscriptions[i] = domain.Inscription{
			ID:       model.ID,
			UserID:   model.UserID,
			CourseID: model.CourseID,
		}
	}
	return inscriptions
}

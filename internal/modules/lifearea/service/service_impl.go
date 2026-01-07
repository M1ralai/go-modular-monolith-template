package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/repository"
)

type lifeAreaService struct {
	repo   repository.LifeAreaRepository
	logger *logger.ZapLogger
}

func NewLifeAreaService(repo repository.LifeAreaRepository, logger *logger.ZapLogger) LifeAreaService {
	return &lifeAreaService{
		repo:   repo,
		logger: logger,
	}
}

func (s *lifeAreaService) Create(ctx context.Context, req *dto.CreateLifeAreaRequest, userID int) (*dto.LifeAreaResponse, error) {
	s.logger.Info("Creating life area", map[string]interface{}{
		"user_id": userID,
		"name":    req.Name,
		"action":  "CREATE_LIFE_AREA",
	})

	lifeArea := &domain.LifeArea{
		UserID:       userID,
		Name:         req.Name,
		Icon:         req.Icon,
		Color:        req.Color,
		DisplayOrder: req.DisplayOrder,
		CreatedAt:    time.Now(),
	}

	created, err := s.repo.Create(ctx, lifeArea)
	if err != nil {
		s.logger.Error("failed to create life area", err, map[string]interface{}{
			"user_id": userID,
			"name":    req.Name,
			"action":  "CREATE_LIFE_AREA_FAILED",
		})
		return nil, err
	}

	s.logger.Info("life area created", map[string]interface{}{
		"life_area_id": created.ID,
		"user_id":      userID,
		"action":       "CREATE_LIFE_AREA",
	})

	return dto.ToLifeAreaResponse(created), nil
}

func (s *lifeAreaService) GetByID(ctx context.Context, id, userID int) (*dto.LifeAreaResponse, error) {
	lifeArea, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lifeArea == nil {
		return nil, errors.New("life area not found")
	}
	if lifeArea.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return dto.ToLifeAreaResponse(lifeArea), nil
}

func (s *lifeAreaService) GetByUserID(ctx context.Context, userID int) ([]*dto.LifeAreaResponse, error) {
	areas, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToLifeAreaResponseList(areas), nil
}

func (s *lifeAreaService) Update(ctx context.Context, id int, req *dto.UpdateLifeAreaRequest, userID int) (*dto.LifeAreaResponse, error) {
	s.logger.Info("Updating life area", map[string]interface{}{
		"user_id":      userID,
		"life_area_id": id,
		"action":       "UPDATE_LIFE_AREA",
	})

	lifeArea, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lifeArea == nil {
		return nil, errors.New("life area not found")
	}
	if lifeArea.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Name != nil {
		lifeArea.Name = *req.Name
	}
	if req.Icon != nil {
		lifeArea.Icon = *req.Icon
	}
	if req.Color != nil {
		lifeArea.Color = *req.Color
	}
	if req.DisplayOrder != nil {
		lifeArea.DisplayOrder = *req.DisplayOrder
	}

	if err := s.repo.Update(ctx, lifeArea); err != nil {
		s.logger.Error("failed to update life area", err, map[string]interface{}{
			"life_area_id": id,
			"user_id":      userID,
			"action":       "UPDATE_LIFE_AREA_FAILED",
		})
		return nil, err
	}

	s.logger.Info("life area updated", map[string]interface{}{
		"life_area_id": id,
		"user_id":      userID,
		"action":       "UPDATE_LIFE_AREA",
	})

	return dto.ToLifeAreaResponse(lifeArea), nil
}

func (s *lifeAreaService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting life area", map[string]interface{}{
		"user_id":      userID,
		"life_area_id": id,
		"action":       "DELETE_LIFE_AREA",
	})

	lifeArea, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if lifeArea == nil {
		return errors.New("life area not found")
	}
	if lifeArea.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete life area", err, map[string]interface{}{
			"life_area_id": id,
			"user_id":      userID,
			"action":       "DELETE_LIFE_AREA_FAILED",
		})
		return err
	}

	s.logger.Info("life area deleted", map[string]interface{}{
		"life_area_id": id,
		"user_id":      userID,
		"action":       "DELETE_LIFE_AREA",
	})

	return nil
}

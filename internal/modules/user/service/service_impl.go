package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo   repository.UserRepository
	logger *logger.ZapLogger
}

func NewUserService(repo repository.UserRepository, logger *logger.ZapLogger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	s.logger.Info("Creating user", map[string]interface{}{
		"email":  req.Email,
		"action": "CREATE_USER",
	})

	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	timezone := req.Timezone
	if timezone == "" {
		timezone = "Europe/Istanbul"
	}

	now := time.Now()
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Timezone:     timezone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", err, map[string]interface{}{
			"email":  req.Email,
			"action": "CREATE_USER_FAILED",
		})
		return nil, err
	}

	s.logger.Info("user created", map[string]interface{}{
		"user_id": created.ID,
		"email":   created.Email,
		"action":  "CREATE_USER",
	})

	return dto.ToUserResponse(created), nil
}

func (s *userService) GetUser(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponseList(users), nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	s.logger.Info("Updating user", map[string]interface{}{
		"user_id": id,
		"action":  "UPDATE_USER",
	})

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Email != nil {
		existing, err := s.repo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user", err, map[string]interface{}{
			"user_id": id,
			"action":  "UPDATE_USER_FAILED",
		})
		return nil, err
	}

	s.logger.Info("user updated", map[string]interface{}{
		"user_id": id,
		"action":  "UPDATE_USER",
	})

	return dto.ToUserResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	s.logger.Info("Deleting user", map[string]interface{}{
		"user_id": id,
		"action":  "DELETE_USER",
	})

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete user", err, map[string]interface{}{
			"user_id": id,
			"action":  "DELETE_USER_FAILED",
		})
		return err
	}

	s.logger.Info("user deleted", map[string]interface{}{
		"user_id": id,
		"action":  "DELETE_USER",
	})

	return nil
}

package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/dto"
	userDomain "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	userRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo userRepo.UserRepository
	logger   *logger.ZapLogger
}

func NewAuthService(userRepo userRepo.UserRepository, logger *logger.ZapLogger) AuthService {
	return &authService{
		userRepo: userRepo,
		logger:   logger,
	}
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	s.logger.Info("Login attempt", map[string]interface{}{
		"email":  req.Email,
		"action": "LOGIN",
	})

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Info("login failed - invalid password", map[string]interface{}{
			"user_id": user.ID,
			"email":   req.Email,
			"action":  "LOGIN_FAILED",
		})
		return nil, errors.New("invalid email or password")
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	s.logger.Info("user logged in", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"action":  "LOGIN",
	})

	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			AvatarURL: user.AvatarURL,
			Timezone:  user.Timezone,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	s.logger.Info("Registering user", map[string]interface{}{
		"email":  req.Email,
		"action": "REGISTER",
	})

	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
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
	user := &userDomain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Timezone:     timezone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", err, map[string]interface{}{
			"email":  req.Email,
			"action": "REGISTER_FAILED",
		})
		return nil, err
	}

	token, expiresAt, err := s.generateToken(created)
	if err != nil {
		return nil, err
	}

	s.logger.Info("user registered", map[string]interface{}{
		"user_id": created.ID,
		"email":   created.Email,
		"action":  "REGISTER",
	})

	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        created.ID,
			Email:     created.Email,
			FullName:  created.FullName,
			AvatarURL: created.AvatarURL,
			Timezone:  created.Timezone,
			CreatedAt: created.CreatedAt,
		},
	}, nil
}

func (s *authService) generateToken(user *userDomain.User) (string, time.Time, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Email,
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

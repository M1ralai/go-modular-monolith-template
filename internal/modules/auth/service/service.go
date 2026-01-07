package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/dto"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
}

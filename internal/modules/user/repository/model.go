package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
)

type UserModel struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     *string   `db:"full_name"`
	AvatarURL    *string   `db:"avatar_url"`
	Timezone     string    `db:"timezone"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (m *UserModel) ToDomain() *domain.User {
	fullName := ""
	if m.FullName != nil {
		fullName = *m.FullName
	}
	avatarURL := ""
	if m.AvatarURL != nil {
		avatarURL = *m.AvatarURL
	}

	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		FullName:     fullName,
		AvatarURL:    avatarURL,
		Timezone:     m.Timezone,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func FromDomain(u *domain.User) *UserModel {
	var fullName, avatarURL *string
	if u.FullName != "" {
		fullName = &u.FullName
	}
	if u.AvatarURL != "" {
		avatarURL = &u.AvatarURL
	}

	return &UserModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     fullName,
		AvatarURL:    avatarURL,
		Timezone:     u.Timezone,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

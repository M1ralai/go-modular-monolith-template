package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/domain"
)

type PersonResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Company      string    `json:"company,omitempty"`
	Relationship string    `json:"relationship,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToPersonResponse(p *domain.Person) *PersonResponse {
	if p == nil {
		return nil
	}
	return &PersonResponse{ID: p.ID, UserID: p.UserID, Name: p.Name, Email: p.Email, Phone: p.Phone, Company: p.Company, Relationship: p.Relationship, Tags: p.Tags, Notes: p.Notes, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func ToPersonResponseList(people []*domain.Person) []*PersonResponse {
	result := make([]*PersonResponse, len(people))
	for i, p := range people {
		result[i] = ToPersonResponse(p)
	}
	return result
}

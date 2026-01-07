package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/domain"
	"github.com/lib/pq"
)

type PersonModel struct {
	ID           int            `db:"id"`
	UserID       int            `db:"user_id"`
	Name         string         `db:"name"`
	Email        *string        `db:"email"`
	Phone        *string        `db:"phone"`
	Company      *string        `db:"company"`
	Relationship *string        `db:"relationship"`
	Tags         pq.StringArray `db:"tags"`
	Notes        *string        `db:"notes"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func (m *PersonModel) ToDomain() *domain.Person {
	if m == nil {
		return nil
	}
	email, phone, company, rel, notes := "", "", "", "", ""
	if m.Email != nil {
		email = *m.Email
	}
	if m.Phone != nil {
		phone = *m.Phone
	}
	if m.Company != nil {
		company = *m.Company
	}
	if m.Relationship != nil {
		rel = *m.Relationship
	}
	if m.Notes != nil {
		notes = *m.Notes
	}
	return &domain.Person{ID: m.ID, UserID: m.UserID, Name: m.Name, Email: email, Phone: phone, Company: company, Relationship: rel, Tags: []string(m.Tags), Notes: notes, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func FromDomain(p *domain.Person) *PersonModel {
	if p == nil {
		return nil
	}
	var email, phone, company, rel, notes *string
	if p.Email != "" {
		email = &p.Email
	}
	if p.Phone != "" {
		phone = &p.Phone
	}
	if p.Company != "" {
		company = &p.Company
	}
	if p.Relationship != "" {
		rel = &p.Relationship
	}
	if p.Notes != "" {
		notes = &p.Notes
	}
	return &PersonModel{ID: p.ID, UserID: p.UserID, Name: p.Name, Email: email, Phone: phone, Company: company, Relationship: rel, Tags: pq.StringArray(p.Tags), Notes: notes, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

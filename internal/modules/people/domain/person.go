package domain

import "time"

type Person struct {
	ID           int
	UserID       int
	Name         string
	Email        string
	Phone        string
	Company      string
	Relationship string
	Tags         []string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (p *Person) HasTag(tag string) bool {
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

package domain

import "time"

type Note struct {
	ID          int
	UserID      int
	CourseID    *int
	ComponentID *int
	LifeAreaID  *int
	Title       string
	Content     string
	IsFavorite  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (n *Note) HasCourse() bool {
	return n.CourseID != nil
}

func (n *Note) HasLifeArea() bool {
	return n.LifeAreaID != nil
}

func (n *Note) GetPreview(maxLength int) string {
	if len(n.Content) <= maxLength {
		return n.Content
	}
	return n.Content[:maxLength] + "..."
}

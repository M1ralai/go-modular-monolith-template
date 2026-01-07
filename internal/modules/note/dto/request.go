package dto

type CreateNoteRequest struct {
	CourseID    *int   `json:"course_id,omitempty"`
	ComponentID *int   `json:"component_id,omitempty"`
	LifeAreaID  *int   `json:"life_area_id,omitempty"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Content     string `json:"content,omitempty"`
}

type UpdateNoteRequest struct {
	CourseID    *int    `json:"course_id,omitempty"`
	ComponentID *int    `json:"component_id,omitempty"`
	LifeAreaID  *int    `json:"life_area_id,omitempty"`
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content     *string `json:"content,omitempty"`
	IsFavorite  *bool   `json:"is_favorite,omitempty"`
}

type CreateNoteLinkRequest struct {
	TargetNoteID int    `json:"target_note_id" validate:"required"`
	LinkText     string `json:"link_text,omitempty"`
}

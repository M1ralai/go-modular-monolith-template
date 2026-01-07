package domain

import "time"

type User struct {
	ID           int
	Email        string
	PasswordHash string
	FullName     string
	AvatarURL    string
	Timezone     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) IsValid() bool {
	return u.Email != "" && u.PasswordHash != ""
}

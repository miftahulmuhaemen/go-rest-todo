package user

import "time"

type User struct {
	ID        string    `json:"_id"`
	RoleID    string    `json:"role_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	Updatet   time.Time `json:"updated_at"`
}

package model

import "time"

type ProfileResponse struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      *string   `db:"name" json:"name"`
	Bio       *string   `db:"bio" json:"bio"`
	Avatar    *string   `db:"avatar" json:"avatar"`
	Phone     *string   `db:"phone" json:"phone"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

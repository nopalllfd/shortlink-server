package model

import "time"

type ProfileResponse struct {
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

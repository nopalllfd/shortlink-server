package dto

import "time"

type Profile struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Name       *string   `json:"name"`
	Bio        *string   `json:"bio"`
	Avatar     *string   `json:"avatar"`
	Phone      *string   `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	LinksTotal int       `json:"links_total"`
}

type UpdateProfileRequest struct {
	Name   *string `json:"name"`
	Bio    *string `json:"bio"`
	Avatar *string `json:"avatar"`
	Phone  *string `json:"phone"`
}

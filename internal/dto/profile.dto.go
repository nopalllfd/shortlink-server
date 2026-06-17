package dto

import "time"

type Profile struct {
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	LinksTotal int       `json:"links_total"`
}

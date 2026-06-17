package dto

import "time"

type CreateLinkRequest struct {
	Link string `json:"link" binding:"required"`
	Slug string `json:"slug,omitempty"`
}

type CreateLinkResponse struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	OriginalUrl string    `json:"original_url"`
	ShortLink   string    `json:"short_link"`
	CreatedAt   time.Time `json:"created_at"`
}

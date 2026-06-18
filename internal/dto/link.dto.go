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
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetLinksWithMeta struct {
	Links []CreateLinkResponse `json:"links"`
	Meta  Meta                 `json:"meta"`
}
type Meta struct {
	Page       int    `json:"page"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	Limit      int    `json:"limit"`
	NextLink   string `json:"next_link"`
	PrevLink   string `json:"prev_link"`
}

type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type DeleteLinkRequest struct {
	ID int `uri:"id" binding:"required"`
}

type RedirectRequest struct {
	Slug string `uri:"slug"`
}

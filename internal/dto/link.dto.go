package dto

import "time"

type CreateLinkRequest struct {
	Link string    `json:"link" binding:"required"`
	Slug string    `json:"slug,omitempty"`
	QR   *QRConfig `json:"qr,omitempty"`
}

type QRConfig struct {
	Size       int    `json:"size,omitempty"`
	Foreground string `json:"foreground,omitempty"`
	Background string `json:"background,omitempty"`
	Style      string `json:"style,omitempty"`
	LogoURL    string `json:"logo_url,omitempty"`
}

type CreateLinkResponse struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	OriginalUrl string    `json:"original_url"`
	ShortLink   string    `json:"short_link"`
	QRUrl       string    `json:"qr_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
	QR          *QRConfig `json:"qr,omitempty"`
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
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

type DeleteLinkRequest struct {
	ID int `uri:"id" binding:"required"`
}

type RedirectRequest struct {
	Slug string `uri:"slug"`
}

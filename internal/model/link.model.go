package model

import "time"

type Link struct {
	ID          int       `db:"id"`
	OriginalUrl string    `db:"original_url"`
	Slug        string    `db:"slug"`
	CreatedAt   time.Time `db:"created_at"`
}

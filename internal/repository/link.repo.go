package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/errs"
	"github.com/nopalllfd/shortlink-server/internal/model"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Create(ctx context.Context, req dto.CreateLinkRequest, userID int) (model.Link, error) {
	query := `INSERT INTO links (user_id, original_url, slug, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, slug, original_url,created_at`

	var link model.Link
	if err := r.db.QueryRow(ctx, query, userID, req.Link, req.Slug).Scan(&link.ID, &link.Slug, &link.OriginalUrl, &link.CreatedAt); err != nil {
		log.Printf("[LinkRepository.Create] error: %v", err)
		return model.Link{}, errs.ErrInternalServer
	}
	return link, nil

}

func (r *LinkRepository) IsSlugExists(ctx context.Context, slug string) (bool, error) {
	query := `EXISTS (SELECT 1 FROM links WHERE slug = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, slug).Scan(&exists); err != nil {
		log.Printf("[LinkRepository.IsSlugExists] error: %v", err)
		return false, errs.ErrInternalServer
	}
	return exists, nil

}

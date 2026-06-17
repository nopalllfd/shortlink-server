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

func (r *LinkRepository) GetByUser(
	ctx context.Context,
	userID int,
	page int,
	limit int,
) ([]dto.CreateLinkResponse, error) {

	offset := (page - 1) * limit

	query := `
		SELECT id, slug, original_url, created_at
		FROM links
		WHERE user_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Printf("[LinkRepository.GetByUser] error: %v", err)
		return nil, errs.ErrInternalServer
	}
	defer rows.Close()

	var result []dto.CreateLinkResponse

	for rows.Next() {
		var link dto.CreateLinkResponse

		if err := rows.Scan(
			&link.ID,
			&link.Slug,
			&link.OriginalUrl,
			&link.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LinkRepository) CountByUser(
	ctx context.Context,
	userID int,
) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM links
		WHERE user_id = $1
		AND deleted_at IS NULL
	`

	var total int

	err := r.db.QueryRow(ctx, query, userID).Scan(&total)
	if err != nil {
		log.Printf("[LinkRepository.CountByUser] error: %v", err)
		return 0, errs.ErrInternalServer
	}

	return total, nil
}

func (r *LinkRepository) Delete(
	ctx context.Context,
	linkID int,
	userID int,
) error {

	query := `
		UPDATE links
		SET deleted_at = NOW()
		WHERE id = $1
		AND user_id = $2
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, linkID, userID)
	if err != nil {
		log.Printf("[LinkRepository.Delete] error: %v", err)
		return errs.ErrInternalServer
	}

	return nil
}

func (r *LinkRepository) GetBySlug(ctx context.Context, slug string) (string, error) {
	query := `SELECT original_url FROM links WHERE slug = $1`

	var original_url string

	if err := r.db.QueryRow(ctx, query, slug).Scan(&original_url); err != nil {
		log.Printf("[LinkRepository.GetBySlug] error: %v", err)
		return "", errs.ErrInternalServer
	}
	return original_url, nil
}

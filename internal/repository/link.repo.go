package repository

import (
	"context"
	"fmt"
	"log"
	"os"

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
	query := `INSERT INTO links (user_id, original_url, slug, created_at, clicks) VALUES ($1, $2, $3, NOW(), 0) RETURNING id, slug, original_url, clicks, created_at`

	var link model.Link
	if err := r.db.QueryRow(ctx, query, userID, req.Link, req.Slug).Scan(&link.ID, &link.Slug, &link.OriginalUrl, &link.Clicks, &link.CreatedAt); err != nil {
		log.Printf("[LinkRepository.Create] error: %v", err)
		return model.Link{}, errs.ErrInternalServer
	}
	return link, nil

}

func (r *LinkRepository) IsSlugExists(ctx context.Context, slug string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM links WHERE slug = $1 AND deleted_at is NULL)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, slug).Scan(&exists); err != nil {
		log.Printf(
			"[LinkRepository.IsSlugExists] failed to check slug=%s error=%v",
			slug,
			err,
		)
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
		SELECT id, slug, original_url, clicks, created_at
		FROM links
		WHERE user_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Printf(
			"[LinkRepository.GetByUser] failed to get links userID=%d page=%d limit=%d offset=%d error=%v",
			userID,
			page,
			limit,
			offset,
			err,
		)
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
			&link.Clicks,
			&link.CreatedAt,
		); err != nil {
			log.Printf(
				"[LinkRepository.GetByUser] failed to scan row userID=%d error=%v",
				userID,
				err,
			)
			return nil, errs.ErrInternalServer
		}
		link.ShortLink = fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), link.Slug)
		result = append(result, link)
	}
	if err := rows.Err(); err != nil {
		log.Printf(
			"[LinkRepository.GetByUser] rows iteration error userID=%d error=%v",
			userID,
			err,
		)
		return nil, errs.ErrInternalServer
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
		log.Printf(
			"[LinkRepository.CountByUser] failed to count links userID=%d error=%v",
			userID,
			err,
		)
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

	result, err := r.db.Exec(
		ctx,
		query,
		linkID,
		userID,
	)

	if err != nil {
		log.Printf(
			"[LinkRepository.Delete] failed to delete linkID=%d userID=%d error=%v",
			linkID,
			userID,
			err,
		)
		return errs.ErrInternalServer
	}

	if result.RowsAffected() == 0 {
		log.Printf(
			"[LinkRepository.Delete] link not found linkID=%d userID=%d",
			linkID,
			userID,
		)
		return errs.ErrLinkNotFound
	}

	return nil
}

func (r *LinkRepository) GetBySlug(ctx context.Context, slug string) (string, error) {
	query := `SELECT original_url FROM links WHERE slug = $1 AND deleted_at is NULL`

	var original_url string

	if err := r.db.QueryRow(ctx, query, slug).Scan(&original_url); err != nil {
		log.Printf(
			"[LinkRepository.GetBySlug] failed to get link slug=%s error=%v",
			slug,
			err,
		)
		return "", err
	}

	return original_url, nil
}

func (r *LinkRepository) GetDeletedByUser(
	ctx context.Context,
	userID int,
	page int,
	limit int,
) ([]dto.CreateLinkResponse, error) {

	offset := (page - 1) * limit

	query := `
		SELECT id, slug, original_url, clicks, created_at
		FROM links
		WHERE user_id = $1
		AND deleted_at IS NOT NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Printf(
			"[LinkRepository.GetDeletedByUser] failed to get deleted links userID=%d page=%d limit=%d offset=%d error=%v",
			userID,
			page,
			limit,
			offset,
			err,
		)
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
			&link.Clicks,
			&link.CreatedAt,
		); err != nil {
			log.Printf(
				"[LinkRepository.GetDeletedByUser] failed to scan row userID=%d error=%v",
				userID,
				err,
			)
			return nil, errs.ErrInternalServer
		}
		link.ShortLink = fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), link.Slug)
		result = append(result, link)
	}
	if err := rows.Err(); err != nil {
		log.Printf(
			"[LinkRepository.GetDeletedByUser] rows iteration error userID=%d error=%v",
			userID,
			err,
		)
		return nil, errs.ErrInternalServer
	}

	return result, nil
}

func (r *LinkRepository) CountDeletedByUser(
	ctx context.Context,
	userID int,
) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM links
		WHERE user_id = $1
		AND deleted_at IS NOT NULL
	`

	var total int

	err := r.db.QueryRow(ctx, query, userID).Scan(&total)
	if err != nil {
		log.Printf(
			"[LinkRepository.CountDeletedByUser] failed to count deleted links userID=%d error=%v",
			userID,
			err,
		)
		return 0, errs.ErrInternalServer
	}

	return total, nil
}

func (r *LinkRepository) IncrementClicks(ctx context.Context, slug string) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE slug = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(ctx, query, slug)
	if err != nil {
		log.Printf(
			"[LinkRepository.IncrementClicks] failed to increment clicks slug=%s error=%v",
			slug,
			err,
		)
		return errs.ErrInternalServer
	}

	return nil
}

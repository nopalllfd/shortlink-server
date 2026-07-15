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

func (r *LinkRepository) Create(
	ctx context.Context,
	req dto.CreateLinkRequest,
	userID *int,
) (model.Link, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("[LinkRepository.Create] failed to begin transaction error: %v", err)
		return model.Link{}, errs.ErrInternalServer
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO links (
			user_id,
			original_url,
			slug,
			created_at,
			clicks
		)
		VALUES ($1, $2, $3, NOW(), 0)
		RETURNING
			id,
			slug,
			original_url,
			clicks,
			created_at
	`

	var link model.Link

	err = tx.QueryRow(
		ctx,
		query,
		userID, // nil -> NULL
		req.Link,
		req.Slug,
	).Scan(
		&link.ID,
		&link.Slug,
		&link.OriginalUrl,
		&link.Clicks,
		&link.CreatedAt,
	)

	if err != nil {
		log.Printf("[LinkRepository.Create] error inserting link: %v", err)
		return model.Link{}, errs.ErrInternalServer
	}

	var fgColor = "#000000"
	var bgColor = "#FFFFFF"
	var dotStyle = "square"
	var eyeStyle = "square"
	var logoURL *string
	var size = 512

	if req.QR != nil {
		if req.QR.Foreground != "" {
			fgColor = req.QR.Foreground
		}
		if req.QR.Background != "" {
			bgColor = req.QR.Background
		}
		if req.QR.Style != "" {
			dotStyle = req.QR.Style
			eyeStyle = req.QR.Style
		}
		if req.QR.LogoURL != "" {
			logoURL = &req.QR.LogoURL
		}
		if req.QR.Size > 0 {
			size = req.QR.Size
		}
	}

	insertQRConfigQuery := `
		INSERT INTO qr_configs (
			link_id,
			foreground_color,
			background_color,
			dot_style,
			eye_style,
			logo_url,
			size
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, insertQRConfigQuery, link.ID, fgColor, bgColor, dotStyle, eyeStyle, logoURL, size)
	if err != nil {
		log.Printf("[LinkRepository.Create] failed to insert qr_config error: %v", err)
		return model.Link{}, errs.ErrInternalServer
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("[LinkRepository.Create] failed to commit transaction error: %v", err)
		return model.Link{}, errs.ErrInternalServer
	}

	return link, nil
}

func (r *LinkRepository) UpdateQRURL(
	ctx context.Context,
	linkID int,
	qrURL string,
) error {

	query := `
		UPDATE links
		SET qr_url = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		qrURL,
		linkID,
	)

	if err != nil {
		log.Printf("[LinkRepository.UpdateQRURL] error: %v", err)
		return errs.ErrInternalServer
	}

	return nil
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
	search string,
) ([]dto.CreateLinkResponse, error) {

	offset := (page - 1) * limit

	query := `
		SELECT 
			l.id, 
			l.slug, 
			l.original_url, 
			l.clicks, 
			l.created_at, 
			l.qr_url,
			q.size,
			q.foreground_color,
			q.background_color,
			q.dot_style,
			q.logo_url
		FROM links l
		LEFT JOIN qr_configs q ON l.id = q.link_id
		WHERE l.user_id = $1
		AND l.deleted_at IS NULL
	`
	args := []interface{}{userID}
	paramIndex := 2

	if search != "" {
		// Build the pattern in Go code to avoid SQL concatenation issues
		pattern := "%" + search + "%"
		log.Printf("[LinkRepository] Applying search pattern: '%s'", pattern)
		query += fmt.Sprintf(` AND (l.original_url ILIKE $%d OR l.slug ILIKE $%d)`, paramIndex, paramIndex)
		args = append(args, pattern)
		paramIndex++
	}

	query += fmt.Sprintf(`
		ORDER BY l.created_at DESC
		LIMIT $%d OFFSET $%d
	`, paramIndex, paramIndex+1)
	args = append(args, limit, offset)

	log.Printf("[LinkRepository.GetByUser] SQL Query: %s", query)
	log.Printf("[LinkRepository.GetByUser] Args: %v", args)

	rows, err := r.db.Query(ctx, query, args...)
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
		var qrURL *string
		var qSize *int
		var qForeground *string
		var qBackground *string
		var qDotStyle *string
		var qLogoURL *string

		if err := rows.Scan(
			&link.ID,
			&link.Slug,
			&link.OriginalUrl,
			&link.Clicks,
			&link.CreatedAt,
			&qrURL,
			&qSize,
			&qForeground,
			&qBackground,
			&qDotStyle,
			&qLogoURL,
		); err != nil {
			log.Printf(
				"[LinkRepository.GetByUser] failed to scan row userID=%d error=%v",
				userID,
				err,
			)
			return nil, errs.ErrInternalServer
		}
		if qrURL != nil {
			link.QRUrl = *qrURL
		}
		if qSize != nil {
			link.QR = &dto.QRConfig{}
			link.QR.Size = *qSize
			if qForeground != nil {
				link.QR.Foreground = *qForeground
			}
			if qBackground != nil {
				link.QR.Background = *qBackground
			}
			if qDotStyle != nil {
				link.QR.Style = *qDotStyle
			}
			if qLogoURL != nil {
				link.QR.LogoURL = *qLogoURL
			}
		}
		link.ShortLink = fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), link.Slug)
		log.Printf("[LinkRepository.GetByUser] Found link: ID=%d, Slug='%s', OriginalUrl='%s'", link.ID, link.Slug, link.OriginalUrl)
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
	search string,
) (int, error) {
	log.Printf("[LinkRepository.CountByUser] Starting, userID=%d, search='%s'", userID, search)

	query := `
		SELECT COUNT(*)
		FROM links
		WHERE user_id = $1
		AND deleted_at IS NULL
	`
	args := []interface{}{userID}
	paramIndex := 2

	if search != "" {
		// Build the pattern in Go code to avoid SQL concatenation issues
		pattern := "%" + search + "%"
		log.Printf("[LinkRepository.CountByUser] Applying search pattern: '%s'", pattern)
		query += fmt.Sprintf(` AND (original_url ILIKE $%d OR slug ILIKE $%d)`, paramIndex, paramIndex)
		args = append(args, pattern)
		paramIndex++
	}

	log.Printf("[LinkRepository.CountByUser] Query: %s, Args: %v", query, args)

	var total int

	err := r.db.QueryRow(ctx, query, args...).Scan(&total)
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
		SELECT 
			l.id, 
			l.slug, 
			l.original_url, 
			l.clicks, 
			l.created_at, 
			l.qr_url,
			q.size,
			q.foreground_color,
			q.background_color,
			q.dot_style,
			q.logo_url
		FROM links l
		LEFT JOIN qr_configs q ON l.id = q.link_id
		WHERE l.user_id = $1
		AND l.deleted_at IS NOT NULL
		ORDER BY l.created_at DESC
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
		var qrURL *string
		var qSize *int
		var qForeground *string
		var qBackground *string
		var qDotStyle *string
		var qLogoURL *string

		if err := rows.Scan(
			&link.ID,
			&link.Slug,
			&link.OriginalUrl,
			&link.Clicks,
			&link.CreatedAt,
			&qrURL,
			&qSize,
			&qForeground,
			&qBackground,
			&qDotStyle,
			&qLogoURL,
		); err != nil {
			log.Printf(
				"[LinkRepository.GetDeletedByUser] failed to scan row userID=%d error=%v",
				userID,
				err,
			)
			return nil, errs.ErrInternalServer
		}
		if qrURL != nil {
			link.QRUrl = *qrURL
		}
		if qSize != nil {
			link.QR = &dto.QRConfig{}
			link.QR.Size = *qSize
			if qForeground != nil {
				link.QR.Foreground = *qForeground
			}
			if qBackground != nil {
				link.QR.Background = *qBackground
			}
			if qDotStyle != nil {
				link.QR.Style = *qDotStyle
			}
			if qLogoURL != nil {
				link.QR.LogoURL = *qLogoURL
			}
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

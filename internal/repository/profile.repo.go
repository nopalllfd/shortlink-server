package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nopalllfd/shortlink-server/internal/errs"
	"github.com/nopalllfd/shortlink-server/internal/model"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) GetById(ctx context.Context, userID int) (model.ProfileResponse, error) {
	query := `SELECT email, created_at FROM users WHERE id = $1`
	var profile model.ProfileResponse
	if err := r.db.QueryRow(ctx, query, userID).Scan(&profile.Email, &profile.CreatedAt); err != nil {
		return model.ProfileResponse{}, err
	}
	return profile, nil
}

func (r *ProfileRepository) CountLinksByUser(
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

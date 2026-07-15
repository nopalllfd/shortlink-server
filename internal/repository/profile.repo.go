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
	query := `SELECT id, email, name, bio, avatar, phone, created_at, updated_at FROM users WHERE id = $1`
	var profile model.ProfileResponse
	if err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.Email,
		&profile.Name,
		&profile.Bio,
		&profile.Avatar,
		&profile.Phone,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		log.Printf("[ProfileRepository.GetById] Failed to get profile: %v", err)
		return model.ProfileResponse{}, err
	}
	return profile, nil
}

func (r *ProfileRepository) UpdateById(ctx context.Context, userID int, name *string, bio *string, avatar *string, phone *string) (model.ProfileResponse, error) {
	query := `
		UPDATE users
		SET 
			name = COALESCE($2, name),
			bio = COALESCE($3, bio),
			avatar = COALESCE($4, avatar),
			phone = COALESCE($5, phone),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, email, name, bio, avatar, phone, created_at, updated_at
	`
	var profile model.ProfileResponse
	if err := r.db.QueryRow(ctx, query, userID, name, bio, avatar, phone).Scan(
		&profile.ID,
		&profile.Email,
		&profile.Name,
		&profile.Bio,
		&profile.Avatar,
		&profile.Phone,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		log.Printf("[ProfileRepository.UpdateById] Failed to update profile: %v", err)
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

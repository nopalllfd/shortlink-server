package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nopalllfd/shortlink-server/internal/model"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(
	ctx context.Context,
	email string,
	passwordHash string,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users(email, password)
		 VALUES ($1, $2)`,
		email,
		passwordHash,
	)

	if err != nil {
		log.Printf("[AuthRepository.Create] error: %v | email=%s", err, email)
		return err
	}

	return nil
}

func (r *AuthRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		log.Printf("[AuthRepository.GetByEmail] error: %v | email=%s", err, email)
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {

	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		log.Printf("[AuthRepository.GetByID] error: %v | id=%d", err, id)
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) EmailExists(
	ctx context.Context,
	email string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		log.Printf("[AuthRepository.EmailExists] error: %v | email=%s", err, email)
		return false, err
	}

	return exists, nil
}

func (r *AuthRepository) ChangePassword(
	ctx context.Context,
	userID int64,
	newPassword string,
) error {

	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, newPassword, userID)

	if err != nil {
		log.Printf("[AuthRepository.ChangePassword] error: %v | userID=%d", err, userID)
		return err
	}

	return nil
}

func (r *AuthRepository) Delete(
	ctx context.Context,
	userID int64,
) error {

	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)

	if err != nil {
		log.Printf("[AuthRepository.Delete] error: %v | userID=%d", err, userID)
		return err
	}

	return nil
}

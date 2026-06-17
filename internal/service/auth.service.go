package service

import (
	"context"
	"log"

	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/errs"
	"github.com/nopalllfd/shortlink-server/internal/repository"
	"github.com/nopalllfd/shortlink-server/pkg"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	log.Println("[Register] START")

	isEmailExists, err := s.authRepo.EmailExists(ctx, req.Email)
	if err != nil {
		log.Printf("[Register] FindByEmail error: %v\n", err)
		return errs.ErrInternalServer
	}

	if isEmailExists {
		log.Printf("[Register] Email already exists: %s\n", req.Email)
		return errs.ErrExistingEmail
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPassword := hc.Hash(req.Password)

	if err := s.authRepo.Create(ctx, req.Email, hashedPassword); err != nil {
		log.Printf("[Register] Create user error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Printf("[Register] User created.")

	return nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	log.Println("[Login] START")

	user, err := s.authRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Println("[Login] service.GetByEmail")
		return dto.AuthResponse{}, err
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	if err := hc.Compare(req.Password, user.Password); err != nil {
		log.Println("[Login] service.Compare")
		return dto.AuthResponse{}, errs.InvalidCredentials
	}

	claims := pkg.NewClaims(int(user.ID), req.Email)

	token, err := claims.GenJWT()
	if err != nil {
		log.Printf("[Login] Generate JWT error: %v\n", err)
		return dto.AuthResponse{}, errs.ErrInternalServer
	}

	log.Println("[Login] SUCCESS")

	data := dto.AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}
	return data, nil

}

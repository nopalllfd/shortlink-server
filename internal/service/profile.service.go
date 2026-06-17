package service

import (
	"context"

	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/repository"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID int) (dto.Profile, error) {
	profile, err := s.profileRepo.GetById(ctx, userID)
	if err != nil {
		return dto.Profile{}, err
	}
	total, err := s.profileRepo.CountLinksByUser(ctx, userID)
	if err != nil {
		return dto.Profile{}, err
	}

	result := dto.Profile{
		Email:      profile.Email,
		CreatedAt:  profile.CreatedAt,
		LinksTotal: total,
	}
	return result, nil
}

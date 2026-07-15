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
		ID:         profile.ID,
		Email:      profile.Email,
		Name:       profile.Name,
		Bio:        profile.Bio,
		Avatar:     profile.Avatar,
		Phone:      profile.Phone,
		CreatedAt:  profile.CreatedAt,
		UpdatedAt:  profile.UpdatedAt,
		LinksTotal: total,
	}
	return result, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID int, req dto.UpdateProfileRequest) (dto.Profile, error) {
	profile, err := s.profileRepo.UpdateById(ctx, userID, req.Name, req.Bio, req.Avatar, req.Phone)
	if err != nil {
		return dto.Profile{}, err
	}
	total, err := s.profileRepo.CountLinksByUser(ctx, userID)
	if err != nil {
		return dto.Profile{}, err
	}

	result := dto.Profile{
		ID:         profile.ID,
		Email:      profile.Email,
		Name:       profile.Name,
		Bio:        profile.Bio,
		Avatar:     profile.Avatar,
		Phone:      profile.Phone,
		CreatedAt:  profile.CreatedAt,
		UpdatedAt:  profile.UpdatedAt,
		LinksTotal: total,
	}
	return result, nil
}

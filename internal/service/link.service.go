package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/errs"
	"github.com/nopalllfd/shortlink-server/internal/repository"
	"github.com/nopalllfd/shortlink-server/pkg"
)

type LinkService struct {
	linkRepo *repository.LinkRepository
}

func NewLinkService(linkRepo *repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

func (s *LinkService) Create(ctx context.Context, req dto.CreateLinkRequest, userID int) (dto.CreateLinkResponse, error) {
	var slug string

	if strings.TrimSpace(req.Slug) == "" {
		genSlug, err := pkg.GenerateRandomSlug(8)
		if err != nil {
			log.Printf("[link.service create] GenerateRandomSlug error: %v\n", err)
			return dto.CreateLinkResponse{}, errs.ErrInternalServer
		}
		slug = genSlug
	} else {
		slug = req.Slug
		exists, err := s.linkRepo.IsSlugExists(ctx, slug)
		if err != nil {
			log.Printf("[link.service create] IsSlugExists error: %v\n", err)
			return dto.CreateLinkResponse{}, errs.ErrInternalServer
		}
		if exists {
			return dto.CreateLinkResponse{}, errs.ErrSlugAlreadyExists
		}
	}
	req.Slug = slug
	data, err := s.linkRepo.Create(ctx, req, userID)
	if err != nil {
		log.Printf("[link.service create] Create error: %v\n", err)
		return dto.CreateLinkResponse{}, errs.ErrInternalServer
	}
	shortLink := fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), data.Slug)

	return dto.CreateLinkResponse{
		ID:          data.ID,
		Slug:        data.Slug,
		OriginalUrl: data.OriginalUrl,
		ShortLink:   shortLink,
	}, nil
}

package service

import (
	"context"
	"fmt"
	"log"
	"math"
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

func (s *LinkService) GetAll(
	ctx context.Context,
	userID int,
	page int,
	limit int,
) (dto.GetLinksWithMeta, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	links, err := s.linkRepo.GetByUser(
		ctx,
		userID,
		page,
		limit,
	)
	if err != nil {
		log.Printf("[LinkService.GetAll] error: %v", err)
		return dto.GetLinksWithMeta{}, errs.ErrInternalServer
	}

	total, err := s.linkRepo.CountByUser(ctx, userID)
	if err != nil {
		log.Printf("[LinkService.GetAll] error: %v", err)
		return dto.GetLinksWithMeta{}, errs.ErrInternalServer
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	meta := dto.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	if page > 1 {
		meta.PrevLink = fmt.Sprintf(
			"/links?page=%d&limit=%d",
			page-1,
			limit,
		)
	}

	if page < totalPages {
		meta.NextLink = fmt.Sprintf(
			"/links?page=%d&limit=%d",
			page+1,
			limit,
		)
	}

	return dto.GetLinksWithMeta{
		Links: links,
		Meta:  meta,
	}, nil
}

func (s *LinkService) Delete(
	ctx context.Context,
	linkID int,
	userID int,
) error {

	err := s.linkRepo.Delete(
		ctx,
		linkID,
		userID,
	)
	if err != nil {
		log.Printf("[LinkService.Delete] error: %v", err)
		return err
	}

	return nil
}

func (s *LinkService) GetBySlug(ctx context.Context, slug string) (string, error) {
	return s.linkRepo.GetBySlug(ctx, slug)
}

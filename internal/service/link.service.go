package service

import (
	"context"
	"errors"
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
			log.Printf(
				"[LinkService.Create] failed to generate slug userID=%d error=%v",
				userID,
				err,
			)
			return dto.CreateLinkResponse{}, errs.ErrInternalServer
		}
		slug = genSlug
	} else {
		reserved := []string{"api", "login", "register", "dashboard"}
		slug = req.Slug
		for _, v := range reserved {
			if req.Slug == v {
				return dto.CreateLinkResponse{}, errs.ErrCannotUserReserveWord
			}
		}
		if len(req.Slug) < 6 {
			return dto.CreateLinkResponse{}, errs.ErrMinimumSlug

		}
		exists, err := s.linkRepo.IsSlugExists(ctx, slug)
		if err != nil {
			log.Printf(
				"[LinkService.Create] failed checking slug=%s userID=%d error=%v",
				slug,
				userID,
				err,
			)
			return dto.CreateLinkResponse{}, errs.ErrInternalServer
		}
		if exists {
			return dto.CreateLinkResponse{}, errs.ErrSlugAlreadyExists
		}
	}

	req.Slug = slug
	data, err := s.linkRepo.Create(ctx, req, userID)
	if err != nil {
		log.Printf(
			"[LinkService.Create] failed creating link userID=%d slug=%s error=%v",
			userID,
			slug,
			err,
		)
		return dto.CreateLinkResponse{}, errs.ErrInternalServer
	}
	shortLink := fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), data.Slug)

	return dto.CreateLinkResponse{
		ID:          data.ID,
		Slug:        data.Slug,
		OriginalUrl: data.OriginalUrl,
		ShortLink:   shortLink,
		Clicks:      data.Clicks,
		CreatedAt:   data.CreatedAt,
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
		log.Printf(
			"[LinkService.GetAll] failed getting links userID=%d page=%d limit=%d error=%v",
			userID,
			page,
			limit,
			err,
		)
		return dto.GetLinksWithMeta{}, errs.ErrInternalServer
	}

	total, err := s.linkRepo.CountByUser(ctx, userID)
	if err != nil {
		log.Printf(
			"[LinkService.GetAll] failed counting links userID=%d error=%v",
			userID,
			err,
		)
		return dto.GetLinksWithMeta{}, errs.ErrInternalServer
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	meta := dto.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
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
		log.Printf(
			"[LinkService.Delete] failed deleting linkID=%d userID=%d error=%v",
			linkID,
			userID,
			err,
		)
		return err
	}

	return nil
}

func (s *LinkService) GetBySlug(ctx context.Context, slug string) (string, error) {
	link, err := s.linkRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, errs.ErrSlugNotFound) {
			return "", errs.ErrSlugNotFound
		}
	}
	return link, nil
}

func (s *LinkService) IsSlugExists(ctx context.Context, slug string) (bool, error) {
	return s.linkRepo.IsSlugExists(ctx, slug)
}

func (s *LinkService) GetAllDeleted(
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

	links, err := s.linkRepo.GetDeletedByUser(
		ctx,
		userID,
		page,
		limit,
	)
	if err != nil {
		log.Printf(
			"[LinkService.GetAllDeleted] failed getting deleted links userID=%d page=%d limit=%d error=%v",
			userID,
			page,
			limit,
			err,
		)
		return dto.GetLinksWithMeta{}, errs.ErrInternalServer
	}

	total, err := s.linkRepo.CountDeletedByUser(ctx, userID)
	if err != nil {
		log.Printf(
			"[LinkService.GetAllDeleted] failed counting deleted links userID=%d error=%v",
			userID,
			err,
		)
		return dto.GetLinksWithMeta{}, errs.ErrInternalServer
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	meta := dto.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	if page > 1 {
		meta.PrevLink = fmt.Sprintf(
			"/links/deleted?page=%d&limit=%d",
			page-1,
			limit,
		)
	}

	if page < totalPages {
		meta.NextLink = fmt.Sprintf(
			"/links/deleted?page=%d&limit=%d",
			page+1,
			limit,
		)
	}

	return dto.GetLinksWithMeta{
		Links: links,
		Meta:  meta,
	}, nil
}

func (s *LinkService) IncrementClicks(ctx context.Context, slug string) error {
	err := s.linkRepo.IncrementClicks(ctx, slug)
	if err != nil {
		log.Printf(
			"[LinkService.IncrementClicks] failed to increment clicks slug=%s error=%v",
			slug,
			err,
		)
		return err
	}
	return nil
}

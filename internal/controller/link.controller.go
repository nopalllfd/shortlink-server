package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/errs"
	"github.com/nopalllfd/shortlink-server/internal/response"
	"github.com/nopalllfd/shortlink-server/internal/service"
	"github.com/nopalllfd/shortlink-server/pkg"
)

type LinkController struct {
	linkService *service.LinkService
}

func NewLinkController(linkService *service.LinkService) *LinkController {
	return &LinkController{
		linkService: linkService,
	}
}

func (c *LinkController) CreateShortLink(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims, _ := token.(pkg.Claims)

	var req dto.CreateLinkRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf(
			"[LinkController.CreateShortLink] failed to bind request body error=%v",
			err,
		)

		errMsg := err.Error()

		if strings.Contains(errMsg, "required") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"field is required",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			errMsg,
		)
		return
	}

	userID := claims.Id

	data, err := c.linkService.Create(
		ctx.Request.Context(),
		req,
		&userID,
	)

	if err != nil {

		log.Printf(
			"[LinkController.CreateShortLink] failed to create short link userID=%d slug=%s error=%v",
			userID,
			req.Slug,
			err,
		)
		if errors.Is(err, errs.ErrCannotUserReserveWord) {
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}
		if errors.Is(err, errs.ErrMinimumSlug) {
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		if errors.Is(err, errs.ErrSlugAlreadyExists) {
			response.Error(
				ctx,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success create link",
		data,
	)
}

func (c *LinkController) CreateShortLinkPublic(ctx *gin.Context) {

	var req dto.CreateLinkRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf(
			"[LinkController.CreateShortLink] failed to bind request body error=%v",
			err,
		)

		errMsg := err.Error()

		if strings.Contains(errMsg, "required") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"field is required",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			errMsg,
		)
		return
	}

	var userID *int

	data, err := c.linkService.Create(
		ctx.Request.Context(),
		req,
		userID,
	)

	if err != nil {

		log.Printf(
			"[LinkController.CreateShortLink] failed to create short link userID=%d slug=%s error=%v",
			userID,
			req.Slug,
			err,
		)
		if errors.Is(err, errs.ErrCannotUserReserveWord) {
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}
		if errors.Is(err, errs.ErrMinimumSlug) {
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		if errors.Is(err, errs.ErrSlugAlreadyExists) {
			response.Error(
				ctx,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success create link",
		data,
	)
}

func (c *LinkController) GetAllLinks(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	log.Printf(
		"[LinkController.GetAllLinks] started userID=%d",
		claims.Id,
	)

	var req dto.PaginationQuery

	if err := ctx.ShouldBind(&req); err != nil {

		log.Printf(
			"[LinkController.GetAllLinks] failed to bind query userID=%d error=%v",
			claims.Id,
			err,
		)

		errMsg := err.Error()

		if strings.Contains(errMsg, "required") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"field is required",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			errMsg,
		)
		return
	}

	userID := claims.Id
	log.Printf("[LinkController.GetAllLinks] Full req object: %+v", req)
	log.Printf("[LinkController.GetAllLinks] Received search query: '%s'", req.Search)

	data, err := c.linkService.GetAll(
		ctx.Request.Context(),
		userID,
		req.Page,
		req.Limit,
		req.Search,
	)

	if err != nil {

		log.Printf(
			"[LinkController.GetAllLinks] failed to get links userID=%d page=%d limit=%d error=%v",
			userID,
			req.Page,
			req.Limit,
			err,
		)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success get all links",
		data,
	)
}

func (c *LinkController) DeleteLink(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var req dto.DeleteLinkRequest

	if err := ctx.ShouldBindUri(&req); err != nil {

		log.Printf(
			"[LinkController.DeleteLink] invalid uri parameter error=%v",
			err,
		)

		response.Error(
			ctx,
			http.StatusBadRequest,
			"invalid link id",
		)
		return
	}

	err := c.linkService.Delete(
		ctx.Request.Context(),
		req.ID,
		claims.Id,
	)

	if err != nil {

		log.Printf(
			"[LinkController.DeleteLink] failed to delete linkID=%d userID=%d error=%v",
			req.ID,
			claims.Id,
			err,
		)

		if errors.Is(err, errs.ErrLinkNotFound) {
			response.Error(
				ctx,
				http.StatusNotFound,
				err.Error(),
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success delete link",
		nil,
	)
}

func (c *LinkController) Redirect(ctx *gin.Context) {
	var req dto.RedirectRequest

	if err := ctx.ShouldBindUri(&req); err != nil {

		log.Printf(
			"[LinkController.Redirect] invalid slug parameter error=%v",
			err,
		)

		ctx.Redirect(http.StatusFound, "https://home.yurl.ink/notfound")
		return
	}

	link, err := c.linkService.GetBySlug(
		ctx,
		req.Slug,
	)

	if err != nil {

		log.Printf(
			"[LinkController.Redirect] failed to get link slug=%s error=%v",
			req.Slug,
			err,
		)
		ctx.Redirect(http.StatusFound, "https://home.yurl.ink/notfound")
		return
	}

	// Increment clicks asynchronously or just call it
	go func() {
		// Use background context since original might be canceled after redirect
		_ = c.linkService.IncrementClicks(context.Background(), req.Slug)
	}()

	log.Printf(
		"[LinkController.Redirect] redirecting slug=%s destination=%s",
		req.Slug,
		link,
	)

	ctx.Redirect(
		http.StatusMovedPermanently,
		link,
	)
}
func (c *LinkController) CheckSlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.Error(ctx, http.StatusBadRequest, "slug must be filled")
		return
	}

	exist, err := c.linkService.IsSlugExists(ctx.Request.Context(), slug)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "slug successfully checked", exist)
}

func (c *LinkController) GetAllDeletedLinks(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	log.Printf(
		"[LinkController.GetAllDeletedLinks] started userID=%d",
		claims.Id,
	)

	var req dto.PaginationQuery

	if err := ctx.ShouldBind(&req); err != nil {

		log.Printf(
			"[LinkController.GetAllDeletedLinks] failed to bind query userID=%d error=%v",
			claims.Id,
			err,
		)

		errMsg := err.Error()

		if strings.Contains(errMsg, "required") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"field is required",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			errMsg,
		)
		return
	}

	userID := claims.Id

	data, err := c.linkService.GetAllDeleted(
		ctx.Request.Context(),
		userID,
		req.Page,
		req.Limit,
	)

	if err != nil {

		log.Printf(
			"[LinkController.GetAllDeletedLinks] failed to get deleted links userID=%d page=%d limit=%d error=%v",
			userID,
			req.Page,
			req.Limit,
			err,
		)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success get all deleted links",
		data,
	)
}

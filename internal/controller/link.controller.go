package controller

import (
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
		errMsg := err.Error()
		if strings.Contains(errMsg, "required") {
			response.Error(ctx, http.StatusBadRequest, "field is required")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, errMsg)
		return
	}

	userID := claims.Id

	data, err := c.linkService.Create(ctx.Request.Context(), req, userID)
	if err != nil {
		if errors.Is(err, errs.ErrSlugAlreadyExists) {
			response.Error(ctx, http.StatusConflict, err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "success create link", data)
}

func (c *LinkController) GetAllLinks(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	log.Printf("[LinkController.started]")

	var req dto.PaginationQuery
	if err := ctx.ShouldBind(&req); err != nil {

		errMsg := err.Error()
		log.Printf("[LinkController.GetAllLinks] error: %v", errMsg)
		if strings.Contains(errMsg, "required") {
			response.Error(ctx, http.StatusBadRequest, "field is required")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, errMsg)
		return
	}

	userID := claims.Id

	data, err := c.linkService.GetAll(ctx.Request.Context(), userID, req.Page, req.Limit)
	if err != nil {
		log.Printf("[LinkController.GetAllLinks] error: %v", err.Error())

		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "success get all links", data)

}

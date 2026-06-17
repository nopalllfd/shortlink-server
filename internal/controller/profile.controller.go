package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/shortlink-server/internal/response"
	"github.com/nopalllfd/shortlink-server/internal/service"
	"github.com/nopalllfd/shortlink-server/pkg"
)

type ProfileController struct {
	profileService *service.ProfileService
}

func NewProfileController(profileService *service.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

func (c *ProfileController) GetProfileByUserId(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims, _ := token.(pkg.Claims)

	user, err := c.profileService.GetProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
	}
	response.Success(ctx, http.StatusOK, "success to get profile", user)
}

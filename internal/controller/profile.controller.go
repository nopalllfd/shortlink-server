package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/shortlink-server/internal/dto"
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

func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims, _ := token.(pkg.Claims)

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ProfileController.UpdateProfile] Failed to bind request: %v", err)
		response.Error(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := c.profileService.UpdateProfile(ctx.Request.Context(), claims.Id, req)
	if err != nil {
		log.Printf("[ProfileController.UpdateProfile] Failed to update profile: %v", err)
		response.Error(ctx, http.StatusInternalServerError, "failed to update profile")
		return
	}
	response.Success(ctx, http.StatusOK, "success to update profile", profile)
}

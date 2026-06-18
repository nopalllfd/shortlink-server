package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/errs"
	"github.com/nopalllfd/shortlink-server/internal/response"
	"github.com/nopalllfd/shortlink-server/internal/service"
	"github.com/nopalllfd/shortlink-server/pkg"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		errMsg := err.Error()

		if strings.Contains(errMsg, "required") {
			response.Error(ctx, http.StatusBadRequest, "field is required")
			return
		}

		if strings.Contains(errMsg, "email") {
			response.Error(ctx, http.StatusBadRequest, "invalid email format")
			return
		}

		if strings.Contains(errMsg, "min") {
			response.Error(ctx, http.StatusBadRequest, "password must be min 8 chars")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, errMsg)
		return
	}
	if err := c.authService.Register(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "registered success", nil)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		errMsg := err.Error()

		if strings.Contains(errMsg, "required") {
			response.Error(ctx, http.StatusBadRequest, "field is required")
			return
		}

		if strings.Contains(errMsg, "email") {
			response.Error(ctx, http.StatusBadRequest, "invalid email format")
			return
		}

		if strings.Contains(errMsg, "min") {
			response.Error(ctx, http.StatusBadRequest, "password must be min 8 chars")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, errMsg)
		return
	}

	user, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, errs.InvalidCredentials) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, errs.ErrUserNotFound) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "login success", user)

}

func (c *AuthController) Logout(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	data := dto.LogoutRequest{
		Token:     tokenString,
		ExpiredAt: claims.ExpiresAt.Time,
	}
	if err := c.authService.Logout(ctx.Request.Context(), data); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())

		return
	}
	response.Success(ctx, http.StatusOK, "logout success", nil)

}

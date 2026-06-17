package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nopalllfd/shortlink-server/internal/controller"
	"github.com/nopalllfd/shortlink-server/internal/repository"
	"github.com/nopalllfd/shortlink-server/internal/service"
)

func RegisterAuthRoute(rg *gin.RouterGroup, db *pgxpool.Pool) {
	AuthRepo := repository.NewAuthRepository(db)
	AuthService := service.NewAuthService(AuthRepo)
	AuthController := controller.NewAuthController(AuthService)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", AuthController.Register)
		auth.POST("/login", AuthController.Login)
	}
}

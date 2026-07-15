package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nopalllfd/shortlink-server/internal/controller"
	"github.com/nopalllfd/shortlink-server/internal/middleware"
	"github.com/nopalllfd/shortlink-server/internal/repository"
	"github.com/nopalllfd/shortlink-server/internal/service"
	"github.com/redis/go-redis/v9"
)

func RegisterProfileRoute(rg *gin.RouterGroup, db *pgxpool.Pool, rc *redis.Client) {
	ProfileRepo := repository.NewProfileRepository(db)
	ProfileService := service.NewProfileService(ProfileRepo)
	ProfileController := controller.NewProfileController(ProfileService)

	profile := rg.Group("/profiles")
	{
		profile.GET("", middleware.VerifyMiddleware(rc), ProfileController.GetProfileByUserId)
		profile.PUT("", middleware.VerifyMiddleware(rc), ProfileController.UpdateProfile)
	}
}

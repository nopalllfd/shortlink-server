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

func RegisterLinkRoute(rg *gin.RouterGroup, db *pgxpool.Pool, rc *redis.Client) {
	LinkRepo := repository.NewLinkRepository(db)
	LinkService := service.NewLinkService(LinkRepo)
	LinkController := controller.NewLinkController(LinkService)

	link := rg.Group("/links")
	{
		link.POST("", middleware.VerifyMiddleware(rc), LinkController.CreateShortLink)
	}
}

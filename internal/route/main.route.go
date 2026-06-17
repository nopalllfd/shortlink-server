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

func InitRoute(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	app.Use(middleware.CORSMiddleware())
	LinkRepo := repository.NewLinkRepository(db)
	LinkService := service.NewLinkService(LinkRepo)
	LinkController := controller.NewLinkController(LinkService)
	app.GET("/:slug", LinkController.Redirect)
	api := app.Group("/api")
	RegisterAuthRoute(api, db)
	RegisterLinkRoute(api, db, rc)
	RegisterProfileRoute(api, db, rc)
}

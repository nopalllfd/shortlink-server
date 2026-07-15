package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nopalllfd/shortlink-server/internal/controller"
	"github.com/nopalllfd/shortlink-server/internal/middleware"
	"github.com/nopalllfd/shortlink-server/internal/repository"
	"github.com/nopalllfd/shortlink-server/internal/service"
	"github.com/nopalllfd/shortlink-server/internal/storage"
	"github.com/redis/go-redis/v9"
)

func InitRoute(
	app *gin.Engine,
	db *pgxpool.Pool,
	rc *redis.Client,
	objectStorage storage.ObjectStorage,
) {
	app.Use(middleware.CORSMiddleware())

	linkRepo := repository.NewLinkRepository(db)

	qrService := service.NewQRService(objectStorage)

	linkService := service.NewLinkService(
		linkRepo,
		qrService,
	)

	linkController := controller.NewLinkController(linkService)

	app.GET("/:slug", linkController.Redirect)

	api := app.Group("/api")

	RegisterAuthRoute(api, db, rc)
	RegisterProfileRoute(api, db, rc)

	RegisterLinkRoute(
		api,
		rc,
		linkController,
	)
}

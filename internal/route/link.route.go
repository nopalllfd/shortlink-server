package route

import (
	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/shortlink-server/internal/controller"
	"github.com/nopalllfd/shortlink-server/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func RegisterLinkRoute(
	rg *gin.RouterGroup,
	rc *redis.Client,
	linkController *controller.LinkController,
) {
	link := rg.Group("/links")
	{
		link.POST("", middleware.VerifyMiddleware(rc), linkController.CreateShortLink)
		link.POST("/public", linkController.CreateShortLinkPublic)
		link.GET("/:slug", linkController.CheckSlug)
		link.GET("", middleware.VerifyMiddleware(rc), linkController.GetAllLinks)
		link.GET("/deleted", middleware.VerifyMiddleware(rc), linkController.GetAllDeletedLinks)
		link.DELETE("/:id", middleware.VerifyMiddleware(rc), linkController.DeleteLink)
	}
}

package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitRoute(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	api := app.Group("/api")
	RegisterAuthRoute(api, db)
	RegisterLinkRoute(api, db, rc)
}

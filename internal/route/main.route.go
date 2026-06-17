package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRoute(app *gin.Engine, db *pgxpool.Pool) {
	api := app.Group("/api")
	RegisterAuthRoute(api, db)
}

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nopalllfd/shortlink-server/internal/config"
	"github.com/nopalllfd/shortlink-server/internal/route"
)

func main() {
	_ = godotenv.Load()

	app := gin.Default()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("DB Connection error : %s", err.Error())
	}

	defer db.Close()

	rc, err := config.ConnectRedis()
	if err != nil {
		log.Fatalf("Redis Connection error : %s", err.Error())
	}

	defer rc.Close()

	route.InitRoute(app, db, rc)

	app.Run("0.0.0.0:8080")

}

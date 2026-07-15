package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nopalllfd/shortlink-server/internal/config"
	"github.com/nopalllfd/shortlink-server/internal/route"
	"github.com/nopalllfd/shortlink-server/internal/storage"
)

func main() {
	_ = godotenv.Load()

	app := gin.Default()

	r2Client, err := config.ConnectR2()
	if err != nil {
		log.Fatalf("R2 Connection error : %v", err)
	}

	r2Storage := storage.NewR2Storage(
		r2Client,
		os.Getenv("R2_BUCKET"),
		os.Getenv("R2_PUBLIC_URL"),
	)

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

	route.InitRoute(app, db, rc, r2Storage)

	app.Run("0.0.0.0:8080")

}

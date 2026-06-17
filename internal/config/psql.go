package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {

	values := make([]any, 0, 5)
	values = append(values, os.Getenv("DB_USER"))
	values = append(values, os.Getenv("DB_PASS"))
	values = append(values, os.Getenv("DB_HOST"))
	values = append(values, os.Getenv("DB_PORT"))
	values = append(values, os.Getenv("DB_NAME"))

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", values...)
	// connStr := "postgresql://shortlink_user:secret123@127.0.0.1:6666/shortlink?sslmode=disable"
	return pgxpool.New(context.Background(), connStr)
}

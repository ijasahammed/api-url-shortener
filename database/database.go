package database

import (
	"context"
	"os"

	"github.com/go-redis/redis"
)

var Ctx = context.Background()

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	URL              string `json:"url"`
	CustomedShortURL string `json:"short_url"`
}

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"), // Redis server address
		Password: os.Getenv("DB_PASS"), // Redis server password
		DB:       dbNo,                 // Database number
	})
	return rdb
}

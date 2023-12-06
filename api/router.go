package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"api-url-shortener/database"

	"github.com/joho/godotenv"
)

type Repository struct {
	ShortUrlDBClient  *redis.Client
	HostCountDBClient *redis.Client
}

func InitializeApp(app *gin.Engine) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	shortUrlClient := database.CreateClient(0)
	hostCountDBClient := database.CreateClient(1)

	repo := Repository{
		ShortUrlDBClient:  shortUrlClient,
		HostCountDBClient: hostCountDBClient,
	}

	repo.SetupRoutes(app)
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}

func (repo *Repository) SetupRoutes(app *gin.Engine) {
	app.POST("/shorten", repo.ShortenURL)
}

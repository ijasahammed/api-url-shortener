package main

import (
	"api-url-shortener/api"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	api.InitializeApp(app)
}

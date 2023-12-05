package main

import (
	"github.com/gin-gonic/gin"
	"api-url-shortener/api"
)

func main() {
	app := gin.Default()
	api.InitializeApp(app)
}

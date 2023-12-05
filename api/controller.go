package api

import (
	"os"
	"fmt"
	"strconv"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"api-url-shortener/database"
	"api-url-shortener/internal/helpers"
)

func (repo *Repository) ShortenURL(c *gin.Context) {

	body := database.Request{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(500, gin.H{
			"Error": err,
		})
		fmt.Println(err)

		return
	}

	// Check if the input URL is a valid URL
	if !govalidator.IsURL(body.Url) {
		c.JSON(400, gin.H{"Error": "Invalid URL"})
	}

	// Handle domain error
	valid,host := helpers.RemoveDomainError(body.Url)
	if !valid {
		c.JSON(400, gin.H{"Error": "Remove domain error"})
	}

	val,err := repo.HostCountDBClient.Get(host).Result()
	if err != nil {
		c.JSON(400, gin.H{"Error": "Get host based count error"})
	}

	count := 0

	if val == ""{
		count = 0
	}else{
		countInt,_ := strconv.Atoi(val)
		count = countInt + 1
	}
	err = repo.HostCountDBClient.Set(host, count,0).Err()

	fmt.Println(host,count)

	id := uuid.New().String()[:6]

	// Enforce HTTPS
	body.Url = helpers.EnforceHTTP(body.Url)

	err = repo.ShortUrlDBClient.Set(id, body.Url,0).Err()
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Unable to connect to the database",
		})
		return
	}

	// Prepare the response
	resp := database.Response{
		URL:              body.Url,
		CustomedShortURL: "",
	}

	resp.CustomedShortURL = os.Getenv("SHORT_BASE_URL") + "/" + id

	c.JSON(200, resp)
}
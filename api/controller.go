package api

import (
	"api-url-shortener/database"
	"api-url-shortener/internal/helpers"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var hostNameKey, urlNamekey = "host_count", "url_data"

func (repo *Repository) ShortenURL(c *gin.Context) {

	body := database.Request{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(500, gin.H{
			"Error": err,
		})

		return
	}

	// Check if the input URL is a valid URL
	if !govalidator.IsURL(body.Url) {
		c.JSON(400, gin.H{"Error": "Invalid URL"})
		return
	}

	// Handle domain error
	valid, host := helpers.RemoveDomainError(body.Url)
	if !valid {
		c.JSON(400, gin.H{"Error": "URL is already short one"})
		return
	}

	// Host based count
	val, err := repo.ShortUrlDBClient.Get(hostNameKey).Result()
	if err != nil && err != redis.Nil {
		c.JSON(400, gin.H{"Error": "Unable to connect to the database"})
		return
	}

	count := 1

	hostCountMap := map[string]int{}

	if val != "" {
		err = json.Unmarshal([]byte(val), &hostCountMap)
		if err != nil {
			c.JSON(400, gin.H{"Error": "Data conversion error"})
			return
		}
		if countInt, exists := hostCountMap[host]; exists {
			count = countInt + 1
		}
	}
	hostCountMap[host] = count
	hostCountJson, err := json.Marshal(hostCountMap)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Data conversion error"})
		return
	}

	err = repo.ShortUrlDBClient.Set(hostNameKey, hostCountJson, 0).Err()

	if err != nil {
		c.JSON(400, gin.H{"Error": "Unable to connect to the database"})
		return
	}

	// url data
	urlDataMap := map[string]string{}

	id := uuid.New().String()[:6]

	// Enforce HTTPS
	body.Url = helpers.EnforceHTTP(body.Url)

	exists, oldID, err := helpers.CheckURLAlreadyExists(repo.ShortUrlDBClient, body.Url, urlNamekey)

	fmt.Println(exists, oldID, err)

	if err != nil {
		c.JSON(400, gin.H{"Error": "Internal problem when checking already exists"})
		return
	}
	if exists {
		resp := database.Response{
			URL:              body.Url,
			CustomedShortURL: "",
		}
		resp.CustomedShortURL = os.Getenv("SHORT_BASE_URL") + "/" + oldID

		c.JSON(200, resp)
		return
	}

	val, err = repo.ShortUrlDBClient.Get(urlNamekey).Result()
	if err != nil && err != redis.Nil {
		c.JSON(400, gin.H{"Error": "Unable to connect to the database"})
		return
	}
	if val != "" {
		err = json.Unmarshal([]byte(val), &urlDataMap)
		if err != nil {
			c.JSON(400, gin.H{"Error": "Data conversion error"})
			return
		}
		for {
			if _, exists := urlDataMap[host]; !exists {
				break
			}
		}
	}
	urlDataMap[id] = body.Url
	urlDataJson, err := json.Marshal(urlDataMap)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Data conversion error"})
		return
	}

	err = repo.ShortUrlDBClient.Set(urlNamekey, urlDataJson, 0).Err()
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

func (repo *Repository) GetHostCount(c *gin.Context) {
	val, err := repo.ShortUrlDBClient.Get(hostNameKey).Result()
	if err != nil && err != redis.Nil {
		c.JSON(400, gin.H{"Error": "Unable to connect to the database"})
		return
	}

	hostCountMap := map[string]int{}

	if val != "" {
		err = json.Unmarshal([]byte(val), &hostCountMap)
		if err != nil {
			c.JSON(400, gin.H{"Error": "Data conversion error"})
			return
		}
	}

	// Create slice of key-value sortData
	sortData := make([][2]interface{}, 0, len(hostCountMap))
	for k, v := range hostCountMap {
		sortData = append(sortData, [2]interface{}{k, v})
	}

	// Sort slice based on values
	sort.Slice(sortData, func(i, j int) bool {
		return sortData[i][1].(int) < sortData[j][1].(int)
	})

	if len(sortData) > 3 {
		sortData = sortData[:3]
	}

	returnData := map[string]int{}

	for _, key := range sortData {
		returnData[key[0].(string)] = key[1].(int)
	}

	resp := database.CountResponse{
		Data: returnData,
	}
	c.JSON(200, resp)

}

func (repo *Repository) ResolveURL(c *gin.Context) {
	urlDataMap := map[string]string{}

	shortUrl := c.Param("url")
	var redirURL string

	val, err := repo.ShortUrlDBClient.Get(urlNamekey).Result()
	if err != nil && err != redis.Nil {
		c.JSON(400, gin.H{"Error": "connect error redis"})
		return
	}
	if val != "" {
		err = json.Unmarshal([]byte(val), &urlDataMap)
		if err != nil {
			c.JSON(400, gin.H{"Error": "Data conversion error"})
			return
		}
		if _, exists := urlDataMap[shortUrl]; !exists {
			c.JSON(404, gin.H{
				"Error": "No URL mapped to this short URL",
			})
			return
		}
		redirURL = urlDataMap[shortUrl]
	} else {
		c.JSON(404, gin.H{
			"Error": "No URL mapped to this short URL",
		})
		return
	}
	c.Redirect(301, redirURL)
}

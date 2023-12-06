package helpers

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

// EnforceHTTP enforces the use of "http://" prefix for a URL if not present.
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) (bool, string) {
	if url == os.Getenv("SHORT_BASE_URL") {
		return false, ""
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("SHORT_BASE_URL") {
		return false, ""
	}

	return true, newURL
}

func CheckURLAlreadyExists(client *redis.Client, url, key string) (bool, string, error) {

	urlDataMap := map[string]string{}
	var orgUrl string

	val, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		return false, "", err
	}
	if val != "" {
		err = json.Unmarshal([]byte(val), &urlDataMap)
		if err != nil {
			return false, "", err
		}
		for key, val := range urlDataMap {
			if val == url {
				orgUrl = key
			}
		}
	} else {
		return false, "", nil
	}
	return true, orgUrl, nil

}

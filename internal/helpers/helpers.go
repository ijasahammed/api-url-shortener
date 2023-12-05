package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP enforces the use of "http://" prefix for a URL if not present.
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) (bool,string) {
	if url == os.Getenv("SHORT_BASE_URL") {
		return false,""
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("SHORT_BASE_URL") {
		return false,""
	}

	return true,newURL
}
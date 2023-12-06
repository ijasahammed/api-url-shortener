// helpers_test.go

package helpers_test

import (
	"testing"
	"encoding/json"

	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"api-url-shortener/internal/helpers"
)

func TestHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helpers Suite")
}

var _ = Describe("Helpers", func() {
	Describe("EnforceHTTP", func() {
		Context("with a URL without HTTP prefix", func() {
			It("should return a URL with 'http://' prefix", func() {
				url := "example.com"
				result := helpers.EnforceHTTP(url)
				Expect(result).To(Equal("http://example.com"))
			})
		})

		// Add more test cases as needed
	})

	Describe("RemoveDomainError", func() {
		Context("with a URL containing 'www.' and 'http://'", func() {
			It("should return a URL with 'www.' and 'http://' removed", func() {
				url := "http://www.example.com"
				valid, newURL := helpers.RemoveDomainError(url)
				Expect(valid).To(BeTrue())
				Expect(newURL).To(Equal("example.com"))
			})
		})

		// Add more test cases as needed
	})

	Describe("CheckURLAlreadyExists", func() {
		Context("with an existing URL in the Redis database", func() {
			It("should return true and the original URL", func() {

				mockClient := redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
					Password: "", 
					DB:       0,
				})
				mockData := map[string]string{"abc123": "http://example.com"}

				mockByteData,_ := json.Marshal(mockData)

				mockClient.Set("url_data", mockByteData, 0)

				exists, orgURL, err2 := helpers.CheckURLAlreadyExists(mockClient, "http://example.com", "url_data")
				Expect(err2).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
				Expect(orgURL).To(Equal("abc123"))
			})
		})

		// Add more test cases as needed
	})
})

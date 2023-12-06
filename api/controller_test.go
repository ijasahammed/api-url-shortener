package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"api-url-shortener/api"

)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("Controller", func() {
	var (
		mockRedisClient    *redis.Client
		mockShortenURLBody 	[]byte
		mockRepo			api.Repository
	)

	BeforeSuite(func() {

		mockRedisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		os.Setenv("SHORT_BASE_URL", "http://localhost:8080")

		mockShortenURLBody = []byte(`{"url":"http://example.com"}`)

		mockRepo = api.Repository{
			ShortUrlDBClient: mockRedisClient,
		}

	})

	BeforeEach(func() {
		mockRedisClient.FlushDB()
	})

	Describe("ShortenURL", func() {
		Context("with a valid request", func() {
			It("should return a 200 status code and a valid response body", func() {
				req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(mockShortenURLBody))
				Expect(err).NotTo(HaveOccurred())

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				mockRepo.ShortenURL(c)

				Expect(w.Code).To(Equal(http.StatusOK))

				var responseBody map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(responseBody["url"]).To(Equal("http://example.com"))
				Expect(responseBody["short_url"]).To(HavePrefix("http://localhost:8080/"))
			})
		})
	})

	Describe("GetHostCount", func() {
		Context("with valid data in the database", func() {
			It("should return a 200 status code and a valid response body", func() {

				req, err := http.NewRequest("GET", "/host-count", nil)
				Expect(err).NotTo(HaveOccurred())

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				mockRepo.GetHostCount(c)

				Expect(w.Code).To(Equal(http.StatusOK))

				var responseBody map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).NotTo(HaveOccurred())

			})
		})

	})

	Describe("ResolveURL", func() {
		Context("with a valid short URL", func() {
			It("should return a 301 status code and redirect to the original URL", func() {

				mockClient := redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
					Password: "", 
					DB:       0,
				})
				mockData := map[string]string{"abc123": "http://example.com"}

				mockByteData,_ := json.Marshal(mockData)

				mockClient.Set("url_data", mockByteData, 0)

				req := httptest.NewRequest("GET", "/abc123", nil)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				q := req.URL.Query()
				q.Add("url", "abc123")
				req.URL.RawQuery = q.Encode()
				c.Request = req

				mockRepo.ResolveURL(c)

				Expect(w.Code).To(Equal(http.StatusMovedPermanently))
				Expect(w.Header().Get("Location")).To(Equal("http://example.com"))
			})
		})

	})
})

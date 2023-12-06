// database_test.go

package database_test

import (
	"os"
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"api-url-shortener/database"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var _ = Describe("Database", func() {
	Describe("CreateClient", func() {
		Context("with valid configuration", func() {
			It("should return a non-nil Redis client", func() {
				// Mock environment variables
				os.Setenv("DB_ADDR", "localhost:6379")
				os.Setenv("DB_PASS", "")

				// Call the CreateClient function
				client := database.CreateClient(0)

				// Validate the returned client
				Expect(client).NotTo(BeNil())
				Expect(client.Options().Addr).To(Equal("localhost:6379"))
				Expect(client.Options().Password).To(Equal(""))
				Expect(client.Options().DB).To(Equal(0))
			})
		})

		// Add more test cases as needed
	})
})

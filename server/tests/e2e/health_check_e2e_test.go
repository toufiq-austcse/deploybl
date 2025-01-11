//go:build integration
// +build integration

package e2e_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/toufiq-austcse/deployit/internal/api/index/controller"
)

var _ = Describe("HealthCheckE2e", func() {
	ginEngine := gin.Default()
	healthCheckRouter := ginEngine.Group("/")
	var testingServer *httptest.Server

	BeforeEach(func() {
		testingServer = httptest.NewServer(ginEngine)
		healthCheckRouter.GET("", controller.Index())
	})
	AfterEach(func() {
		testingServer.Close()
	})

	Context("When health check is called", func() {
		It("should return 200", func() {
			resp, err := http.Get(testingServer.URL)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})

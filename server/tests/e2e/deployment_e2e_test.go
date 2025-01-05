package e2e_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/di"
	"github.com/toufiq-austcse/deployit/internal/app"
	"github.com/toufiq-austcse/deployit/internal/server"
)

var _ = Describe("DeploymentE2e", Ordered, func() {
	var testingServer *httptest.Server

	BeforeAll(func() {
		if err := config.Init(); err != nil {
			Fail("Error in initializing config")
		}
		mainServer := server.NewServer()
		testingServer = httptest.NewServer(mainServer.GinEngine)
		container, err := di.NewE2eDiContainer()
		if err != nil {
			Fail("Error in creating DI container")
		}
		err = app.SetupRouters(mainServer, container)
		if err != nil {
			fmt.Println(err.Error())
			Fail("Error in setting up routers ")
		}
	})

	AfterAll(func() {
		testingServer.Close()
	})

	When(" list deployments is called without Auth Token", func() {
		It("should return 401", func() {
			resp, err := http.Get(testingServer.URL + "/api/v1/deployments")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})

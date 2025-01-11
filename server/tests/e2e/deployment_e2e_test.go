//go:build integration
// +build integration

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/di"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
	"github.com/toufiq-austcse/deployit/internal/app"
	"github.com/toufiq-austcse/deployit/internal/server"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/firebase"
)

var _ = Describe("DeploymentE2e", Ordered, func() {
	var testingServer *httptest.Server
	var accessToken string

	BeforeAll(func() {
		if err := config.Init(); err != nil {
			Fail("Error in initializing config")
		}
		mainServer := server.NewServer()
		gin.SetMode(gin.TestMode)
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

		invokeErr := container.Invoke(func(deploymentController *controller.DeploymentController,

			firebaseClient *firebase.Client,
		) {
			customToken, err := firebaseClient.AuthClient.CustomToken(context.Background(), config.AppConfig.TEST_UID)
			if err != nil {
				Fail("Error in getting access token")
			}

			verifyRes, err := firebaseClient.VerifyCustomToken(customToken)
			if err != nil {
				Fail("Error in verifying custom token")
			}
			accessToken = verifyRes.IDToken
		})
		if invokeErr != nil {
			Fail("Error in invoking container")
		}
	})

	AfterAll(func() {
		testingServer.Close()
	})

	Describe("list deployments", func() {
		When("list deployments is called without Auth Token", func() {
			It("should return 401", func() {
				response, err := resty.New().SetBaseURL(testingServer.URL).R().Get("api/v1/deployments")
				Expect(err).To(BeNil())
				Expect(response.StatusCode()).To(Equal(http.StatusUnauthorized))
			})
		})
		When("list deployments is called with wrong Auth Token", func() {
			It("should return 401", func() {
				response, err := resty.
					New().
					SetBaseURL(testingServer.URL).
					R().
					SetHeader("Authorization", "wrong token").
					Get("api/v1/deployments")
				Expect(err).To(BeNil())
				Expect(response.StatusCode()).To(Equal(http.StatusUnauthorized))
			})
		})
		When("list deployments is called with valid Auth Token", func() {
			It("should return 200", func() {
				response, err := resty.
					New().
					SetBaseURL(testingServer.URL).
					R().
					SetHeader("Authorization", accessToken).
					Get("api/v1/deployments")
				Expect(err).To(BeNil())
				Expect(response.StatusCode()).To(Equal(http.StatusOK))

				var responseBody api_response.PaginationResponse
				err = json.Unmarshal(response.Body(), &responseBody)

				Expect(err).To(BeNil())
				Expect(reflect.TypeOf(responseBody.Data).Kind()).To(Equal(reflect.Slice))
				Expect(responseBody.Pagination).ShouldNot(BeNil())
			})
		})
	})
})

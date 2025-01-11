package deployment_controller_test_suite

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
	"github.com/toufiq-austcse/deployit/internal/api/users/model"
)

var _ = Describe("DeploymentController", func() {
	var (
		deploymentController *controller.DeploymentController
		router               *gin.Engine
		w                    *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.Use(func(context *gin.Context) {
			context.Set("user", &model.User{})
		})
		w = httptest.NewRecorder()
		deploymentController = controller.NewDeploymentController(nil, nil, nil, nil, nil, nil, nil)

		router.GET("/api/v1/deployments", deploymentController.DeploymentIndex)
	})

	Context("DeploymentIndex", func() {
		When("there are no deployments", func() {
			It("should return an empty array", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
	})
})

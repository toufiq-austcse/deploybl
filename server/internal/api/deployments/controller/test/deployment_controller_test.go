package deployment_controller_test_suite

import (
	"net/http"
	"net/http/httptest"

	mock_service "github.com/toufiq-austcse/deployit/internal/api/deployments/mocks/service"
	model2 "github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"go.uber.org/mock/gomock"

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

		mockDeploymentService *mock_service.MockIDeploymentService
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.Use(func(context *gin.Context) {
			context.Set("user", &model.User{})
		})
		w = httptest.NewRecorder()

		mockController := gomock.NewController(GinkgoT())
		mockDeploymentService = mock_service.NewMockIDeploymentService(mockController)

		deploymentController = controller.NewDeploymentController(nil, mockDeploymentService, nil, nil, nil, nil, nil)

		router.GET("/api/v1/deployments", deploymentController.DeploymentIndex)
	})

	Context("DeploymentIndex", func() {
		When("there are no deployments", func() {
			It("should return an empty array", func() {
				mockDeploymentService.EXPECT().
					ListDeployment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*model2.Deployment{}, nil, nil)
				req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
	})
})

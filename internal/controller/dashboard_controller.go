package controller

import (
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DashboardController struct {
	DashboardService *service.DashboardService
}

func NewDashboardController(dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{
		DashboardService: dashboardService,
	}
}

// GetFinancialOverviewHandler godoc
// @Summary 	Get financial overview
// @Description Get user's financial overview including current balance, monthly income, monthly expense, and monthly savings
// @Tags 		dashboard
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Success 	200 {object} response.ApiResponse{data=response.RespFinancialOverview}
// @Failure 	401 {object} response.ApiResponse
// @Failure 	500 {object} response.ApiResponse
// @Router 		/dashboard/overview [get]
func (c *DashboardController) GetFinancialOverviewHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		logrus.Errorf("Failed to get user ID from context: %v", err)
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	overview, err := c.DashboardService.GetFinancialOverview(userID)
	if err != nil {
		logrus.Errorf("Error getting financial overview: %v", err)
		utility.ServerErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get Financial Overview successful",
		Data:            overview,
	})
}

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
// @Success 	200 {object} response.SuccessResponse{data=response.RespFinancialOverview}
// @Failure 	401 {object} response.SuccessResponse
// @Failure 	500 {object} response.SuccessResponse
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
		utility.InternalServerErrorResponse(ctx, "Failed to get financial overview", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get Financial Overview successful",
		Data:            overview,
	})
}

// GetDashboardChartsHandler godoc
// @Summary 	Get dashboard charts data
// @Description Get user's dashboard charts including income vs expense, category distribution, and top expenses
// Tags 		dashboard
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Success 	200 {object} response.SuccessResponse{data=response.RespDashboardCharts}
// @Failure 	401 {object} response.ErrorResponse
// @Failure 	500 {object} response.ErrorResponse
// @Router 		/dashboard/charts [get]
func (c *DashboardController) GetDashboardChartsHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		logrus.Errorf("Failed to get user ID from context: %v", err)
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	charts, err := c.DashboardService.GetDashboardCharts(userID)
	if err != nil {
		logrus.Errorf("Error getting dashboard charts: %v", err)
		utility.InternalServerErrorResponse(ctx, "Failed to get dashboard charts", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get Dashboard Charts successful",
		Data:            charts,
	})
}

package router

import (
	"github.com/gin-gonic/gin"
	"go-manajemen-keuangan/internal/controller"
	"go-manajemen-keuangan/internal/payload/response"
	"go-manajemen-keuangan/internal/service"
	"gorm.io/gorm"
	"net/http"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	// init user service dan controller
	userService := &service.UserService{DB: db}
	userController := &controller.UserController{UserService: userService}

	// Healtcheck endpoint
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.ApiResponse{
			ResponseStatus:  true,
			ResponseMessage: "ok",
			Data:            nil,
		})
	})

	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", userController.RegisterHandler)
		userRouter.POST("/login", userController.LoginHandler)
	}
}

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

	// API routes group
	api := r.Group("/api")
	{
		api.GET("/health-check", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.ApiResponse{
				ResponseStatus:  true,
				ResponseMessage: "ok",
				Data:            nil,
			})
		})

		// user endpoint
		userRouter := api.Group("/user")
		{
			userRouter.POST("/register", userController.RegisterHandler)
			userRouter.POST("/login", userController.LoginHandler)
		}
	}

	// serve frontend static file
	r.Static("/js", "./web/dist/js")
	r.Static("/css", "./web/dist/css")
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// handle SPA routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}

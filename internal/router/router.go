package router

import (
	"go-fintrack/internal/controller"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/service"
	"go-fintrack/middleware"
	"net/http"
	"strings"

	_ "go-fintrack/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	// init user service dan controller
	userService := &service.UserService{DB: db}
	userController := &controller.UserController{UserService: userService}

	// init category
	categoryService := &service.CategoryService{DB: db}
	categoryController := &controller.CategoryController{CategoryService: categoryService}

	// init transaction
	transactionService := &service.TransactionService{DB: db}
	transactionController := &controller.TransactionController{TransactionService: transactionService}

	// swagger enpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes group
	api := r.Group("/api")
	{
		api.GET("/health-check", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, response.ApiResponse{
				ResponseStatus:  true,
				ResponseMessage: "ok",
				Data:            nil,
			})
		})

		// user endpoint
		userRouter := api.Group("/auth")
		{
			userRouter.POST("/register", userController.RegisterHandler)
			userRouter.POST("/login", userController.LoginHandler)
		}

		// admin endpoint
		adminRouter := api.Group("/admin")
		adminRouter.Use(middleware.Authentication(), middleware.AdminOnly())
		{
			//	router admin
		}

		// category endpoint
		categoryRouter := api.Group("/category")
		categoryRouter.Use(middleware.Authentication())
		{
			categoryRouter.GET("", categoryController.GetAllCategoriesHandler)
			categoryRouter.GET("/:id", categoryController.GetCategoryIdHandler)
			categoryRouter.POST("", categoryController.CreateCategoryHandler)
			categoryRouter.PUT("/:id", categoryController.UpdateCategoryHandler)
			categoryRouter.DELETE("/:id", categoryController.DeleteCategoryHandler)
		}

		// transaction endpoint
		transactionRouter := api.Group("/transaction")
		transactionRouter.Use(middleware.Authentication())
		{
			transactionRouter.GET("", transactionController.GetTransactionHandler)
			transactionRouter.POST("", transactionController.CreateTransactionHandler)
			transactionRouter.PUT("/:id", transactionController.UpdateTransactionHandler)
			transactionRouter.DELETE("/:id", transactionController.DeleteTransactionHandler)
			transactionRouter.GET("/export", transactionController.ExportTransactionsExcelHandler)
		}

		chatRouter := api.Group("/chat")
		chatRouter.Use(middleware.Authentication())
		{
			chatRouter.POST("/stream", controller.StreamChat)
		}
	}

	// serve frontend static file
	r.Static("/js", "./web/dist/js")
	r.Static("/css", "./web/dist/css")
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// handle SPA routing
	r.NoRoute(func(ctx *gin.Context) {
		// not found enpoint
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			ctx.JSON(http.StatusNotFound, response.ApiResponse{
				ResponseStatus:  false,
				ResponseMessage: "Endpoint not found",
				Data:            nil,
			})
			return
		}

		ctx.File("./web/dist/index.html")
	})
}

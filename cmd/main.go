package main

import (
	"go-fintrack/config"
	"go-fintrack/internal/router"
	"go-fintrack/internal/utility"
	"go-fintrack/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title 			Financial Management Application API
// @version 		1.0
// @description 	API for Financial Management Application
// @temsOfService 	http://swagger.io/terms/

// @contact.name 	Rafli Bima Pratandra
// @contact.email 	raflibima1106@gmail.com

// @license.name 	Apache 2.0
// @license.url 	http://www.apache.org/licenses/LICENSE-2.0.html

// @host 			localhost:8080
// @BasePath 		/api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect database
	db := config.ConnectDB()
	if db == nil {
		logrus.Fatal("Failed to connect database")
	}
	logrus.Info("Database connected!")

	// init google oauth config
	config.InitGoogleOauthConfig()

	// set gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin router
	r := gin.Default()
	r.Use(utility.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CorsMiddleware())

	// setup router
	router.InitRoutes(r, db)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// start HTTTP server
	logrus.Infof("Starting server on port :%s", serverPort)
	if err := r.Run(":" + serverPort); err != nil {
		logrus.Fatalf("HTTP server failed to start: %v", err)
	}
}

package main

import (
	"go-manajemen-keuangan/config"
	"go-manajemen-keuangan/internal/router"
	"go-manajemen-keuangan/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

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

	// set gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin router
	r := gin.Default()
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

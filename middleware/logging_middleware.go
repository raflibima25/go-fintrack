package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func setupLogger() *logrus.Logger {
	logger := logrus.New()

	if err := os.MkdirAll("logs", 0755); err != nil {
		logger.Fatal("Failed to create logs directory", err)
	}

	filename := filepath.Join("logs", fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("Failed to open log file", err)
	}

	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

var logger = setupLogger()

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		ctx.Set("RequestID", requestID)

		start := time.Now()

		ctx.Next()

		// Log duration
		duration := time.Since(start)
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()

		fields := logrus.Fields{
			"request_id": requestID,
			"client_ip":  clientIP,
			"duration":   duration,
			"method":     ctx.Request.Method,
			"path":       ctx.Request.URL.Path,
			"status":     statusCode,
			"user_agent": userAgent,
			"query":      ctx.Request.URL.RawQuery,
		}

		if err, exists := ctx.Get("Error"); exists {
			fields["error"] = err
		}

		if statusCode >= 500 {
			logger.WithFields(fields).Error("Internal server error")
		} else if statusCode >= 400 {
			logger.WithFields(fields).Warn("Client Error")
		} else {
			logger.WithFields(fields).Info("Request processed")
		}
	}
}

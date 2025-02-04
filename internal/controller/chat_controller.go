package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"go-fintrack/internal/utility"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func cleanResponse(text string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(text, "")

	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

func StreamChat(ctx *gin.Context) {
	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Set headers for SSE
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("X-Accel-Buffering", "no") // Disable nginx buffering if using nginx

	// Prepare Ollama request body
	ollamaBody := map[string]interface{}{
		"model":  "deepseek-r1:7b",
		"prompt": req.Message,
		"stream": true,
	}

	ollamaBodyBytes, err := json.Marshal(ollamaBody)
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}

	// Create Ollama request
	ollamaReq, err := http.NewRequest("POST", "http://localhost:11434/api/generate",
		bytes.NewBuffer(ollamaBodyBytes))
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}
	ollamaReq.Header.Set("Content-Type", "application/json")

	// Send request to Ollama
	ollamaResp, err := http.DefaultClient.Do(ollamaReq)
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}
	defer ollamaResp.Body.Close()

	// Create reader for response body
	reader := bufio.NewReader(ollamaResp.Body)

	// Stream response
	for {
		select {
		case <-ctx.Request.Context().Done():
			return // Client disconnected
		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				ctx.SSEvent("error", gin.H{"error": err.Error()})
				return
			}

			// Parse response
			var response struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}
			if err := json.Unmarshal(line, &response); err != nil {
				continue
			}

			// Send chunk to client
			if cleanedResponse := cleanResponse(response.Response); cleanedResponse != "" {
				ctx.SSEvent("message", response.Response)
				ctx.Writer.Flush()
			}

			// Check if generation is complete
			if response.Done {
				return
			}
		}
	}
}

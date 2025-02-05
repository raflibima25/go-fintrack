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
	"unicode"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func cleanResponse(text string) string {
	// Hapus tag XML tetapi pertahankan kontennya
	re := regexp.MustCompile(`<think>(.*?)</think>`)
	text = re.ReplaceAllString(text, "$1")

	re = regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")

	// Pisahkan kata-kata yang menempel
	var result []rune
	var lastRune rune
	for i, r := range text {
		if i > 0 && lastRune != ' ' {
			// Jika bertemu huruf kapital dan sebelumnya huruf kecil, tambahkan spasi
			if unicode.IsUpper(r) && unicode.IsLower(lastRune) {
				result = append(result, ' ')
			}
			// Jika bertemu huruf dan sebelumnya angka atau sebaliknya, tambahkan spasi
			if (unicode.IsLetter(r) && unicode.IsNumber(lastRune)) ||
				(unicode.IsNumber(r) && unicode.IsLetter(lastRune)) {
				result = append(result, ' ')
			}
		}
		result = append(result, r)
		lastRune = r
	}
	text = string(result)

	// Handle markdown formatting
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "**$1**")
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "*$1*")
	text = regexp.MustCompile(`\_(.*?)\_`).ReplaceAllString(text, "_$1_")

	// Handle lists
	text = regexp.MustCompile(`(?m)^(\d+)\.\s`).ReplaceAllString(text, "$1. ")
	text = regexp.MustCompile(`(?m)^[-*]\s`).ReplaceAllString(text, "• ")

	// Fix multiple spaces
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Fix newlines
	text = strings.ReplaceAll(text, `\n`, "\n")
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

func prosesChunk(chunk string) string {
	var processed []string
	current := ""

	for _, char := range chunk {
		if current == "" {
			current += string(char)
			continue
		}

		lastChar := rune(current[len(current)-1])
		if unicode.IsLower(lastChar) && unicode.IsUpper(char) {
			processed = append(processed, current)
			current = string(char)
		} else {
			current += string(char)
		}
	}

	if current != "" {
		processed = append(processed, current)
	}

	return strings.Join(processed, current)
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
	ctx.Header("X-Accel-Buffering", "no")

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
	var responseBuffer strings.Builder

	// Stream response
	for {
		select {
		case <-ctx.Request.Context().Done():
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				if responseBuffer.Len() > 0 {
					processedText := cleanResponse(responseBuffer.String())
					if processedText != "" {
						ctx.SSEvent("message", processedText)
						ctx.Writer.Flush()
					}
				}
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

			if response.Response != "" {
				processed := prosesChunk(response.Response)
				responseBuffer.WriteString(processed)

				if strings.ContainsAny(response.Response, ".!?\n") {
					processedText := cleanResponse(responseBuffer.String())
					if processedText != "" {
						ctx.SSEvent("message", processedText)
						ctx.Writer.Flush()
						responseBuffer.Reset()
					}
				}
			}

			// Check if generation is complete
			if response.Done {
				if responseBuffer.Len() > 0 {
					processedText := cleanResponse(responseBuffer.String())
					if processedText != "" {
						ctx.SSEvent("message", processedText)
						ctx.Writer.Flush()
					}
				}
				return
			}
		}
	}
}

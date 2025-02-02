package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
)

type Response struct {
	APIResponse string `json:"api_response"`
}

func main() {
	http.HandleFunc("/", HandleHome())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func HandleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == "GET" {
			http.ServeFile(w, r, "index.html")
			return
		}

		if r.Method == "POST" {
			// Parse the multipart form
			err := r.ParseMultipartForm(10 << 20) // 10 MB max
			if err != nil {
				http.Error(w, "Unable to parse form", http.StatusBadRequest)
				return
			}

			// Get the file from form
			file, _, err := r.FormFile("image")
			if err != nil {
				http.Error(w, "Error retrieving file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Read the file content
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}

			// Convert to base64
			base64String := base64.StdEncoding.EncodeToString(fileBytes)
			context := r.FormValue("context-ai")

			// Here you would make your API call and get the response
			apiResponse := APICall(base64String, context) // Replace with actual API call

			response := Response{
				APIResponse: apiResponse,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}
}

func APICall(image string, context string) string {
	client := &http.Client{}

	// OpenAI API endpoint
	url := "https://api.openai.com/v1/chat/completions"

	// Prepare request body
	requestBody := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": context,
					},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": "data:image/jpeg;base64," + image,
						},
					},
				},
			},
		},
		"max_tokens": 500,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "Error preparing request: " + err.Error()
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "Error creating request: " + err.Error()
	}

	// Add headers
	key := getAPIKey()
	if key == "-1" {
		return "Error getting API key"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return "Error making request: " + err.Error()
	}
	defer resp.Body.Close()

	// Read response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Error reading response: " + err.Error()
	}

	fmt.Println(result)
	// Extract response text
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content
				}
			}
		}
	}

	return "Error processing response"
}

func getAPIKey() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Try to read from .env file
		content, err := os.ReadFile(".env")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "OPENAI_API_KEY=") {
					apiKey = strings.TrimPrefix(line, "OPENAI_API_KEY=")
					apiKey = strings.TrimSpace(apiKey)
					break
				}
			}
		}
	}
	if apiKey == "" {
		log.Println("OPENAI_API_KEY is required in environment or .env file")
		return "-1"
	}
	return apiKey
}

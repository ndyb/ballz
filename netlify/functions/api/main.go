package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is the structure for our API response
type Response struct {
	Scores interface{} `json:"scores,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// Score represents a player score
type Score struct {
	CreatedAt int64 `json:"created_at"`
	Score     int   `json:"score"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get API key from environment
	apiKey := os.Getenv("API_KEY")

	// Connect to the database
	db, err := Connect()
	if err != nil {
		return errorResponse(500, "Failed to connect to database: "+err.Error())
	}
	defer db.Close()

	// Route based on path and method
	path := request.Path
	method := request.HTTPMethod

	// Get scores - public endpoint
	if path == "/scores" && method == "GET" {
		scores, err := db.GetHighscores(10) // Pass db connection to GetHighscores
		if err != nil {
			return errorResponse(500, "Failed to retrieve scores: "+err.Error())
		}
		return successResponse(scores)
	}

	// Insert score - protected endpoint
	if path == "/scores" && method == "POST" {
		// Verify API key
		authHeader := request.Headers["apikey"]
		if authHeader != apiKey || apiKey == "" {
			return errorResponse(401, "Unauthorized")
		}

		// Parse the score from the request body
		var score Score
		if err := json.Unmarshal([]byte(request.Body), &score); err != nil {
			return errorResponse(400, "Invalid score data")
		}

		// Insert the score
		if err := db.InsertHighscore(score.Score); err != nil { // Pass db connection to InsertHighscore
			return errorResponse(500, "Failed to insert score: "+err.Error())
		}

		return successResponse(map[string]string{"message": "Score added successfully"})
	}

	// If no route matches
	if path != "/scores" {
		return successResponse(42)
	}

	return errorResponse(404, "Not Found")
}

func successResponse(data interface{}) (events.APIGatewayProxyResponse, error) {
	resp := Response{
		Scores: data,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return errorResponse(500, "Error generating response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

func errorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	resp := Response{
		Error: message,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		body = []byte(`{"error":"Error generating response"}`)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	// Start the Lambda handler
	lambda.Start(handler)
}

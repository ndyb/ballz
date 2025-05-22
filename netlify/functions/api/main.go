package main

import (
	"encoding/json"
	"log"
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

	log.Println("Request path:", path)
	log.Println("Request method:", method)
	log.Println("Request body:", request.Body)

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
		log.Println("No route matches, returning asdf")

		// Sample data for testing purposes
		sampleData :=
			`[
  {
    "score": 393872,
    "created": "2025-05-22 09:37:34.987437+00"
  },
  {
    "score": 65,
    "created": "2025-05-22 09:37:25.175401+00"
  },
  {
    "score": 54,
    "created": "2025-05-22 09:37:19.346369+00"
  },
  {
    "score": 43,
    "created": "2025-05-22 09:37:10.335945+00"
  },
  {
    "score": 5,
    "created": "2025-05-21 18:26:59.365226+00"
  },
  {
    "score": 3,
    "created": "2025-05-21 17:20:08.739607+00"
  },
  {
    "score": 2,
    "created": "2025-05-21 16:58:29.489417+00"
  },
  {
    "score": 1,
    "created": "2025-05-22 09:37:43.338297+00"
  },
  {
    "score": 1,
    "created": "2025-05-21 16:47:34.136628+00"
  }
]`

		return successResponse(sampleData)
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

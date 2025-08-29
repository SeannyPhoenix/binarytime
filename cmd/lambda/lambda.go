package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/seannyphoenix/binarytime/pkg/binarytime"
)

type Response struct {
	BinaryTime string `json:"binaryTime"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get current binarytime
	now := binarytime.Now()

	// Get format from query parameters, default to "hex"
	format := request.QueryStringParameters["format"]
	if format != "base64" {
		format = "hex"
	}

	// Format the binarytime according to the requested format
	var timeStr string
	if format == "base64" {
		timeStr = now.Fixed128().Base64()
	} else {
		timeStr = now.Fixed128().String()
	}

	// Create response body
	response := Response{
		BinaryTime: timeStr,
	}

	// Convert response to JSON
	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"Failed to generate response"}`,
		}, nil
	}

	// Return successful response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}

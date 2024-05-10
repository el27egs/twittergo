package models

import "github.com/aws/aws-lambda-go/events"

type ApiResponse struct {
	Status         int
	Body           string
	ActualResponse *events.APIGatewayProxyResponse
}

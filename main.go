package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/starlingapps/twittergo/awsgo"
	"github.com/starlingapps/twittergo/db"
	"github.com/starlingapps/twittergo/handlers"
	"log"
	"os"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.LoadDefaultConfig()

	envVars := []string{"SECRET_ID", "BUCKET_NAME"}
	if !checkMandatoryEnvVars(envVars) {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Please, set minimum required environment variables: %v", envVars),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}

	settings, err := awsgo.GetSecret(os.Getenv("SECRET_ID"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error on reading settings: %s\n", err),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}
	fmt.Printf("> Secret has:\n%+v\n", settings)
	err = db.Connect(settings)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error on connecting to database: %s", err),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	}
	response := handlers.ApiHandler(settings.JwtSeed, request)
	if response.ActualResponse == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: response.Status,
			Body:       response.Body,
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return res, nil
	} else {
		return response.ActualResponse, nil
	}
}

func checkMandatoryEnvVars(envVars []string) bool {
	for _, envVar := range envVars {
		_, exists := os.LookupEnv(envVar)
		if !exists {
			log.Printf("%s is not set", envVar)
			return false
		}
	}
	return true
}

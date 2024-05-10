package handlers

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/starlingapps/twittergo/models"
	"github.com/starlingapps/twittergo/routers"
	"github.com/starlingapps/twittergo/security"
)

func ApiHandler(seed string, request events.APIGatewayProxyRequest) models.ApiResponse {
	var res = models.ApiResponse{
		Status: 404,
		Body:   "Invalid call",
	}
	httpMethod := request.HTTPMethod
	// fmt.Printf("> HttpRequest %+v\n", request)
	path := request.PathParameters["twittergo"] // URL_PREFIX was not required here
	fmt.Printf("> Processing HTTP call: %s %s\n", httpMethod, path)
	validToken, status, message, claims := validateToken(path, seed, request)
	fmt.Printf("> Claims after validateToken call %v\n", claims)
	if !validToken {
		res.Status = status
		res.Body = message
		return res
	}
	switch httpMethod {
	case "POST":
		switch path {
		case "signup":
			return routers.SingUp(request)
		case "login":
			return routers.Login(seed, request)
		case "tweets":
			return routers.CreateTweet(request, claims.ID.Hex())
		case "images":
			return routers.UploadImage(request, claims.ID.Hex())
		}
	case "GET":
		switch path {
		case "users":
			return routers.GetUser(request)
		case "tweets":
			return routers.GetTweets(request)
		case "images":
			return routers.DownloadImage(request, claims.ID.Hex())
		}
	case "PUT":
		switch path {
		case "users":
			return routers.UpdateUser(request, claims.ID.Hex())

		}
	case "DELETE":
		switch path {
		case "tweets":
			return routers.DeleteTweet(request, claims.ID.Hex())

		}
	}
	return res
}

func validateToken(path, seed string, req events.APIGatewayProxyRequest) (bool, int, string, models.Claims) {
	if path == "signup" || path == "login" || path == "avatar" || path == "banner" {
		return true, 200, "", models.Claims{}
	}
	token := req.Headers["Authorization"]

	claims, valid, msg, err := security.ProcessJwtToken(token, seed)
	if !valid {
		if err != nil {
			fmt.Printf("> Error on validating bearer token: %s", err)
			return false, 401, err.Error(), models.Claims{}
		} else {
			fmt.Printf("> Error on validating bearer token on server: %s", msg)
			return false, 401, msg, models.Claims{}
		}
	}
	return true, 200, "", claims
}

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
	fmt.Printf("> HttpRequest %+v\n", request)
	path := request.PathParameters["twittergo"] // URL_PREFIX was not required here
	fmt.Printf("> Processing HTTP call: %s %s\n", httpMethod, path)
	fmt.Printf("> Processing HTTP call: %s %s\n", httpMethod, request.Path)

	validToken, status, message, _ := validateToken(path, seed, request)

	if !validToken {
		res.Status = status
		res.Body = message
		return res
	}

	switch httpMethod {
	case "POST":
		switch path {
		case "singup":
			return routers.SingUp(request)
		}
	case "GET":
		switch path {

		}
	case "PUT":
		switch path {

		}
	case "DELETE":
		switch path {

		}
	}
	return res
}

func validateToken(path, seed string, req events.APIGatewayProxyRequest) (bool, int, string, models.Claims) {
	if path == "singup" || path == "login" || path == "avatar" || path == "banner" {
		return true, 200, "", models.Claims{}
	}
	token := req.Headers["Authorization"]

	claim, valid, msg, err := security.ProcessJwtToken(token, seed)
	if !valid {
		if err != nil {
			fmt.Printf("> Error al validar token: %s", err)
			return false, 401, err.Error(), models.Claims{}
		} else {
			fmt.Printf("> Error al validar token: %s", msg)
			return false, 401, msg, models.Claims{}
		}
	}
	return true, 200, "", *claim
}

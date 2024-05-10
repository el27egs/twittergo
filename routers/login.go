package routers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/starlingapps/twittergo/db"
	"github.com/starlingapps/twittergo/models"
	"github.com/starlingapps/twittergo/security"
	"net/http"
	"time"
)

func Login(seed string, req events.APIGatewayProxyRequest) models.ApiResponse {
	var payload models.User
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into Login method\n")
	bodyReq := req.Body
	err := json.Unmarshal([]byte(bodyReq), &payload)
	if err != nil {
		res.Body = "Input login data is not correct"
		fmt.Printf("> Error Unmarshall %s", res.Body)
		return res
	}
	if len(payload.Email) == 0 {
		res.Body = "Email is mandatory"
		fmt.Printf("> Error validating email %s", res.Body)
		return res
	}
	user, found := db.LoginUser(payload.Email, payload.Password)
	if !found {
		res.Status = 404
		res.Body = "User not found"
		fmt.Printf("> User with email[%s] was not found on database\n", payload.Email)
		return res
	}
	token, err := security.CreateNewJwt(seed, user)
	if err != nil {
		res.Status = 500
		res.Body = "Error on generating token for authenticating user"
		fmt.Printf("> Error on CreateNewJt %s", err.Error())
		return res
	}
	loginResponse := models.LoginResponse{
		Token: token,
	}
	bodyRes, err2 := json.Marshal(loginResponse)
	if err2 != nil {
		res.Status = 500
		res.Body = "Error on formatting token for authenticated user"
		fmt.Printf("> Error on Marshal loginResponse object %s", err2.Error())
		return res
	}
	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	actualResponse := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bodyRes),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-type":                "application/json",
			"Set-Cookie":                  cookie.String(),
		},
	}
	res.Status = 200
	res.Body = string(bodyRes)
	res.ActualResponse = actualResponse
	fmt.Printf("> Returning login respose %+v\n", res)
	return res
}

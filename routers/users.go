package routers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/starlingapps/twittergo/db"
	"github.com/starlingapps/twittergo/models"
)

func GetUser(req events.APIGatewayProxyRequest) models.ApiResponse {
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into GEtProfile method\n")
	userId := req.QueryStringParameters["id"]
	if len(userId) == 0 {
		res.Body = "Parameter 'id' is mandatory"
		return res
	}
	user, err := db.FindUserById(userId)
	if err != nil {
		res.Status = 404
		res.Body = "Error on finding user profile"
		fmt.Printf("> Error on finding user profile %s\n", err)
		return res
	}
	bodyRes, err2 := json.Marshal(user)
	if err2 != nil {
		res.Status = 500
		res.Body = "Error on formatting user profile"
		fmt.Printf("> Error on Marshal user profile object %s", err2.Error())
		return res
	}
	res.Status = 200
	res.Body = string(bodyRes)
	fmt.Printf("> Returning user profile respose %+v\n", res)
	return res
}

func UpdateUser(req events.APIGatewayProxyRequest, userId string) models.ApiResponse {
	var payload models.User
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into UpdateProfile method for userId %s\n", userId)
	bodyReq := req.Body
	err := json.Unmarshal([]byte(bodyReq), &payload)
	if err != nil {
		res.Body = "Input profile data is not correct"
		fmt.Printf("> Error Unmarshall %s", res.Body)
		return res
	}
	status, err := db.UpdateProfileById(userId, payload)
	if err != nil || !status {
		res.Status = 404
		res.Body = "Error on updating user profile"
		fmt.Printf("> Error on updating user profile %s\n", err)
		return res
	}
	res.Status = 200
	res.Body = "User Profile updated"
	return res
}

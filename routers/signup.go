package routers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/starlingapps/twittergo/db"
	"github.com/starlingapps/twittergo/models"
)

func SingUp(req events.APIGatewayProxyRequest) models.ApiResponse {

	var payload models.User
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into SingUp method\n")
	body := req.Body
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		res.Body = err.Error()
		fmt.Printf("> Error Unmarshall %s", res.Body)
		return res
	}
	if len(payload.Email) == 0 {
		res.Body = "Email is mandatory"
		fmt.Printf("> Error validating email %s", res.Body)
		return res
	}
	if len(payload.Password) < 6 {
		res.Body = "Password must be at least 6 characters"
		fmt.Printf("> Error validating password %s", res.Body)
		return res
	}
	_, found, _ := db.FindUserByEmail(payload.Email)
	if found {
		res.Body = "Email is already registered"
		fmt.Printf("> Error validating password %s", res.Body)
		return res
	}
	_, status, err := db.SaveUser(payload)
	if err != nil || !status {
		res.Status = 500
		res.Body = fmt.Sprintf("> Error on saving data on database due %s", err.Error())
		fmt.Printf("> Error saving user %s", res.Body)
		return res
	}
	res.Status = 200
	res.Body = "New User saved successfully"
	fmt.Printf("> %s", res.Body)
	return res
}

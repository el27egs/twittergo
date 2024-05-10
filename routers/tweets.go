package routers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/starlingapps/twittergo/db"
	"github.com/starlingapps/twittergo/models"
	"strconv"
	"time"
)

func CreateTweet(req events.APIGatewayProxyRequest, userId string) models.ApiResponse {
	var payload models.TweetRequest
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into CreateTweet method for userId %s\n", userId)
	bodyReq := req.Body
	err := json.Unmarshal([]byte(bodyReq), &payload)
	if err != nil {
		res.Body = "Input tweet data is not correct"
		fmt.Printf("> Error Unmarshall %s", res.Body)
		return res
	}
	record := models.TweetDbModel{
		UserId:    userId,
		Message:   payload.Message,
		CreatedAt: time.Now(),
	}
	_, status, err := db.CreateTweet(record)
	if err != nil || !status {
		res.Status = 404
		res.Body = "Error on creating tweet"
		fmt.Printf("> Error on creating tweet %s\n", err)
		return res
	}
	res.Status = 200
	res.Body = "Tweet created successfully"
	return res
}

func GetTweets(req events.APIGatewayProxyRequest) models.ApiResponse {
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into GetTweets method\n")
	userId := req.QueryStringParameters["userId"]
	paramPage := req.QueryStringParameters["page"]
	if len(userId) == 0 {
		res.Body = "Parameter 'userId' is mandatory"
		return res
	}
	if len(paramPage) == 0 {
		paramPage = "1"
	}
	page, err := strconv.Atoi(paramPage)
	if err != nil {
		res.Status = 404
		res.Body = "Page must be a valid number"
		fmt.Printf("> Page must be a valid number greater or equal to 1 %s\n", err)
		return res
	}
	tweets, err := db.GetTweets(userId, int64(page))
	if err != nil {
		res.Status = 500
		res.Body = "Error on reading tweets"
		fmt.Printf("> Error on reading tweets from database%s\n", err)
		return res
	}
	bodyRes, err := json.Marshal(tweets)
	if err != nil {
		res.Status = 500
		res.Body = "Error on parsing tweets"
		fmt.Printf("> Error on parsing tweets to JSON format %s\n", err)
		return res
	}
	fmt.Println("> Reading tweets was successful")
	res.Status = 200
	res.Body = string(bodyRes)
	return res
}

func DeleteTweet(req events.APIGatewayProxyRequest, userId string) models.ApiResponse {
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into DeleteTweet method\n")
	tweetId := req.QueryStringParameters["tweetId"]
	if len(tweetId) == 0 {
		res.Body = "Parameter 'tweetId' is mandatory"
		return res
	}
	err := db.DeleteTweet(tweetId, userId)
	if err != nil {
		res.Status = 500
		res.Body = "Error on deleting tweet"
		fmt.Printf("> Error on deleting tweet%s\n", err)
		return res
	}
	fmt.Println("> Deleting tweets was successful")
	res.Status = 200
	res.Body = "Tweet deleted"
	return res
}

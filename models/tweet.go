package models

import "time"

type TweetRequest struct {
	Message string `bson:"message" json:"message"`
}

type TweetDbModel struct {
	UserId    string    `bson:"userId" json:"userId"`
	Message   string    `bson:"message" json:"message"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

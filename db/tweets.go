package db

import (
	"context"
	"github.com/starlingapps/twittergo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTweet(tweet models.TweetDbModel) (string, bool, error) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("tweets")
	result, err := col.InsertOne(ctx, tweet)
	if err != nil {
		return "", false, err
	}
	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

func GetTweets(userId string, page int64) ([]models.TweetDbModel, error) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("tweets")
	filter := bson.M{"userId": userId}
	tweetsPerPage := int64(20)
	opts := options.Find()
	opts.SetLimit(tweetsPerPage)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	opts.SetSkip((page - 1) * tweetsPerPage)

	cur, err := col.Find(ctx, filter, opts)
	var tweets []models.TweetDbModel
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var tweet models.TweetDbModel
		// Note: Here the unmarshalling is between mongo record to golang struct
		// that is the reason we are using cur.Decode method instead of
		//json.Unmarshal, the latter is used to convert from json string to golang struct
		err := cur.Decode(&tweet)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

func DeleteTweet(tweetId, userId string) error {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("tweets")
	// Note, Here _id is of type ObjectID, userId is only a string into tweets collection
	// therefore for userId we do not need the conversion to ObjectID
	objID, err := primitive.ObjectIDFromHex(tweetId)
	filter := bson.M{
		"_id":    objID,
		"userId": userId,
	}
	_, err = col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

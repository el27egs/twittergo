package db

import (
	"context"
	"github.com/starlingapps/twittergo/models"
	"github.com/starlingapps/twittergo/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUserByEmail(email string) (models.User, bool, string) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("users")
	condition := bson.M{"email": email}
	var user models.User
	err := col.FindOne(ctx, condition).Decode(&user)
	ID := user.ID.Hex()
	if err != nil {
		return user, false, ID
	}
	return user, true, ID
}

func SaveUser(user models.User) (string, bool, error) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("users")
	user.Password, _ = security.EncryptPassword(user.Password)
	result, err := col.InsertOne(ctx, user)
	if err != nil {
		return "", false, err
	}
	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

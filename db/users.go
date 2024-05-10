package db

import (
	"context"
	"fmt"
	"github.com/starlingapps/twittergo/models"
	"github.com/starlingapps/twittergo/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func LoginUser(email, password string) (models.User, bool) {
	user, found, _ := FindUserByEmail(email)
	if !found {
		fmt.Printf("> User with email [%s] was not found\n", email)
		return models.User{}, false
	}
	fmt.Printf("> User with email [%s] was found\n", email)

	valid := security.ValidatePasswords([]byte(user.Password), []byte(password))
	if !valid {
		fmt.Printf("> Error passwords mismatch [%s],[%s]\n", user.Password, password)
		return models.User{}, false
	}
	fmt.Printf("> Password validation was successful\n")

	fmt.Printf("> Login was successful\n")
	return user, true
}

func FindUserByEmail(email string) (models.User, bool, string) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("users")
	condition := bson.M{"email": email}
	var user models.User
	err := col.FindOne(ctx, condition).Decode(&user)
	id := user.ID.Hex()
	if err != nil {
		return user, false, id
	}
	return user, true, id
}

func FindUserById(id string) (models.User, error) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("users")
	mongoId, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{"_id": mongoId}
	var user models.User
	err := col.FindOne(ctx, condition).Decode(&user)
	user.Password = ""
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdateProfileById(id string, user models.User) (bool, error) {
	ctx := context.TODO()
	db := MongoClient.Database(DbName)
	col := db.Collection("users")
	data := make(map[string]any)
	if len(user.Name) > 0 {
		data["name"] = user.Name
	}
	if len(user.LastName) > 0 {
		data["lastName"] = user.LastName
	}
	data["DOB"] = user.DOB
	if len(user.Avatar) > 0 {
		data["avatar"] = user.Avatar
	}
	if len(user.Banner) > 0 {
		data["banner"] = user.Banner
	}
	if len(user.Bio) > 0 {
		data["bio"] = user.Bio
	}
	if len(user.Location) > 0 {
		data["location"] = user.Location
	}
	if len(user.Web) > 0 {
		data["web"] = user.Web
	}
	// We could use user instead of data in the below line
	// however, all empty fields like password or email un user instance
	// will be clean in Mongo due user instance has no values on those fields
	// so, we can conclude only all fields passed into the updateOne method are replaced
	// with or without a value based on the values sent as input.
	update := bson.M{
		"$set": data,
	}
	mongoId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": mongoId}
	_, err := col.UpdateOne(ctx, filter, update)

	fmt.Printf("> Filter %+v\n", filter)
	fmt.Printf("> Data to upate %+v\n", update)

	if err != nil {
		fmt.Printf("> Error on updating user data %s\n", err)
		return false, err
	}
	return true, nil
}

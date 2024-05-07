package db

import (
	"context"
	"fmt"
	"github.com/starlingapps/twittergo/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	DbName      string
)

func Connect(settings *models.Settings) error {
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/rettyWrites=true&w=majority", settings.DbUsername, settings.DbPassword, settings.DbHost)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connStr))
	if err != nil {
		fmt.Printf("Error Connect %s\n", err)
		return err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error Ping %s\n", err)
		return err
	}
	fmt.Println("> Conexion existosa con la base de datos")
	MongoClient = client
	DbName = settings.DbName
	return nil
}

func IsConnectionAlive() bool {
	err := MongoClient.Ping(context.TODO(), nil)
	return err == nil
}

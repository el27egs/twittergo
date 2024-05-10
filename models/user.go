package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name,omitempty"`
	LastName string             `bson:"lastName" json:"lastName,omitempty"`
	DOB      time.Time          `bson:"DOB" json:"DOB,omitempty"`
	Email    string             `bson:"email" json:"email,omitempty"`
	Password string             `bson:"password" json:"password,omitempty"`
	Avatar   string             `bson:"avatar" json:"avatar,omitempty"`
	Banner   string             `bson:"banner" json:"banner,omitempty"`
	Bio      string             `bson:"bio" json:"bio,omitempty"`
	Location string             `bson:"location" json:"location,omitempty"`
	Web      string             `bson:"web" json:"web,omitempty"`
}

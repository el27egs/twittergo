package models

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	Email string             `json:"email"`
	ID    primitive.ObjectID `bson:"ID" json:"ID,omitempty"`
	jwt.RegisteredClaims
}

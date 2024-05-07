package security

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/starlingapps/twittergo/models"
	"strings"
)

var Email string
var IdUsuario string

func ProcessJwtToken(bearerToken, seed string) (*models.Claims, bool, string, error) {
	var claims models.Claims
	bearerToken = strings.Replace(bearerToken, "Bearer", "", -1)
	bearerToken = strings.TrimSpace(bearerToken)
	if len(bearerToken) == 0 {
		return &claims, false, string(""), errors.New("Invalid Bearer Token")
	}
	token, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(seed), nil
	})
	if err == nil {
		//Rutina que checke contra la BD
	}
	if !token.Valid {
		return &claims, false, string(""), errors.New("Invalid JWT Token")
	}
	return &claims, true, string(""), nil
}

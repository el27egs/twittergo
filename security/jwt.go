package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/starlingapps/twittergo/models"
	"strings"
	"time"
)

func ProcessJwtToken(bearerToken, seed string) (models.Claims, bool, string, error) {
	var claims models.Claims
	bearerToken = strings.Replace(bearerToken, "Bearer", "", -1)
	bearerToken = strings.TrimSpace(bearerToken)
	if len(bearerToken) == 0 {
		return claims, false, string(""), errors.New("Invalid Bearer Token")
	}
	token, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(seed), nil
	})
	if err != nil {
		return claims, false, string(""), err
	}
	if !token.Valid {
		return claims, false, string(""), errors.New("Invalid JWT Token")
	}
	fmt.Printf("> Claims tiene::: \n\n%+v\n", claims)
	return claims, true, string(""), nil
}

func CreateNewJwt(seed string, user models.User) (string, error) {
	payload := jwt.MapClaims{
		"email":    user.Email,
		"name":     user.Name,
		"lastName": user.LastName,
		"DOB":      user.DOB,
		"bio":      user.Bio,
		"location": user.Location,
		"avatar":   user.Avatar,
		"banner":   user.Banner,
		"web":      user.Web,
		"ID":       user.ID.Hex(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(seed))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

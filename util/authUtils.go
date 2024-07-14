package util

import (
	"log"
	"synapse/auth"
	"github.com/golang-jwt/jwt/v5"
)

var(
	secret_key = []byte("Synapse_Rocks")
)


func GetUserNameFromToken(tokenStr string) string{
	token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(t *jwt.Token) (interface{}, error) { return secret_key, nil })
	if err != nil {
		log.Println("Error in parsing Token to extract userName")
		return ""
	}

	claim := token.Claims.(*auth.Claims)
	return claim.Id
}
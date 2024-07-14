package util

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

var(
	secret_key = []byte("Synapse_Rocks")
)

type claims struct {
	Id string
	*jwt.RegisteredClaims
}

func GetUserNameFromToken(tokenStr string) string{
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) { return secret_key, nil })
	if err != nil {
		log.Println("Error in parsing Token to extract userName")
		return ""
	}

	claim := token.Claims.(*claims)
	return claim.Id
}
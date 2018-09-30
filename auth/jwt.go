package auth

import (
	"log"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWT represents a generated json web token
type JWT struct{ Secret string }

// Encode a Mobius auth token as a JSON Web token
func (j *JWT) Encode(t *Token, options map[string]interface{}) string {
	var claims map[string]interface{}
	if options != nil {
		for i, v := range options {
			claims[i] = v
		}
	}
	iat, _ := strconv.ParseInt(string(t.TX.Tx.TimeBounds.MinTime), 10, 32)
	exp, _ := strconv.ParseInt(string(t.TX.Tx.TimeBounds.MaxTime), 10, 32)
	claims["jti"] = t.Hash("hex")
	claims["sub"] = t.Address
	claims["iat"] = iat
	claims["exp"] = exp
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims.(jwt.MapClaims))
	tokenStr, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		log.Fatalf("failed to sign jwt, err: %v", err)
		os.Exit(1)
	}
}

// Decode and verify a JSON Web Token
func (j *JWT) Decode(payload string) *jwt.Token {
	token, err := jwt.ParseWithClaims(payload, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		log.Fatalf("failed to parse jwt, err: %v", err)
	}
	return token
}

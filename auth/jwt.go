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
	iat, _ := strconv.ParseInt(strconv.Itoa(int(t.TX.Tx.Tx.TimeBounds.MinTime)), 10, 64)
	exp, _ := strconv.ParseInt(strconv.Itoa(int(t.TX.Tx.Tx.TimeBounds.MaxTime)), 10, 64)
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = t.Hash("hex")
	claims["sub"] = t.GetAddress()
	claims["iat"] = iat
	claims["exp"] = exp
	if options != nil {
		for i, v := range options {
			claims[i] = v
		}
	}
	tokenStr, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		log.Fatalf("failed to sign jwt, err: %v", err)
		os.Exit(1)
	}
	return tokenStr
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

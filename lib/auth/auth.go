package auth

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	currentKey    *rsa.PrivateKey
	oldKey        *rsa.PrivateKey
	tokenLifeTime time.Duration
)

func Init(current *rsa.PrivateKey, old *rsa.PrivateKey, lifeTime time.Duration) {
	currentKey = current
	oldKey = old
	tokenLifeTime = lifeTime
}

func Sign(userID string) (authToken string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(tokenLifeTime).Unix(),
	})
	return token.SignedString(currentKey)
}

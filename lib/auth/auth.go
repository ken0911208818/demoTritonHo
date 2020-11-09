package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
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

type Token struct {
	UserId string
	Exp    int64
	jwt.StandardClaims
}

func Sign(userID string) (authToken string, err error) {
	tk := Token{
		userID,
		time.Now().Add(tokenLifeTime).Unix(),
		jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, tk)
	ss, err := token.SignedString(currentKey)
	return ss, err
}

func Verify(authToken string) (userId string, err error) {
	// parse and vertify the token string
	//token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
	//	// make sure the JWT token is using RSA alg
	//	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
	//		return nil, errors.New("Unexpected signing method")
	//	}
	//	return &currentKey.PublicKey, nil
	//})
	token, err := jwt.ParseWithClaims(authToken, &Token{}, func(t *jwt.Token) (interface{}, error) {
		// make sure the JWT token is using RSA alg
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return &currentKey.PublicKey, nil
	})
	if err != nil {
		return ``, err
	}

	if token.Valid == false { // make sure token is Valid
		return ``, errors.New("Wrong jwt token")
	}

	if s, ok := token.Claims.(*Token); !ok {
		return ``, errors.New("Improper JWT Token")
	} else {
		userId = s.UserId
	}
	fmt.Println(userId)
	return userId, nil
}

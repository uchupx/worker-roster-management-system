package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWT(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j *JWT) Create(ttl time.Duration, content interface{}) (result *string, err error) {
	now := time.Now()

	claims := make(jwt.MapClaims)
	claims["dat"] = content
	claims["exp"] = now.Add(ttl).Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		return
	}

	result = &token

	return
}

func (j *JWT) Validate(token string) (resp interface{}, err error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, err
		}

		return j.publicKey, nil
	})

	if err != nil {
		return
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return
	}

	resp = claims["dat"]

	return
}

func (j *JWT) Content(token string) (map[string]interface{}, error) {
	jwtToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}

		return j.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return jwtToken.Claims.(jwt.MapClaims), nil
}

package lib

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Id int
	Email string
	jwt.StandardClaims
}

type AccessKey struct {
	Token string `json:"accessKey"`
}

func (a *AccessKey) ToJson() ([]byte, error){

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Claims) CreateToken(key []byte) (*AccessKey, error){

	expirationTime := time.Now().Add(30 * time.Minute)

	if c.Id == 0 || c.Email == ""{
		return nil, errors.New("Id and Email field must not be empty!")
	}

	c.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &AccessKey{Token:tokenString}, nil

}

func ClaimsFromToken(key []byte, token []byte) (*Claims, error){

	claims := &Claims{}
	t, err := jwt.ParseWithClaims(string(token), claims, func(token *jwt.Token) (i interface{}, e error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("Token Invalid")
	}

	return claims, nil

}
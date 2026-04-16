package utils

import (
	"errors"
	"hiliriset_ecoprint_golang/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type UserClaims struct{
	jwt.RegisteredClaims
	Username string `json:"username"`
	Email string `json:"email"`
}

func createClaims(username, email string) UserClaims {
	expiryMinute, err := strconv.Atoi(config.APPConfig.JWTExpireMinutes);
	
	if err != nil {
		expiryMinute = 1800
	}

	login_expiration_duration := time.Now().Add(time.Duration(expiryMinute) * time.Minute)

	return UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "Ecoprint Golang",
			ExpiresAt: jwt.NewNumericDate(login_expiration_duration),
		},
		Username: username,
		Email: email,
	}
}

func GenerateToken(username, email string) (string, error){
	claims := createClaims(username, email) 

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(config.APPConfig.JWTSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(config.APPConfig.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		return errors.New("Authentication Failed, Used Invalid Token")
	}

	return nil
}


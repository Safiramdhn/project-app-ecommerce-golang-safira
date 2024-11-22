package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func GenerateToken(userId string, config Configuration) (string, error) {
	claim := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	if config.Jwtkey != "" {
		secretKey = []byte(config.Jwtkey)
	}
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

func VerifyToken(tokenString string, config Configuration) (jwt.MapClaims, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		var secretKey []byte
		if config.Jwtkey != "" {
			secretKey = []byte(config.Jwtkey)
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims if the token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check the expiration time
	if exp, ok := claims["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return nil, fmt.Errorf("token is expired")
		}
	} else {
		return nil, fmt.Errorf("expiration claim (exp) is missing or invalid")
	}

	return claims, nil
}

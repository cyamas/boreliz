package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJWT(username string, id int) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading file: %v", err)
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := Claims{
		Username: username,
		UserID:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "boreliz",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	return tokenString
}

func DecodeJWT(tokenString string) (jwt.MapClaims, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading file: %v", err)
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

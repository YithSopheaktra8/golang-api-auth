package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(
	userID string,
	email string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "access",
		"exp": time.Now().
			Add(15 * time.Minute).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}

func GenerateRefreshToken(userID string) (string, string, error) {

	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"jti":     jti,
		"exp": time.Now().
			Add(30 * 24 * time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)

	if err != nil {
		return "", "", err
	}

	return tokenString, jti, nil
}

func ValidateToken(
	tokenString string,
) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			// Prevent algorithm attack
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, errors.New(
					"invalid signing method",
				)

			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {

		return nil, err

	}

	if !token.Valid {

		return nil, errors.New(
			"invalid token",
		)

	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {

		return nil, errors.New(
			"invalid claims",
		)

	}

	return claims, nil

}

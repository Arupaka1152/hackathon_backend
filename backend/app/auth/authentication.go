package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateToken(accountId string) (signedToken string, err error) {
	claims := jwt.MapClaims{
		"account_id": accountId,
		"exp":        time.Now().Add(60 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//secretKey := os.Getenv("secret_key")
	secretKey := "secret"

	return token.SignedString([]byte(secretKey))
}

func ParseToken(signedToken string) (accountId string, err error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		//secretKey := os.Getenv("secret_key")
		secretKey := "secret"
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("claims not found")
	}
	accountId, ok = claims["account_id"].(string)
	if !ok {
		return "", errors.New("account_id not found")
	}

	return accountId, nil
}

package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leomarzochi/facebooklike/cmd/config"
)

func CreateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()
	claims["userID"] = id

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(r *http.Request) error {
	tokenAsString := getTokenFromHeader(r)
	token, err := jwt.Parse(tokenAsString, verificationKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func getTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func verificationKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}

func GetIDFromToken(r *http.Request) (uint64, error) {
	tokenString := getTokenFromHeader(r)
	token, err := jwt.Parse(tokenString, verificationKey)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userID"]), 10, 64)
		if err != nil {
			return 0, err
		}

		return id, nil
	}

	return 0, err
}

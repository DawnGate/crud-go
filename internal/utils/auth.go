package utils

import (
	"errors"
	"strconv"
	"time"

	"sos-crud/internal/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var secRetToken = "super-secret"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func GenerateJwT(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	userID := strconv.FormatUint(uint64(user.ID), 10)

	claims := &jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secRetToken))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (uint, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secRetToken), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)

	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID, err := strconv.Atoi(claims.Subject)

	if err != nil {
		return 0, errors.New("Invalid user Id")
	}

	return uint(userID), nil
}

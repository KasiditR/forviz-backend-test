package utils

import (
	"github.com/KasiditR/forviz-backend-api-test/internal/config"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"strings"
	"time"
)

type SignedDetails struct {
	ID       string
	UserName string
	jwt.StandardClaims
}

var ACCESS_TOKEN_SECRET string = config.LoadConfig().AccessTokenSecret
var REFRESh_TOKEN_SECRET string = config.LoadConfig().RefreshTokenSecret

func TokenGenerator(id string, username string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		ID:       id,
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(REFRESh_TOKEN_SECRET))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	signedToken = strings.TrimPrefix(signedToken, "Bearer ")

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(ACCESS_TOKEN_SECRET), nil
	})

	if err != nil {
		msg = err.Error()
		log.Println(msg)
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "token invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return
	}

	return claims, msg
}

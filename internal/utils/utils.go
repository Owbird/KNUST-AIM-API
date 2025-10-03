package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/golang-jwt/jwt/v5"
)

// Spins up a new browser with no sandbox for linux systems
func NewBrowser() *rod.Browser {
	controlUrl := launcher.New().NoSandbox(true).MustLaunch()

	browser := rod.New().ControlURL(controlUrl).MustConnect().WithPanic(func(i interface{}) {
		log.Println("[!] Headerless browser proberly lost context.")
	})

	return browser
}


func GetCookiesFromJWT(tokenString string) (models.UserCookies, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return models.UserCookies{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return models.UserCookies{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.UserCookies{}, fmt.Errorf("failed to extract claims")
	}

	tokenData, ok := claims["token"].(map[string]interface{})
	if !ok {
		return models.UserCookies{}, fmt.Errorf("token claim not found or invalid format")
	}

	userCookies := models.UserCookies{}

	if antiforgery, ok := tokenData["Antiforgery"].(string); ok {
		userCookies.Antiforgery = antiforgery
	}

	if session, ok := tokenData["Session"].(string); ok {
		userCookies.Session = session
	}

	if identity, ok := tokenData["Identity"].(string); ok {
		userCookies.Identity = identity
	}

	return userCookies, nil
}


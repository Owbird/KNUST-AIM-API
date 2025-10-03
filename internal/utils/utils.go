package utils

import (
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
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return models.UserCookies{}, err
	}

	data, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return models.UserCookies{}, err
	}

	userCookies := models.UserCookies{}

	userCookies.Antiforgery = data["token"].(map[string]interface{})["antiforgery"].(string)
	userCookies.Session = data["token"].(map[string]interface{})["session"].(string)
	userCookies.Identity = data["token"].(map[string]interface{})["identity"].(string)

	return userCookies, nil
}

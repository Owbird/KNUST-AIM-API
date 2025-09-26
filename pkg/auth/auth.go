package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/internal/utils"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/go-rod/rod/lib/proto"
	"github.com/golang-jwt/jwt/v5"
)

type AuthFunctions struct{}

func NewAuthFunctions() *AuthFunctions {
	return &AuthFunctions{}
}

func (af *AuthFunctions) AuthenticateUser(payload models.UserAuthPayload) (string, error) {
	browser := utils.NewBrowser()

	page := browser.MustPage()

	defer page.Close()

	err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: config.UserAgent,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = page.Navigate(config.BaseUrl)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = page.WaitLoad()
	if err != nil {
		log.Println(err)
		return "", err
	}

	page.MustWaitStable()

	form := page.MustElement("form")

	usernameInput := form.MustElement("input[name='studentUsername']")
	usernameInput.MustInput(payload.Username)

	passwordInput := form.MustElement("input[name='Password']")
	passwordInput.MustInput(payload.Password)

	studentIdInput := form.MustElement("input[name='StudentId']")
	studentIdInput.MustInput(payload.StudentId)

	loginBtn := form.MustElement("button[type='submit']")
	loginBtn.MustClick()

	page.MustWaitNavigation()

	page.MustWaitLoad()

	cookies, err := page.Cookies([]string{config.BaseUrl})
	if err != nil {
		log.Println(err)
		return "", err
	}

	userCookies := models.UserCookies{}

	for _, cookie := range cookies {
		switch cookie.Name {
		case ".AspNetCore.Antiforgery.oBcnM5PKSJA":
			userCookies.Antiforgery = cookie.Value
		case ".AspNetCore.Session":
			userCookies.Session = cookie.Value
		case ".AspNetCore.Identity.Application":
			userCookies.Identity = cookie.Value
		}
	}

	if userCookies.Identity != "" {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"token": userCookies,
				"exp":   time.Now().Add(time.Hour * 24).Unix(),
			})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Println(err)

			return "", err
		}

		return tokenString, nil

	}

	return "", fmt.Errorf("could not authenticate user")
}

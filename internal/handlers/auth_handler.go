package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Owbird/KNUST-AIM-API/config"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary Authenticate a user
// @Description Authenticates the user the based on the credentials and returns a token which will be used to authorize requests
// @Tags Auth
// @Produce json
// @Accept  json
// @Param  username body string true "Username"
// @Param  password body string true "Password"
// @Param  studentId body string true "Student ID"
// @Success 200 {object} models.UserResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *Handlers) AuthHandler(c *gin.Context) {

	var authPayload models.UserAuthPayload

	c.BindJSON(&authPayload)

	var browser = rod.New().MustConnect().WithPanic(func(i interface{}) {
		log.Println("[!] Headerless browser proberly lost context.")
	})

	page := browser.MustPage()

	defer page.Close()

	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: config.UserAgent,
	})

	page.Navigate(config.BaseUrl)

	page.WaitLoad()

	page.MustWaitStable()

	form := page.MustElement("form")

	usernameInput := form.MustElement("input[name='studentUsername']")
	usernameInput.MustInput(authPayload.Username)

	passwordInput := form.MustElement("input[name='Password']")
	passwordInput.MustInput(authPayload.Password)

	studentIdInput := form.MustElement("input[name='StudentId']")
	studentIdInput.MustInput(authPayload.StudentId)

	loginBtn := form.MustElement("button[type='submit']")
	loginBtn.MustClick()

	page.MustWaitNavigation()

	page.MustWaitLoad()

	cookies, err := page.Cookies([]string{config.BaseUrl})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't authorize user. Please try again",
		})
	}

	if len(cookies) == 1 && cookies[0].Name == ".AspNetCore.Antiforgery.oBcnM5PKSJA" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Credentials are incorrect. Please try again",
		})
		return
	}

	if len(cookies) >= 3 {
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

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"token": userCookies,
				"exp":   time.Now().Add(time.Hour * 24).Unix(),
			})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Message: "Couldn't authorize user. Please try again",
			})

			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			Message: "User authorized successfully",
			Token:   tokenString,
		})

		return
	}

	c.JSON(http.StatusInternalServerError, models.ErrorResponse{
		Message: "Couldn't authorize user. Please try again",
	})
}

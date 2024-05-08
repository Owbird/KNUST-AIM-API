package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Please provide a user token"})

		return
	}

	auth := strings.ReplaceAll(authHeader, "Bearer ", "")

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Couldn't authorize user. Invalid token"})

		return
	}

	data, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Couldn't authorize user. Please try again"})

		return
	}

	userCookies := models.UserCookies{}

	userCookies.Antiforgery = data["token"].(map[string]interface{})["antiforgery"].(string)
	userCookies.Session = data["token"].(map[string]interface{})["session"].(string)
	userCookies.Identity = data["token"].(map[string]interface{})["identity"].(string)

	c.Set("userCookies", userCookies)

	c.Next()
}

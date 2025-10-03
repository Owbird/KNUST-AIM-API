package middlewares

import (
	"net/http"
	"strings"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Please provide a user token"})

		return
	}

	tokenString := strings.ReplaceAll(authHeader, "Bearer ", "")

	c.Set("userCookies", tokenString)

	c.Next()
}

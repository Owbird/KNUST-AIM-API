package handlers

import (
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/Owbird/KNUST-AIM-API/pkg/auth"
	"github.com/gin-gonic/gin"
)

var authFunctions = auth.NewAuthFunctions()

// @Summary Authenticate a user
// @Description Authenticates the user based on the credentials and returns a token which will be used to authorize requests as a bearer token
// @Tags Auth
// @Produce json
// @Accept json
// @Param authPayload body models.UserAuthPayload true "User authentication credentials"
// @Success 200 {object} models.UserResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *Handlers) AuthHandler(c *gin.Context) {
	var authPayload models.UserAuthPayload
	err := c.BindJSON(&authPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request payload",
		})
		return
	}
	token, err := authFunctions.AuthenticateUser(authPayload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Couldn't authorize user. Please try again",
		})
		return
	}
	c.JSON(http.StatusOK, models.UserResponse{
		Message: "User authorized successfully",
		Token:   token,
	})
}

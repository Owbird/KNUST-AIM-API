package handlers

import (
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/Owbird/KNUST-AIM-API/pkg/auth"
	"github.com/gin-gonic/gin"
)

// @Summary Authenticate a user
// @Description Authenticates the user the based on the credentials and returns a token which will be used to authorize requests as a bearer token
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

	err := c.BindJSON(&authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't authorize user. Please try again",
		})
	}

	token, err := auth.AuthenticateUser(authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't authorize user. Please try again",
		})
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Message: "User authorized successfully",
		Token:   token,
	})
}

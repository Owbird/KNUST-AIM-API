package handlers

import (
	"net/http"

	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/Owbird/KNUST-AIM-API/pkg/user"
	"github.com/gin-gonic/gin"
)

var userFunctions = user.NewUserFunctions()

// @Summary Get User Data
// @Description Returns personal, programme and contact user data
// @Tags User
// @Produce json
// @Success 200 {object} models.UserDataResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /user [get]
// @Security ApiKeyAuth
func (h *Handlers) GetUserData(c *gin.Context) {
	cookies, _ := c.Get("userCookies")

	userData, err := userFunctions.GetUserData(cookies.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't get user data. Please try again",
		})
	}

	c.JSON(http.StatusOK, models.UserDataResponse{
		Message:  "Fetched user data successfully",
		UserData: userData,
	})
}

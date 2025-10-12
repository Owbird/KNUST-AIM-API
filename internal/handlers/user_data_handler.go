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

// @Summary User image
// @Description Serves up the user image based on the student id
// @Tags User
// @Produce json
// @Param  studentId path string true "Student ID"
// @Success 200 {string} string "OK"
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /user/image/{studentId} [get]
func (h *Handlers) GetUserImage(c *gin.Context) {
	id, ok := c.Params.Get("id")

	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid user id",
		})

		return
	}

	imageBytes, err := userFunctions.GetUserImage(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Couldn't fetch user image",
		})

		return
	}

	c.Data(http.StatusOK, "image/jpg", imageBytes)
}

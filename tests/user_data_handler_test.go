package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Owbird/KNUST-AIM-API/internal/handlers"
	"github.com/Owbird/KNUST-AIM-API/internal/middlewares"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestUserDataHandler(t *testing.T) {

	godotenv.Load("../.env")

	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.Use(middlewares.AuthMiddleware)

	router.GET("/user-data", handlers.GetUserData)

	req, err := http.NewRequest(http.MethodGet, "/user-data", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("TEST_JWT"))

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response models.UserDataResponse

	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, "Fetched user data successfully", response.Message)

	assert.NotEqual(t, response.UserData.Personal.Username, "")
	assert.NotEqual(t, response.UserData.Personal.Surname, "")
	assert.NotEqual(t, response.UserData.Personal.OtherNames, "")
	assert.NotEqual(t, response.UserData.Personal.Gender, "")
	assert.NotEqual(t, response.UserData.Personal.DateOfBirth, "")
	assert.NotEqual(t, response.UserData.Personal.Country, "")
	assert.NotEqual(t, response.UserData.Personal.Region, "")
	assert.NotEqual(t, response.UserData.Personal.Religion, "")

	assert.NotEqual(t, response.UserData.Programme.StudentId, "")
	assert.NotEqual(t, response.UserData.Programme.IndexNo, "")
	assert.NotEqual(t, response.UserData.Programme.ProgrammeStream, "")

	assert.NotEqual(t, response.UserData.Contact.SchoolEmail, "")
	assert.NotEqual(t, response.UserData.Contact.PersonalEmail, "")
	assert.NotEqual(t, response.UserData.Contact.KNUSTMobile, "")
	assert.NotEqual(t, response.UserData.Contact.PersonalMobile, "")
	assert.NotEqual(t, response.UserData.Contact.AltPersonalMobile, "")
	assert.NotEqual(t, response.UserData.Contact.PostalAddress, "")
	assert.NotEqual(t, response.UserData.Contact.ResidentialAddress, "")

}

func TestGetUserImageHandler(t *testing.T) {

	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.GET("/image/:id", handlers.GetUserImage)

	req, err := http.NewRequest(http.MethodGet, "/image/20812892", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	assert.Equal(t, res.Header().Get("Content-Type"), "image/jpg")

	content, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, len(content), 0)

}

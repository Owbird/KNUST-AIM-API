package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Owbird/KNUST-AIM-Desktop-API/internal/handlers"
	"github.com/Owbird/KNUST-AIM-Desktop-API/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
)

type testCase struct {
	name        string
	expectedMsg string
	expected    int
	payload     models.UserAuthPayload
}

func TestAuthHandler(t *testing.T) {
	godotenv.Load("../.env")

	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.POST("/auth", handlers.AuthHandler)

	testCases := []testCase{
		{
			name:        "Valid credentials",
			expected:    http.StatusOK,
			expectedMsg: "User authorized successfully",
			payload: models.UserAuthPayload{
				Username:  os.Getenv("TEST_USERNAME"),
				Password:  os.Getenv("TEST_PASSWORD"),
				StudentId: os.Getenv("TEST_STUDENTID"),
			},
		},
		{
			name:        "Invalid credentials",
			expected:    http.StatusUnauthorized,
			expectedMsg: "Credentials are incorrect. Please try again",
			payload: models.UserAuthPayload{
				Username:  "wrong_username",
				Password:  "wrong_password",
				StudentId: "wrong_student_id",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payloadJSON, err := json.Marshal(tc.payload)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payloadJSON))

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			assert.Equal(t, tc.expected, res.Code)

			var response models.UserResponse

			json.NewDecoder(res.Body).Decode(&response)

			assert.Equal(t, tc.expectedMsg, response.Message)
		})
	}
}

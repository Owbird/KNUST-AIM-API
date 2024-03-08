package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestResultSelectionHandler(t *testing.T) {

	err := godotenv.Load("../.env")

	if err != nil {
		t.Fatal(err)
	}

	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.Use(middlewares.AuthMiddleware)

	router.GET("/results", handlers.ResultSelectionHandler)

	req, err := http.NewRequest(http.MethodGet, "/results", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("TEST_JWT"))

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response models.ResultsSelectionResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Fetched results successfully", response.Message)

	assert.NotEqual(t, len(response.Results.Sems), 0)
	assert.NotEqual(t, len(response.Results.Years), 0)
}

func TestGetResultsHandler(t *testing.T) {

	err := godotenv.Load("../.env")

	if err != nil {
		t.Fatal(err)
	}

	router := gin.Default()
	router.Use(middlewares.AuthMiddleware)

	handlers := handlers.NewHandlers()

	router.POST("/results", handlers.GetResultsHandler)

	testCases := []struct {
		name        string
		expectedMsg string
		expected    int
		payload     models.GetResultsPayload
	}{
		{
			name:        "Invalid year",
			expected:    http.StatusBadRequest,
			expectedMsg: "Unknown year. Please use the available academic years.",
			payload: models.GetResultsPayload{
				Year: "2096",
				Sem:  "2",
			},
		},
		{
			name:        "Invalid sem",
			expected:    http.StatusBadRequest,
			expectedMsg: "Unknown sem. Please use the available academic sems.",
			payload: models.GetResultsPayload{
				Year: "2022",
				Sem:  "56",
			},
		},
		{
			name:        "Results not available",
			expected:    http.StatusBadRequest,
			expectedMsg: "Results not available for the chosen academic year and sem",
			payload: models.GetResultsPayload{
				Year: "2024",
				Sem:  "1",
			},
		},
		{
			name:        "Results available",
			expected:    http.StatusOK,
			expectedMsg: "Fetched results successfully",
			payload: models.GetResultsPayload{
				Year: "2021",
				Sem:  "1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payloadJSON, err := json.Marshal(tc.payload)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/results", bytes.NewBuffer(payloadJSON))

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TEST_JWT")))

			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			assert.Equal(t, tc.expected, res.Code)

			var response models.GetResultsResponse

			err = json.NewDecoder(res.Body).Decode(&response)

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expectedMsg, response.Message)

			if res.Code == http.StatusOK {
				assert.NotEqual(t, response.PersonalData.Name, "")
				assert.NotEqual(t, response.PersonalData.Year, "")
				assert.NotEqual(t, response.PersonalData.IndexNo, "")
				assert.NotEqual(t, response.PersonalData.Programme, "")
				assert.NotEqual(t, response.PersonalData.StudentID, "")
				assert.NotEqual(t, response.PersonalData.Date, "")
				assert.NotEqual(t, response.PersonalData.Option, "")
				assert.NotEqual(t, response.PersonalData.Username, "")

				for _, result := range response.Results {
					assert.NotEqual(t, result.CourseCode, "")
					assert.NotEqual(t, result.CourseName, "")
					assert.NotEqual(t, result.Credits, "")
					assert.NotEqual(t, result.Grade, "")
					assert.NotEqual(t, result.TotalMark, "")
				}

				assert.NotEqual(t, response.Summary.CreditsRegistered.Semester, "")
				assert.NotEqual(t, response.Summary.CreditsRegistered.Cumulative, "")
				assert.NotEqual(t, response.Summary.CreditsObtained.Semester, "")
				assert.NotEqual(t, response.Summary.CreditsObtained.Cumulative, "")
				assert.NotEqual(t, response.Summary.CreditsCalculated.Semester, "")
				assert.NotEqual(t, response.Summary.CreditsCalculated.Cumulative, "")
				assert.NotEqual(t, response.Summary.WeightedMarks.Semester, "")
				assert.NotEqual(t, response.Summary.WeightedMarks.Cumulative, "")
				assert.NotEqual(t, response.Summary.CWA.Semester, "")
				assert.NotEqual(t, response.Summary.CWA.Cumulative, "")
			}
		})
	}
}

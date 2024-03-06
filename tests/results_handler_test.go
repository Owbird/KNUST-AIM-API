package tests

import (
	"encoding/json"
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

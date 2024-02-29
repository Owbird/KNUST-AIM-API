package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owbird/KNUST-AIM-Desktop-API/internal/handlers"
	"github.com/Owbird/KNUST-AIM-Desktop-API/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewsHandler(t *testing.T) {

	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.GET("/news", handlers.GetNewsHandler)

	req, err := http.NewRequest(http.MethodGet, "/news", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response models.NewsResponse

	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, "Fetched news successfully", response.Message)

	for index, news := range response.News {
		t.Run(fmt.Sprintf("Non Empty field for #%v", index), func(t *testing.T) {
			assert.NotEqual(t, news.Category, "")
			assert.NotEqual(t, news.Date, "")
			assert.NotEqual(t, news.Title, "")
			assert.NotEqual(t, news.Description, "")
			assert.NotEqual(t, news.Slug, "")
		})

	}
}

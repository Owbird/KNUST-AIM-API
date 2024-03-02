package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owbird/KNUST-AIM-API/internal/handlers"
	"github.com/Owbird/KNUST-AIM-API/models"
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

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		log.Fatal(err)
	}

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

func TestNewsDetailsHandler(t *testing.T) {

	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.GET("/news/:slug", handlers.GetNewsDetailsHandler)

	req, err := http.NewRequest(http.MethodGet, "/news/knust-set-represent-ghana-jessup-international-law-moot-court-competition-usa", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response models.NewsDetailsResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "Fetched news successfully", response.Message)

	assert.NotEqual(t, response.News.Title, "")
	assert.NotEqual(t, response.News.Date, "")
	assert.NotEqual(t, response.News.Source, "")
	assert.NotEqual(t, response.News.FeaturedImage, "")

	assert.NotEqual(t, len(response.News.Content), 0)

	for index, content := range response.News.Content {
		t.Run(fmt.Sprintf("Non Empty field for #%v", index), func(t *testing.T) {
			assert.NotEqual(t, content.Type, "")
			assert.NotEqual(t, content.Value, "")
		})
	}

}

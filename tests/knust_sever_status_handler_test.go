package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owbird/KNUST-AIM-API/internal/handlers"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestKNUSTServerStatusHandler(t *testing.T) {
	router := gin.Default()
	handlers := handlers.NewHandlers()

	router.GET("/knust-server-status", handlers.KNUSTServerStatusHandler)

	req, err := http.NewRequest(http.MethodGet, "/knust-server-status", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var response models.KNUSTServerStatusResponse

	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, "Fetched server status successfully", response.Message)

	assert.NotEqual(t, response.Badge, "")

	for _, server := range response.Servers {
		t.Run(fmt.Sprintf("Server status for %s", server.Url), func(t *testing.T) {
			assert.Contains(t, []string{"Up", "Down"}, server.Status)
		})
	}

}

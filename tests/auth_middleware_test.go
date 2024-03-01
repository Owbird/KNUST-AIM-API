package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owbird/KNUST-AIM-API/internal/middlewares"
	"github.com/Owbird/KNUST-AIM-API/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	router := gin.Default()

	router.Use(middlewares.AuthMiddleware)

	router.GET("/user", func(ctx *gin.Context) {})

	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "No auth token",
			token: "",
		},
		{
			name:  "Invalid token",
			token: "TEST_TOKEN",
		},
		{
			name:  "Expired token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxODE5NzIsInRva2VuIjp7ImFudGlmb3JnZXJ5IjoiQ2ZESjhQeEZkRDQ3d3lKRm1hV2ZtM0JQYTlMNHIwSUc5VnNBRDVHWmVwWTFmSzNLTGtSVkQ4LXFWb1R0QzJUMXdzbFlPOXFVMjBaNjJvRC1PSmxaMVFxQlByZDR0MlNqblJKTUlqNnh2M2hzUm5YSlZjeExTb0RfN1JMMTJtZlNtaEUzb3ljbGticWFOLVRqa3l3MGk5eUFaV0EiLCJzZXNzaW9uIjoiQ2ZESjhQeEZkRDQ3d3lKRm1hV2ZtM0JQYTlLUk5DZkJBMkliN0lYU1AwZkNGViUyQnVLSWFCcmFUVlZweW8lMkY4WlpyYXdqbGRFWURSTkt5alUyRjBIaDNDNVFGdmJ6UkdHNTJTbVVFeDNlNmt1SGRtRGt1NEZIVDBXNWpUYndQazhLZzNyRk9mRHFjOVg3elVCcTM0ViUyQmo5Vm5tQUtKNTVUaGFsUSUyQkh2ZXlGRk1DZHZ4eSIsImlkZW50aXR5IjoiQ2ZESjhQeEZkRDQ3d3lKRm1hV2ZtM0JQYTlJcEI0OE8tWW0ybW80OWkyNkFwMl9Ick4tUERVQ0Y1ZFJodVZDaDRtbnlDV1lKNnhXdXBNaHlwODlCT2NULUVVZzBYT0lHQzFtQy1jZmlPR3Mzel9qOU9CQU9pTElDVmlabzRzSWVacllHUHlDTG9jT0h1SE4zSkxYajA0ZGQxTXF6U0VSNG93bWsxMElkYTVjMWpVN0VVVGpWVXlsYWR1NGhOQkZwakEzODkzVDBRZ1hDMnVGNV91TFlqWF9WTlcycHcyaFhsQ3pHQ1dUaHNZTjhBRWUzS1N4akpiMDdMYTE4Zl90TUU4eDM0Tm91VFNGSnQ5OGFIcXJjQ1dIN2FVb3Nlajk2YVhGaDJqdk9KUF8tOXFfRXJsRHl1a0FPeE9udEhpNjk1dWJZa0VVa0tEMklfWWtWNjRycTJaelgtMHRIOUlJRnBtMzVRaU91aGF5V2ZIMWY4RUZHeWcyMWJrUEo2M0hOR2dSYzY5bmNGSmNtd3lOc1dKbkd5SjV5NEZRRzlpbjAxanZmakNRTm11UWFLaTZuNGJfemxKRDh1c1RPVnVoRU43dHlkbHhONFVIbWVjemFTU1R6WlEwY1BncnMtMWhQV2szV3pHa21hbXBSTWVmYjBNNXphRVlraFhkTkYyNlcwUVVwLUQzWnlXb2pVem8zbW84TExHLU4tNEtBQ19udkwxY3R1dm51OGk3SV81SzRUbHVwWFhMU1hGSjRxNFlZQkU2d05OUDNpbnR2b2dKN1IwU1JBdUNHYzE0MENhU1JBYW1YSWYzbUpLaElIazZpWlBnTDVwZ25ZUzhpZ0pkdHNVZ3kwR2M1RUVMM24yTlJHSURBVXVVRFFLZ2d1R1VhUEdsRGhUbUJwRHJvLVlYZ3g0TjZBUUFjRDJtWHp6Mldta3FDQWtkcDd2UFhKN0REekpaVEk2eUZMWHNWUkhJOE1nYWI5NXA2Zm9nVndUTEIxYW1FSDNaS21nME1XXy0tMEZPU3RQak9nZ2pEaG5hT3ZGOXE4VEVRSV95RWxRNGx5S2hobGdVVW93a2RnaUFVLW9IakZ5YmNJa0dMQ1Zpb2FReGNBV0RhWnlwbV9Gb1VYNlVuSFF0M3BFWSJ9fQ.Y-UAyxOw5CzJzs8EL7QKID6ZybL4bqfZwoH9dDVSR5c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, "/user", nil)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))

			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			assert.Equal(t, http.StatusUnauthorized, res.Code)

			var response models.ErrorResponse

			json.NewDecoder(res.Body).Decode(&response)

			assert.Equal(t, "Couldn't authorize user", response.Message)
		})
	}
}

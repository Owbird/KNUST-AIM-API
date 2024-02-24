package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Owbird/KNUST-AIM-Desktop-API/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	handlers := handlers.NewHandlers()

	router.GET("/", handlers.HelloHandler)

	api := router.Group("/api")

	apiV1 := api.Group("/v1")

	auth := apiV1.Group("/auth")
	{
		auth.POST("/login", handlers.AuthHandler)
	}

	return router
}

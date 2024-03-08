package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Owbird/KNUST-AIM-API/docs"
	"github.com/Owbird/KNUST-AIM-API/internal/handlers"
	"github.com/Owbird/KNUST-AIM-API/internal/middlewares"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	docs.SwaggerInfo.Title = "KNUST AIM API"
	docs.SwaggerInfo.Description = "An unofficial API for KNUST AIM"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api")
	{
		apiV1 := api.Group("/v1")
		{
			techServers := apiV1.Group("/knust-server-status")
			{
				techServers.GET("/", handlers.KNUSTServerStatusHandler)
				techServers.GET("/badge", handlers.KNUSTServerStatusBadgeHandler)
			}

			auth := apiV1.Group("/auth")
			{
				auth.POST("/login", handlers.AuthHandler)
			}

			news := apiV1.Group("/news")
			{
				news.GET("/", handlers.GetNewsHandler)
				news.GET("/:slug", handlers.GetNewsDetailsHandler)
			}

			user := apiV1.Group("/user")
			{
				user.Use(middlewares.AuthMiddleware)
				user.GET("/", handlers.GetUserData)
				user.GET("/image/:id", handlers.GetUserImage)
				user.POST("/results", handlers.GetResultsHandler)
				user.GET("/results/selection", handlers.ResultSelectionHandler)
			}
		}
	}

	return router
}

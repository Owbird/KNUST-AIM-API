package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Handlers struct {
	Browser *rod.Browser
}

func NewHandlers() *Handlers {
	controlUrl := launcher.New().NoSandbox(true).MustLaunch()

	var browser = rod.New().ControlURL(controlUrl).MustConnect().WithPanic(func(i interface{}) {
		log.Println("[!] Headerless browser proberly lost context.")
	})

	return &Handlers{
		Browser: browser,
	}
}

func (h *Handlers) HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

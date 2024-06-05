package utils

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Spins up a new browser with no sandbox for linux systems
func NewBrowser() *rod.Browser {
	controlUrl := launcher.New().NoSandbox(true).MustLaunch()

	browser := rod.New().ControlURL(controlUrl).MustConnect().WithPanic(func(i interface{}) {
		log.Println("[!] Headerless browser proberly lost context.")
	})

	return browser
}

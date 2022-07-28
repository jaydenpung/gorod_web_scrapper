package main

import (
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	var (
		username = os.Getenv("WEB_USERNAME")
		password = os.Getenv("WEB_PASSWORD")
		headless = os.Getenv("HEADLESS_MODE") == "true"
	)

	l := launcher.New().
		Headless(headless).
		Devtools(false)

	defer l.Cleanup() // remove launcher.FlagUserDataDir

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url)

	if !headless {
		browser.Trace(true).
			SlowMotion(1 * time.Second)
		launcher.Open(browser.ServeMonitor(""))
	}

	browser.MustConnect()
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage("https://www.unstructureddataterminal.com/PagesRegister")

	page.MustElement("#input-with-icon-adornment").Input(username)
	page.MustElement("#standard-adornment-password").Input(password)
	page.MustElement("form").MustElement("button").MustClick()

	if !headless {
		time.Sleep(1 * time.Hour)
	}
}

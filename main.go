package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

type CheckList struct {
	homeWatchlist int
	homeSec       int
	homeSedar     int
	gridWatchlist int
	gridRecent    int
	gridComplete  int
}

func main() {

	// load env files
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	err = godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load(exPath + "/.env")
		if err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}

	var checkList CheckList

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
			SlowMotion(200 * time.Millisecond)
		launcher.Open(browser.ServeMonitor(""))
	}

	browser.MustConnect()
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage("https://www.unstructureddataterminal.com/PagesRegister")

	// login
	page.MustElement("#input-with-icon-adornment").Input(username)
	page.MustElement("#standard-adornment-password").Input(password)
	page.MustElement("form").MustElement("button").MustClick()

	// home
	page.MustElement("div.app-header-menu > button:nth-child(1)").MustClick()
	wait()
	// watchlist
	page.MustElement("div.card-header > div.MuiButtonGroup-root > button:nth-child(1)").MustClick()
	wait()
	checkList.homeWatchlist = len(page.MustElement(".ag-center-cols-container").MustElements("[role=\"row\"]"))
	// sec
	page.MustElement("div.card-header > div.MuiButtonGroup-root > button:nth-child(2)").MustClick()
	wait()
	checkList.homeSec = len(page.MustElement(".ag-center-cols-container").MustElements("[role=\"row\"]"))
	// sedar
	page.MustElement("div.card-header > div.MuiButtonGroup-root > button:nth-child(3)").MustClick()
	wait()
	checkList.homeSedar = len(page.MustElement(".ag-center-cols-container").MustElements("[role=\"row\"]"))

	// grid
	page.MustElement("div.app-header-menu > button:nth-child(2)").MustClick()
	wait()
	// watchlist
	page.MustElement("div.MuiGrid-root.MuiGrid-item.MuiGrid-grid-xs-8 > div > div > div > div:nth-child(3) > div.MuiButtonGroup-root > button:nth-child(1)").MustClick()
	wait()
	checkList.gridWatchlist = len(page.MustElement(".ag-pinned-left-cols-container").MustElements("[role=\"row\"]"))
	// recent
	page.MustElement("div.MuiGrid-root.MuiGrid-item.MuiGrid-grid-xs-8 > div > div > div > div:nth-child(3) > div.MuiButtonGroup-root > button:nth-child(2)").MustClick()
	wait()
	checkList.gridRecent = len(page.MustElement(".ag-pinned-left-cols-container").MustElements("[role=\"row\"]"))
	// complete
	page.MustElement("div.MuiGrid-root.MuiGrid-item.MuiGrid-grid-xs-8 > div > div > div > div:nth-child(3) > div.MuiButtonGroup-root > button:nth-child(3)").MustClick()
	wait()
	checkList.gridComplete = len(page.MustElement(".ag-pinned-left-cols-container").MustElements("[role=\"row\"]"))

	fmt.Printf("P \"unstructureddataterminal\" homewatchlist=%v;0;0|homesec=%v;0;0|homesedar=%v;0;0|gridwatchlist=%v;0;0|gridrecent=%v;0;0|gridcomplete=%v;0;0 flashboard rows", checkList.homeWatchlist, checkList.homeSec, checkList.homeSedar, checkList.gridWatchlist, checkList.gridRecent, checkList.gridComplete)

	if !headless {
		time.Sleep(1 * time.Hour)
	}
}

func wait() {
	time.Sleep(5 * time.Second)
}

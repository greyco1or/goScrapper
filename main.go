package main

import (
	"github.com/greyco1or/goScrapper/scrapper"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment("jobs.csv", "job.csv")
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}

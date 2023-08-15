package main

import "github.com/greyco1or/goScrapper/scrapper"

const (
	goLang   string = "golang"
	javaLang string = "java"
)

func main() {
	scrapper.Scrape(javaLang)
}

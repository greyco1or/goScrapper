package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

var (
	baseURL  string = "https://www.jobkorea.co.kr/Search/?stext="
	goLang   string = "golang"
	javaLang string = "java&tabType=recruit"
)

type extractedJob struct {
	id       string
	title    string
	location string
	lang     string
	sort     string
}

func main() {
	totalPages := getPageNumber()
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageURL := baseURL + javaLang + "&Page_No=" + strconv.Itoa(page+1)
	fmt.Println("Requesting page: ", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".lists")
	fmt.Println("get card")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Find(".list-post").Attr("data-gno")
		if id != "false" && id != "" {
			fmt.Println(id)
			info := card.Find("data-gainfo")
			//fmt.Println(info.Find())
		}
	})
}

func getPageNumber() int {
	pages := 0
	res, err := http.Get(baseURL + javaLang)
	checkErr(err)
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".recruit-info .tplPagination.newVer.wide").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("ul li").Length()
	})

	defer res.Body.Close()

	return pages
}

func checkErr(err error) {
	if err != nil {
		//kill the program
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status code: ", res.StatusCode)
	}
}

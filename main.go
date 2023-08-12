package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	baseURL  string = "https://www.jobkorea.co.kr/Search/?stext="
	goLang   string = "golang"
	javaLang string = "java&tabType=recruit"
)

type extractJobData struct {
	ID       string `json:"dimension42"`
	Title    string `json:"dimension45"`
	Location string `json:"dimension46"`
	Company  string `json:"dimension48"`
	Sort     string `json:"dimension43"`
	Work     string `json:"dimension44"`
	/*
		Dimension65 string `json:"dimension65"`
		Dimension66 string `json:"dimension66"`
		Dimension70 string `json:"dimension70"`
		Dimension47 string `json:"dimension47"`
	*/
}

func main() {
	var jobs []extractJobData
	c := make(chan []extractJobData)
	totalPages := getPageNumber()
	for i := 1; i < totalPages+1; i++ {
		go getPage(i, c)
	}

	for i := 1; i < totalPages+1; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	writeJobs(jobs)
	fmt.Printf("Done, Writing on csv files: %v 개\n", len(jobs))
}

func getPage(page int, mainC chan<- []extractJobData) {
	var jobs []extractJobData
	c := make(chan extractJobData)
	pageURL := baseURL + javaLang + "&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting page: ", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".lists .clear .list-post")
	fmt.Println("GET CARD")
	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		if job.ID == "" {
			continue
		}
		jobs = append(jobs, job)
		fmt.Println("GET PAGE FINISH")
		//fmt.Printf("jobs: %+v\n", jobs)
	}
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractJobData) {
	id, _ := card.Attr("data-gno")
	if id == "" {
		c <- extractJobData{}
		return
	}
	infoData, _ := card.Attr("data-gainfo")
	if infoData == "" {
		c <- extractJobData{}
		return
	}
	var jsonData extractJobData
	err := json.Unmarshal([]byte(infoData), &jsonData)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("jsonData: %+v\n", jsonData)
	fmt.Println("EXTRACT JOB FINISH")
	c <- jsonData
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

func writeJobs(jobs []extractJobData) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Company", "Sort", "Work"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.jobkorea.co.kr/Recruit/GI_Read/" + job.ID, job.Title, job.Location, job.Company, job.Sort, job.Work}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

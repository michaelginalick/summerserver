package main

import (
	"errors"
	"log"
	"net/http"
	"./db"
	"./structs"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"../sharedFunctions/stringutil"
    "../sharedFunctions/numutil"
)

const link = "https://www.choosechicago.com/events-and-shows/festivals-guide/"

func scrapeEventPage() {
	// Request the HTML page.
	res, err := http.Get(link)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	eventList := event.List()

	// Find the events
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		name := s.Find("a").Text()
		link, _ := s.Find("a").Attr("href")
		date := s.Text()
		month, i, err := str.ExtractMonthDate(date)
		days := str.ExtractDays(date, i)
		year, _ := extractYear(days)
		individualDays := make([]string, 0)

		if err != nil {
			log.Println("Requested item not found")
		}

		//remove year from days slice
		if len(days) > 0 {
			days = days[:len(days)-1]
		}

		if len(days) >= 1 {
			firstInt, lastInt := str.FirstAndLastElement(days)
			individualDays = nums.GetIndividualDays(firstInt, lastInt)
		}

		newEvent := &event.Event{0, name, link, month, days, year, individualDays, 0, "", nil}
		eventList = event.AddBeginning(newEvent, eventList)
	})

	res.Body.Close()

	db.SaveRecords(eventList)
}

func main() {
	scrapeEventPage()
}

func extractYear(days []string) (string, error) {

	if len(days) > 0 {
		return days[len(days)-1], nil
	}
	return "", errors.New("cannot return year")
}

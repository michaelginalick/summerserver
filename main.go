package main

import(
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"./calendar"
	"errors"
)

const link = "https://www.choosechicago.com/events-and-shows/festivals-guide/"


type event struct {
	name string
	link string
	date string
	next *event
}


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
	eventList := &event{}

  // Find the events
  doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		name := s.Find("a").Text()
		link, _ := s.Find("a").Attr("href")
		date := s.Text()

		month, err := extractMonthDate(date)
		
		if err != nil {
			fmt.Println(err)
		}

		newEvent := &event{name, link, month, nil }

		eventList = addBeginning(newEvent, eventList)
	})
	
	res.Body.Close()

	// printList(eventList)
}

func main() {
  scrapeEventPage()
}

func addBeginning(newEvent, eventList *event) *event {
	newEvent.next = eventList
	return newEvent
}


func printList(eventList *event) {
	for i := eventList; i != nil; i = i.next {
		fmt.Println(i.name, i.date, i.link)
	}
}


func extractMonthDate(s string) (string, error) {
	parseDate := strings.Fields(s)

	for i := 0; i < len(parseDate); i++ {

		monthValue := calendar.GetMonth(parseDate[i])

		if monthValue != "" {
			extractDays(parseDate, i+1)
			return monthValue, nil
		}
	}
	return "", errors.New("No date is listed with this event")
}


func extractDays(s []string, i int) {
	
	x := ""
	for j := i; j < len(s); j++ {
		x += " "
		x += string(s[j])
	}

	fmt.Println(x)
}

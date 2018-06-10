package main

import(
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
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

		newEvent := &event{name, link, date, nil }

		eventList = addBeginning(newEvent, eventList)
	})
	
	res.Body.Close()

	printList(eventList)
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
		fmt.Println(i.name, i.date)
	}
}

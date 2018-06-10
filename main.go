package main

import(
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
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

  // Find the events
  doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		event := s.Find("a").Text()
		link, _ := s.Find("a").Attr("href")
		date := s.Text()
    fmt.Printf("%s - %s - %s\n", event, date, link)
  })
}

func main() {
  scrapeEventPage()
}

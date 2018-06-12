package main

import(
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"./calendar"
	"errors"
	"regexp"
	"strconv"
)

const link = "https://www.choosechicago.com/events-and-shows/festivals-guide/"


type event struct {
	name string
	link string
	month string
	days []string
	length int
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

		month, i, err := extractMonthDate(date)
		days := extractDays(date, i)

		//remove year
		if len(days) > 0 {
			days = days[:len(days)-1]
		}

		if err != nil {
			fmt.Println(err)
		}

		newEvent := &event{name, link, month, days, 5, nil }

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
		fmt.Println(i.name, i.link, i.month, i.days)
	}
}


func extractMonthDate(s string) (string, int, error) {
	parseDate := parseFields(s)

	for i := 0; i < len(parseDate); i++ {

		monthValue := calendar.GetMonth(parseDate[i])

		if monthValue != "" {
			return monthValue, i+1, nil
		}
	}
	return "", 0, errors.New("No date is listed with this event")
}


func extractDays(s string, i int) []string {
	date := parseFields(s)

	x := ""
	for j := i; j < len(date); j++ {
		x += " "
		x += string(date[j])
	}

	re := regexp.MustCompile("[0-9]+")
	return re.FindAllString(x, -1)
}


func parseFields(s string) []string {
	return strings.Fields(s)
}

//https://play.golang.org/p/BZdhROeZf2T

func convInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	
	if err != nil {
		return 0, errors.New("cannot convert to integer")
	}
	
	return i, nil
}


func getDays(firstInt, lastInt int) []int {
	var s []int
	
	i := firstInt 
	
	for i <= lastInt {
		s = append(s, i)
		i++
	}
	return s
}

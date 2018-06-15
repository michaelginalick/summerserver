package main

import (
	"./calendar"
	"./structs"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

	eventList := &event.Event{}

	// Find the events
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		name := s.Find("a").Text()
		link, _ := s.Find("a").Attr("href")
		date := s.Text()
		month, i, err := extractMonthDate(date)
		days := extractDays(date, i)
		individualDays := make([]int, 0)

		//remove year
		if len(days) > 0 {
			days = days[:len(days)-1]
		}

		if err != nil {
			log.Println("Requested item not found")
		}


		if len(days) > 1 {
			first := days[0]
			last := days[len(days)-1]
			
			firstInt, _ := convInt(first)
			lastInt, _ := convInt(last)
			
			individualDays = getIndividualDays(firstInt, lastInt)
		}

		newEvent := &event.Event{name, link, month, days, individualDays, len(individualDays), nil}

		eventList = addBeginning(newEvent, eventList)
	})

	res.Body.Close()

	printList(eventList)
	printListByMonth(eventList, "July")
	// saveEventListToDatabase()
}

func main() {
	scrapeEventPage()
}


func printListByMonth(e *event.Event, s string) {
	for i:=e; i != nil; i = i.Next {
		if i.Month == s {
			fmt.Println(i.Name, i.Link, i.Month, i.Days, i.IndividualDays)
		}
	}
}


func extractMonthDate(s string) (string, int, error) {
	parseDate := parseFields(s)

	for i := 0; i < len(parseDate); i++ {

		monthValue := calendar.GetMonth(parseDate[i])

		if monthValue != "" {
			return monthValue, i + 1, nil
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


func convInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("cannot convert to integer")
	}

	return i, nil
}

func getIndividualDays(firstInt, lastInt int) []int {
	var s []int

	i := firstInt

	for i <= lastInt {
		s = append(s, i)
		i++
	}
	return s
}


func addBeginning(newEvent, eventList *event.Event) *event.Event {
	newEvent.Next = eventList
	return newEvent
}

func printList(eventList *event.Event) {
	for i := eventList; i != nil; i = i.Next {
		fmt.Println(i.Name, i.Link, i.Month, i.Days, i.IndividualDays)
	}
}

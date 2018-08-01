package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"./calendar"
	"./db"
	"./structs"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
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
			firstInt, lastInt := firstAndLastElement(days)
			individualDays = getIndividualDays(firstInt, lastInt)
		}

		newEvent := &event.Event{0, name, link, month, days, year, individualDays, nil}
		eventList = event.AddBeginning(newEvent, eventList)
	})

	res.Body.Close()

	db.SaveRecords(eventList)
}

func main() {
	scrapeEventPage()
}

func firstAndLastElement(days []string) (int, int) {
	first := days[0]
	last := days[len(days)-1]

	firstInt, _ := convInt(first)
	lastInt, _ := convInt(last)

	return firstInt, lastInt
}

func extractMonthDate(s string) (string, int, error) {
	parseDate := parseFields(s)

	for i := 0; i < len(parseDate); i++ {

		monthValue := calendar.GetMonth(parseDate[i])

		if monthValue.Name != "" {
			return strings.ToLower(monthValue.Name), i + 1, nil
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

func getIndividualDays(firstInt, lastInt int) []string {
	var s []string

	i := firstInt

	for i <= lastInt {
		j := strconv.Itoa(i)
		s = append(s, j)
		i++
	}
	return s
}

func extractYear(days []string) (string, error) {

	if len(days) > 0 {
		return days[len(days)-1], nil
	}
	return "", errors.New("cannot return year")
}

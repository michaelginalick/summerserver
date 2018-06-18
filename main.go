package main

import (
	"./calendar"
	"./structs"
	"./routes/api/v1"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"./db"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
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
		individualDays := make([]string, 0)

		//remove year assuming this is for 2018
		if len(days) > 0 {
			days = days[:len(days)-1]
		}

		if err != nil {
			log.Println("Requested item not found")
		}


		if len(days) > 1 {
			firstInt, lastInt := firstAndLastElement(days)
			individualDays = getIndividualDays(firstInt, lastInt)
		}

		newEvent := &event.Event{0, name, link, month, days, individualDays, len(individualDays), nil}

		eventList = event.AddBeginning(newEvent, eventList)
	})

	res.Body.Close()

	db.SaveRecords(eventList)
}

func main() {
	// scrapeEventPage()

	router := mux.NewRouter()
	router.HandleFunc("/events", eventsController.GetEvents).Methods("GET")

	router.HandleFunc("/events/{id}", eventsController.GetEvent).Methods("GET")
	router.HandleFunc("/events/{month}", eventsController.GetEventsByMonth).Methods("GET")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) 
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(allowedOrigins, allowedMethods)(router)))
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
			return monthValue.Name, i + 1, nil
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

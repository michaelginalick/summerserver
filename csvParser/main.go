package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"../db"
	nums "../sharedFunctions/numutil"
	str "../sharedFunctions/stringutil"
	event "../structs"
	_ "github.com/lib/pq"
)

func main() {
	csvFile, _ := os.Open("events_2019.csv")
	fileReader := csv.NewReader(bufio.NewReader(csvFile))
	eventList := event.List()

	parseDocument(fileReader, eventList)
}

func parseDocument(fileReader *csv.Reader, eventList *event.Event) {
	for {
		line, error := fileReader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		dateInfo := line[1]
		month, i, _ := str.ExtractMonthDate(dateInfo)
		days := str.ExtractDays(dateInfo, i)
		year := extractYear(dateInfo)
		individualDays := make([]string, 0)
		singleDay := extractDay(dateInfo)
		name := strings.Trim(line[0], " ")

		if len(days) >= 1 {
			firstInt, lastInt := str.FirstAndLastElement(days)
			individualDays = nums.GetIndividualDays(firstInt, lastInt)
		}

		newEvent := &event.Event{0, name, line[3], month, days, year, individualDays, singleDay, line[2], nil}
		eventList = event.AddBeginning(newEvent, eventList)
	}

	db.SaveRecords(eventList)
}

func extractDay(line string) int {
	return 0
}

func extractYear(line string) string {
	return "2019"
}

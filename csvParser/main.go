package main

import (
	"../sharedFunctions/numutil"
	"../sharedFunctions/stringutil"
	"../webScrapper/db"
	"../webScrapper/structs"
	"bufio"
	"encoding/csv"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
)

func main() {
	csvFile, _ := os.Open("events.csv")
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

		if len(days) >= 1 {
			firstInt, lastInt := str.FirstAndLastElement(days)
			individualDays = nums.GetIndividualDays(firstInt, lastInt)
		}

		newEvent := &event.Event{0, line[0], line[3], month, days, year, individualDays, singleDay, nil}
		eventList = event.AddBeginning(newEvent, eventList)
	}

	db.SaveRecords(eventList)
}

func extractDay(line string) int {
	return 0
}

func extractYear(line string) string {
	return "2018"
}

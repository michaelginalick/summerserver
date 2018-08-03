package main

import (
    "../webScrapper/structs"
    "../webScrapper/db"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
    "../sharedFunctions/stringutil"
    "../sharedFunctions/numutil"
    _ "github.com/lib/pq"
)

func main() {
	csvFile, _ := os.Open("events.csv")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    eventList := event.List()

	for {
		line, error := reader.Read()
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


func parseDocument() {

}

func extractDay(line string) int {
	return 0
}

func extractYear(line string) string {
	return "2018"
}

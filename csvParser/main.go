package main

import (
    "bufio"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "os"
)

type Event struct {
    Event    string   `json:"Event"`
		Date     string   `json:"Date"`
		Location string   `json:"Location"`
		Link     string   `json:"Link"`
}

func main() {
    csvFile, _ := os.Open("events.csv")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    var event []Event
    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }
        event = append(event, Event{
            Event:    line[0],
						Date:     line[1],
						Location: line[2],
						Link:     line[3],
        })
    }
    EventJSON, _ := json.Marshal(event)
    fmt.Println(string(EventJSON))
}

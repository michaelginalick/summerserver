package eventscontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../../../db"
	event "../../../structs/"
	"github.com/gorilla/mux"
)

// CreateEvent creates an event
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	newEvent := event.Event{}

	err := decoder.Decode(&newEvent)
	if err != nil {
		panic(err)
	}
	newEvent = event.Event{0, newEvent.Name, newEvent.Link,
		newEvent.Month, nil,
		newEvent.Year, newEvent.IndividualDays,
		0, newEvent.Location, nil}
	db.SaveRecords(&newEvent)
	setHeaders(w)
	return
}

// GetEvents :  returns all the events
func GetEvents(w http.ResponseWriter, r *http.Request) {

	db := db.OpenDB()
	defer db.Close()
	sqlStatement := `select events.id, name, link, month, year, day, location
									 from events
									 inner join days on events.id = days.event_id
									 order by days.day;`
	rows, _ := db.Query(sqlStatement)

	formatAndReturnJSONResponse(rows, w)
	db.Close()

	return
}

// GetEventByID : by event by id
func GetEventByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	db := db.OpenDB()
	defer db.Close()
	sqlStatement := `select events.id, name, link, month, year, day, location
									 from events 
									 inner join days on events.id=days.event_id  
									 where events.id=$1
									 order by days.day;`

	rows, _ := db.Query(sqlStatement, id)
	formatAndReturnJSONResponse(rows, w)
	db.Close()

	return
}

// GetEventsByMonth : by events by month
func GetEventsByMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	month := strings.ToLower(vars["month"])
	year := strings.ToLower(vars["year"])
	db := db.OpenDB()
	defer db.Close()
	sqlStatement := `select events.id, name, link, month, year, day, location
									 from events 
									 inner join days on events.id=days.event_id 
									 where events.month=$1 and events.year=$2
									 order by days.day;`
	rows, err := db.Query(sqlStatement, month, year)

	if err != nil {
		panic(err)
	}
	formatAndReturnJSONResponse(rows, w)
	db.Close()

	return
}

func formatAndReturnJSONResponse(rows *sql.Rows, w http.ResponseWriter) {

	defer rows.Close()
	events := []event.Event{}
	for rows.Next() {
		event := event.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Link,
			&event.Month,
			&event.Year,
			&event.Day,
			&event.Location,
		)

		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	out, err := json.Marshal(events)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows.Close()

	setHeaders(w)
	fmt.Fprintf(w, string(out))
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

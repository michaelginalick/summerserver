package eventsController

import (
	"database/sql"
	"github.com/lib/pq"
	"net/http"
	"../../../db"
	"../../../structs"
	"encoding/json"
	"fmt"
	"strconv"
)

func parseToIntSlice(e *[]int) *[]int {
	h := *e
	newDays := make([]int, 0)
	for i := 0; i < len(h); i++ {
		s, _ := strconv.Atoi(string(h[i]))
		newDays = append(newDays, s)
	}

	return &newDays
}


// GetEvents :  returns all the events
func GetEvents(w http.ResponseWriter, r *http.Request) {

	db := db.OpenDB()
	defer db.Close()
	sqlStatement := `select id, name, link, month, days, individual_days, festival_length from events;`
	rows, _ := db.Query(sqlStatement)
	defer rows.Close()

	

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	formatAndReturnJSONResponse(rows, w)

	db.Close()
	return
}


// GetEvent : by event by id
func GetEvent(w http.ResponseWriter, r *http.Request) {

}


// GetEventsByMonth : by events by month
func GetEventsByMonth(w http.ResponseWriter, r *http.Request) {

}


func formatAndReturnJSONResponse(rows *sql.Rows, w http.ResponseWriter) {
	events := []event.Event{}



	for rows.Next() {
		event := event.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Link,
			&event.Month,
			pq.Array(&event.Days),
			pq.Array(&event.IndividualDays),
			&event.FestivalLength,
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
	
	fmt.Fprintf(w, string(out))
}

package eventsController

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../../../db"
	"../../../structs"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
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
	sqlStatement := `select * from events;`
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
	sqlStatement := `select * from events where id = $1;`
	rows, _ := db.Query(sqlStatement, id)
	formatAndReturnJSONResponse(rows, w)
	db.Close()

	return
}

// GetEventsByMonth : by events by month
func GetEventsByMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	month := vars["month"]
	db := db.OpenDB()
	defer db.Close()
	sqlStatement := `select * from events where month=$1;`
	rows, err := db.Query(sqlStatement, month)

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

	rows.Close()

	setHeaders(w)
	fmt.Fprintf(w, string(out))
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

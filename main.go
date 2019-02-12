package main

import (
	"./routes/api/v1"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/events", eventscontroller.GetEvents).Methods("GET")
	router.HandleFunc("/events/{id}", eventscontroller.GetEventByID).Methods("GET")
	router.HandleFunc("/events_by_month/{month}/{year}", eventscontroller.GetEventsByMonth).Methods("GET")
	router.HandleFunc("/create_event", eventscontroller.CreateEvent).Methods("POST")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(allowedOrigins, allowedMethods)(router)))
}

package main

import (
	"./routes/api/v1"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/events", eventsController.GetEvents).Methods("GET")

	router.HandleFunc("/events/{id}", eventsController.GetEventByID).Methods("GET")
	router.HandleFunc("/eventsByMonth/{month}", eventsController.GetEventsByMonth).Methods("GET")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(allowedOrigins, allowedMethods)(router)))
}


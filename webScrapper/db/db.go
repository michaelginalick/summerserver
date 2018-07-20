package db

import (
	"database/sql"
	"fmt"
	"log"

	event "../structs"
	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "summerserver"
)

// OpenDB :  opens a connection to the database
func OpenDB() *sql.DB {
	var db *sql.DB

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

// SaveRecords :  saves event records to the database
func SaveRecords(e *event.Event) {
	db := OpenDB()

	defer db.Close()

	for i := e; i != nil; i = i.Next {
		if !isPresent(i, db) {
			_, err := db.Exec(
				"INSERT INTO  events (name, link, month, days, year, individual_days, festival_length) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				i.Name, i.Link, i.Month, pq.Array(i.Days), i.Year, pq.Array(i.IndividualDays), i.FestivalLength,
			)
			checkErr(err)
		}
	}
	db.Close()
}

func isPresent(e *event.Event, db *sql.DB) (exists bool) {

	sqlStatement := `select exists(select 1 from events where name = $1);`
	rows, err := db.Query(sqlStatement, e.Name)

	checkErr(err)

	for rows.Next() {
		err := rows.Scan(&exists)
		checkErr(err)
	}
	return exists
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

package db

import (
	"database/sql"
	"fmt"
	"log"

	event "../structs"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "michaelginalick"
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
		id := 0
		if !isPresent(i, db) {
			sqlStatement := `INSERT INTO  events (name, link, month, year, location) VALUES ($1, $2, $3, $4, $5) RETURNING id`
			err := db.QueryRow(sqlStatement, i.Name, i.Link, i.Month, i.Year, i.Location).Scan(&id)
			checkErr(err)
			for j := 0; j <= len(i.IndividualDays); j++ {
				_, err := db.Exec(
					"INSERT INTO days (day, event_id) VALUES ($1, $2)",
					i.IndividualDays[j], id,
				)
				checkErr(err)
			}
		}
	}
	db.Close()
}

func isPresent(e *event.Event, db *sql.DB) (exists bool) {

	sqlStatement := `select exists(select 1 from events where name = $1 and year = $2);`
	rows, err := db.Query(sqlStatement, e.Name, e.Year)

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

package db

import(
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "summerserver"
)



// OpenDB :  opens a connection to a database
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

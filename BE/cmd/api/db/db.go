package db

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "123"
	dbname   = "ChatApp"
)

var DB *sql.DB

func Init() {
	var err error

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	if err = DB.Ping(); err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return
	}

	fmt.Println("The database is connected")
}

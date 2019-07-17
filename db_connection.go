package main

import (
	"database/sql"
	"fmt"
)

func getDBConnection() *sql.DB {

	host := "localhost"
	port := 5433
	user := "user1"
	password := "pass1"
	dbname := "baseball"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}

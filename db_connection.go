package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

func getDBConnection() *sql.DB {

	host := os.Getenv("BASEBALL_HOST")
	port, _ := strconv.Atoi(os.Getenv("BASEBALL_PORT"))
	user := os.Getenv("BASEBALL_USER")
	password := os.Getenv("BASEBALL_PASS")
	dbname := os.Getenv("BASEBALL_DB")

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

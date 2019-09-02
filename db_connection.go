package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func getDBConnection() *sql.DB {

	host := os.Getenv("BASEBALL_HOST")

	log.Println("Host:", host)

	port, portErr := strconv.Atoi(os.Getenv("BASEBALL_PORT"))

	if portErr != nil {
		log.Println("Could not read port env var", portErr)
	} else {
		log.Println("Port:", port)
	}

	user := os.Getenv("BASEBALL_USER")
	log.Println("User:", user)
	password := os.Getenv("BASEBALL_PASS")
	dbname := os.Getenv("BASEBALL_DB")
	log.Println("dbname:", dbname)

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

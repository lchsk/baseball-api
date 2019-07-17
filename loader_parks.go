package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type RawPark struct {
	ParkID    string    `json:"park_id"`
	Name      string    `json:"name"`
	Nickname  string    `json:"nickname"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	League    string    `json:"league"`
}

func readRawPark(line []string) *RawPark {
	return &RawPark{
		ParkID:    line[0],
		Name:      line[1],
		Nickname:  line[2],
		City:      line[3],
		State:     line[4],
		StartDate: parseUSDate(line[5]),
		EndDate:   parseUSDate(line[6]),
		League:    line[7],
	}
}

func loadParks(path string) {
	csvFile, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var parks []*RawPark

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		parks = append(parks, readRawPark(line))
	}

	stmt := Statements["insertPark"]

	log.Println("Inserting parks")

	for _, park := range parks {
		_, err := stmt.Exec(park.ParkID, park.Name, park.Nickname, park.City, park.State, park.StartDate, park.EndDate, park.League)

		if err != nil {
			log.Printf("Error when inserting park %v %s", park, err)
		}
	}

	log.Printf("Inserted park")
}

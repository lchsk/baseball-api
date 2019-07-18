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

type RawPerson struct {
	PersonID     string    `json:"person_id"`
	LastName     string    `json:"last_name"`
	FirstName    string    `json:"first_name"`
	PlayerDebut  time.Time `json:"player_debut"`
	ManagerDebut time.Time `json:"manager_debut"`
	CoachDebut   time.Time `json:"coach_debut"`
	UmpireDebut  time.Time `json:"umpire_debut"`
}

func readRawPerson(line []string) *RawPerson {
	return &RawPerson{
		PersonID:     line[0],
		LastName:     line[1],
		FirstName:    line[2],
		PlayerDebut:  parseUSDate(line[3]),
		ManagerDebut: parseUSDate(line[4]),
		CoachDebut:   parseUSDate(line[5]),
		UmpireDebut:  parseUSDate(line[6]),
	}
}

func loadPeople(path string) {
	csvFile, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var people []*RawPerson

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		people = append(people, readRawPerson(line))
	}

	stmt := Statements["insertPerson"]

	log.Println("Inserting people")

	for _, person := range people {
		_, err := stmt.Exec(person.PersonID, person.LastName, person.FirstName, person.PlayerDebut, person.ManagerDebut, person.CoachDebut, person.UmpireDebut)

		if err != nil {
			log.Printf("Error when inserting person %v %s", person, err)
		}
	}

	log.Printf("Inserted people")
}

package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/lchsk/baseballapi/dbconnection"
	_ "github.com/lib/pq"
)

type RawTeam struct {
	TeamSymbol string
	Founded    int
	League     string
	Location   string
	Name       string
}

func readRawTeam(line []string) *RawTeam {
	return &RawTeam{
		TeamSymbol: line[0],
		League:     line[1],
		Location:   line[2],
		Name:       line[3],
		Founded:    parseInt(line[4]),
	}
}

func loadTeams(path string) {
	csvFile, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var teams []*RawTeam

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		teams = append(teams, readRawTeam(line))
	}

	stmt := dbconnection.Statements["insertTeam"]

	log.Println("Inserting teams")

	for _, team := range teams {
		_, err := stmt.Exec(team.TeamSymbol, team.Founded, team.League, team.Location, team.Name)

		if err != nil {
			log.Printf("Error when inserting team %v %s", team, err)
		}
	}

	log.Printf("Inserted teams")
}

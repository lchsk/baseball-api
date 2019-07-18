package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

// TODO: Put them in a struct
var db *sql.DB
var TEAMS map[string]*RawTeam
var PARKS map[string]*RawPark

// add statements

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		w.Header().Add("Content-Type", "application/json")

		log.Println(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func loadTeamsDataFromDB() {
	stmt := Statements["selectAllTeams"]

	rows, err := stmt.Query()

	if err != nil {
		panic(err)
	}

	TEAMS = make(map[string]*RawTeam)

	for rows.Next() {
		var team RawTeam

		rows.Scan(
			&team.TeamSymbol,
			&team.Founded,
			&team.League,
			&team.Location,
			&team.Name,
		)

		TEAMS[team.TeamSymbol] = &team
	}
}

func loadParksDataFromDB() {
	stmt := Statements["selectAllParks"]

	rows, err := stmt.Query()

	if err != nil {
		panic(err)
	}

	PARKS = make(map[string]*RawPark)

	for rows.Next() {
		var park RawPark

		rows.Scan(
			&park.ParkID,
			&park.Name,
			&park.Nickname,
			&park.City,
			&park.State,
			&park.StartDate,
			&park.EndDate,
			&park.League,
		)

		PARKS[park.ParkID] = &park
	}
}

func serveAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/games/{date}/{teams}", getGameSummary).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/teams/{team}", getTeam).Methods(http.MethodGet)

	log.Println("Serving api")
	log.Fatal(http.ListenAndServe(":8000", commonMiddleware(router)))
}

func main() {
	var loadData = flag.Bool("load-data", false, "Load game log data")
	var gameLogsDir = flag.String("game-logs", "", "Path to game logs directory")
	var teamsFile = flag.String("teams", "", "Path to teams file")
	var parksFile = flag.String("parks", "", "Path to parks file")
	var peopleFile = flag.String("people", "", "Path to people file")

	flag.Parse()

	db = getDBConnection()
	prepareQueries(db)

	loadTeamsDataFromDB()
	loadParksDataFromDB()

	if *loadData {
		if *gameLogsDir != "" {
			loadGameLogs(*gameLogsDir)
		}

		if *teamsFile != "" {
			loadTeams(*teamsFile)
		}

		if *parksFile != "" {
			loadParks(*parksFile)
		}

		if *peopleFile != "" {
			loadPeople(*peopleFile)
		}

		return
	}

	serveAPI()
}

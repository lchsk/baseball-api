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
	"github.com/lchsk/baseballapi/dbconnection"
)

// TODO: Put them in a struct
var db *sql.DB
var TEAMS map[string]*RawTeam

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		w.Header().Add("Content-Type", "application/json")

		log.Println(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func loadTeamsDataFromDB() {
	stmt := dbconnection.Statements["selectAllTeams"]

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

	flag.Parse()

	db = dbconnection.GetDBConnection()
	dbconnection.PrepareQueries(db)

	loadTeamsDataFromDB()

	if *loadData {
		if *gameLogsDir != "" {
			loadGameLogs(*gameLogsDir)
		}

		if *teamsFile != "" {
			loadTeams(*teamsFile)
		}

		return
	}

	serveAPI()
}

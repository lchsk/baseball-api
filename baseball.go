package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gorilla/mux"
)

// TODO: Put them in a struct
var db *sql.DB
var TEAMS map[string]*RawTeam
var PARKS map[string]*RawPark
var PositionSymbolsMap = make(map[int]string)
var PositionNamesMap = make(map[int]string)

// add statements

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
	router.HandleFunc("/api/v1/games/{date}/{teams}/lineups", getGameSummaryLineups).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/games/{date}/{teams}/stats", getGameSummaryStats).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/teams/{team}", getTeam).Methods(http.MethodGet)

	log.Println("Serving api")
	log.Fatal(http.ListenAndServe(":8000", commonMiddleware(router)))
}

func initPositionConstants() {
	PositionSymbolsMap[1] = "P"
	PositionSymbolsMap[2] = "C"
	PositionSymbolsMap[3] = "1B"
	PositionSymbolsMap[4] = "2B"
	PositionSymbolsMap[5] = "3B"
	PositionSymbolsMap[6] = "SS"
	PositionSymbolsMap[7] = "LF"
	PositionSymbolsMap[8] = "CF"
	PositionSymbolsMap[9] = "RF"
	PositionSymbolsMap[10] = "DH"

	PositionNamesMap[1] = "pitcher"
	PositionNamesMap[2] = "catcher"
	PositionNamesMap[3] = "first baseman"
	PositionNamesMap[4] = "second baseman"
	PositionNamesMap[5] = "third baseman"
	PositionNamesMap[6] = "shortstop"
	PositionNamesMap[7] = "left fielder"
	PositionNamesMap[8] = "center fielder"
	PositionNamesMap[9] = "right fielder"
	PositionNamesMap[10] = "designated hitter"
}

func main() {
	lumberjackLog := &lumberjack.Logger{
		Filename:   "./baseball.log",
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	log.SetOutput(io.MultiWriter(lumberjackLog, os.Stderr))

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

	initPositionConstants()

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

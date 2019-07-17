package main

import (
	"database/sql"
	"encoding/json"
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

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// URL?
}

type GameSummaryTeam struct {
	Symbol       string `json:"symbol"`
	FullTeamName string `json:"full_team_name"`
	TeamName     string `json:"team_name"`
	TeamLocation string `json:"team_location"`
	League       string `json:"league"`
	GameNumber   int    `json:"game_number"`
	Score        int    `json:"runs"`
	Hits         int    `json:"hits"`
	Errors       int    `json:"errors"`
	Manager      Person `json:"manager"`
}

type GameSummary struct {
	Date string `json:"date"`
	// TODO: Consider different field name, double_header_information?
	NumberOfGame         string          `json:"number_of_game"`
	DayOfWeek            string          `json:"day_of_week"`
	VisitingTeam         GameSummaryTeam `json:"visiting_team"`
	HomeTeam             GameSummaryTeam `json:"home_team"`
	GameLengthInOuts     int             `json:"game_length_in_outs"`
	TimeOfGameInMins     int             `json:"game_length_in_mins"`
	DayNightIndicator    string          `json:"day_night_indicator"`
	Attendance           int             `json:"attendance"`
	WinningPitcher       Person          `json:"winning_pitcher"`
	LosingPitcher        Person          `json:"losing_pitcher"`
	SavingPitcher        Person          `json:"saving_pitcher"`
	GameWinningRBIBatter Person          `json:"game_winning_rbi_batter"`
	// park struct
	// line score
	// team names
}

type GameSummaryResponse struct {
	Games []GameSummary `json:"games"`
}

type Error struct {
	Message string `json:"message"`
}

type GameSummaryError struct {
	Errors []Error `json:"errors"`
}

func getGameSummary(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	date := params["date"]
	teams := strings.Split(params["teams"], "@")

	if len(teams) != 2 {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(GameSummaryError{
			Errors: []Error{{Message: "Must provide two teams"}},
		})
		return
	}

	visitingTeam := teams[0]
	homeTeam := teams[1]

	games, err := loadGames(date, visitingTeam, homeTeam)

	if err != nil {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(GameSummaryError{
			Errors: []Error{{Message: "Invalid input"}},
		})
		return
	}
	if len(games) == 0 {
		w.WriteHeader(404)

		json.NewEncoder(w).Encode(GameSummaryError{
			Errors: []Error{{Message: "No games were found"}},
		})
		return
	}

	var data []GameSummary

	for _, game := range games {
		// TODO: Move it away from here
		visitingFullTeamName := ""
		visitingTeamName := ""
		visitingTeamLocation := ""
		visitingTeamData, ok := TEAMS[game.VisitingTeam]

		if ok {
			visitingFullTeamName = fmt.Sprintf("%s %s", visitingTeamData.Location, visitingTeamData.Name)
			visitingTeamLocation = visitingTeamData.Location
			visitingTeamName = visitingTeamData.Name
		}

		homeFullTeamName := ""
		homeTeamName := ""
		homeTeamLocation := ""
		homeTeamData, ok := TEAMS[game.HomeTeam]

		if ok {
			homeFullTeamName = fmt.Sprintf("%s %s", homeTeamData.Location, homeTeamData.Name)
			homeTeamLocation = homeTeamData.Location
			homeTeamName = homeTeamData.Name
		}

		data = append(data, GameSummary{
			Date:         game.Date.Format("2006-01-02"),
			NumberOfGame: game.NumberOfGame,
			DayOfWeek:    game.DayOfWeek,
			VisitingTeam: GameSummaryTeam{
				Symbol:       game.VisitingTeam,
				TeamName:     visitingTeamName,
				TeamLocation: visitingTeamLocation,
				FullTeamName: visitingFullTeamName,
				League:       game.VisitingTeamLeague,
				GameNumber:   game.VisitingGameNumber,
				Score:        game.VisitingTeamScore,
				Hits:         game.VisitingH,
				Errors:       game.VisitingErrors,
				Manager: Person{
					ID:   game.VisitingManagerID,
					Name: game.VisitingManagerName,
				},
			},
			HomeTeam: GameSummaryTeam{
				Symbol:       game.HomeTeam,
				TeamName:     homeTeamName,
				TeamLocation: homeTeamLocation,
				FullTeamName: homeFullTeamName,
				League:       game.HomeTeamLeague,
				GameNumber:   game.HomeTeamGameNumber,
				Score:        game.HomeTeamScore,
				Hits:         game.HomeH,
				Errors:       game.HomeErrors,
				Manager: Person{
					ID:   game.HomeManagerID,
					Name: game.HomeManagerName,
				},
			},
			GameLengthInOuts:  game.GameLengthInOuts,
			TimeOfGameInMins:  game.TimeOfGameInMins,
			DayNightIndicator: game.DayNightIndicator,
			Attendance:        game.Attendance,
			WinningPitcher: Person{
				ID:   game.WinningPitcherID,
				Name: game.WinningPitcherName,
			},
			LosingPitcher: Person{
				ID:   game.LosingPitcherID,
				Name: game.LosingPitcherName,
			},
			SavingPitcher: Person{
				ID:   game.SavingPitcherID,
				Name: game.SavingPitcherName,
			},
			GameWinningRBIBatter: Person{
				ID:   game.GameWinningRBIBatterID,
				Name: game.GameWinningRBIBatterName,
			},
		})
	}

	json.NewEncoder(w).Encode(GameSummaryResponse{
		Games: data,
	})
}

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

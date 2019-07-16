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
	"github.com/lchsk/baseball/dbconnection"
)

var db *sql.DB

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// URL?
}

type GameSummaryTeam struct {
	// Team name
	Symbol     string `json:"symbol"`
	League     string `json:"league"`
	GameNumber int    `json:"game_number"`
	Score      int    `json:"score"`
	Hits       int    `json:"hits"`
	Errors     int    `json:"errors"`
	Manager    Person `json:"manager"`
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
		data = append(data, GameSummary{
			Date:         game.Date.Format("2006-01-02"),
			NumberOfGame: game.NumberOfGame,
			DayOfWeek:    game.DayOfWeek,
			VisitingTeam: GameSummaryTeam{
				Symbol:     game.VisitingTeam,
				League:     game.VisitingTeamLeague,
				GameNumber: game.VisitingGameNumber,
				Score:      game.VisitingTeamScore,
				Hits:       game.VisitingH,
				Errors:     game.VisitingErrors,
				Manager: Person{
					ID:   game.VisitingManagerID,
					Name: game.VisitingManagerName,
				},
			},
			HomeTeam: GameSummaryTeam{
				Symbol:     game.HomeTeam,
				League:     game.HomeTeamLeague,
				GameNumber: game.HomeTeamGameNumber,
				Score:      game.HomeTeamScore,
				Hits:       game.HomeH,
				Errors:     game.HomeErrors,
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

func serveAPI() {
	db = dbconnection.GetDBConnection()
	dbconnection.PrepareQueries(db)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/games/{date}/{teams}", getGameSummary).Methods(http.MethodGet)

	log.Println("Serving api")
	log.Fatal(http.ListenAndServe(":8000", commonMiddleware(router)))
}

func main() {
	var loadData = flag.Bool("load-data", false, "Load game log data")
	var gameLogsDir = flag.String("game-logs", "./", "Path to game logs directory")

	flag.Parse()

	if *loadData {
		loadGameLogs(*gameLogsDir)
		return
	}

	serveAPI()
}

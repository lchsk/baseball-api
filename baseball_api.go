package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/lchsk/baseball/dbconnection"
)

var db *sql.DB

const ResponseStatusSuccess = "success"

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
	ID   string `json:"id"`
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
	Status string        `json:"status"`
	Data   []GameSummary `json:"data"`
}

func getGameSummary(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	gameID := params["gameId"]

	gameIdentifiers := strings.Split(gameID, "-")

	// TODO: Add error handling
	if len(gameIdentifiers) != 5 {
		panic("invalid game identifier")
	}

	visitingTeam := gameIdentifiers[0]
	homeTeam := gameIdentifiers[1]
	year := gameIdentifiers[2]
	month := gameIdentifiers[3]
	day := gameIdentifiers[4]

	date := fmt.Sprintf("%s-%s-%s", year, month, day)

	games := loadGames(date, visitingTeam, homeTeam)

	var data []GameSummary

	for _, game := range games {
		uniqueGameID := fmt.Sprintf("%s-%s", gameID, game.NumberOfGame)
		data = append(data, GameSummary{
			ID:           uniqueGameID,
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
		Status: ResponseStatusSuccess,
		Data:   data,
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

func main() {
	db = dbconnection.GetDBConnection()
	dbconnection.PrepareQueries(db)

	router := mux.NewRouter()
	router.HandleFunc("/api/games/{gameId}", getGameSummary).Methods(http.MethodGet)

	log.Println("Serving api")
	log.Fatal(http.ListenAndServe(":8000", commonMiddleware(router)))
}

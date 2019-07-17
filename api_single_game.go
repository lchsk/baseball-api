package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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

type GameSummaryPark struct {
	// Add URL
	ParkID string `json:"venue_id"`
	Name   string `json:"name"`
	City   string `json:"city"`
	State  string `json:"state"`
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
	Park                 GameSummaryPark `json:"venue"`
	// line score
	// team names
}

type GameSummaryResponse struct {
	Games []GameSummary `json:"games"`
}

func getGameSummary(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	date := params["date"]
	teams := strings.Split(params["teams"], "@")

	if len(teams) != 2 {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(ResponseErrors{
			Errors: []Error{{Message: "Must provide two teams"}},
		})
		return
	}

	visitingTeam := teams[0]
	homeTeam := teams[1]

	games, err := loadGames(date, visitingTeam, homeTeam)

	if err != nil {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(ResponseErrors{
			Errors: []Error{{Message: "Invalid input"}},
		})
		return
	}
	if len(games) == 0 {
		w.WriteHeader(404)

		json.NewEncoder(w).Encode(ResponseErrors{
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

		parkName := ""
		parkCity := ""
		parkState := ""
		park, ok := PARKS[game.ParkID]

		if ok {
			parkName = park.Name
			parkCity = park.City
			parkState = park.State
		}

		data = append(data, GameSummary{
			Date:         game.Date.Format("2006-01-02"),
			NumberOfGame: game.NumberOfGame,
			DayOfWeek:    game.DayOfWeek,
			Park: GameSummaryPark{
				ParkID: game.ParkID,
				Name:   parkName,
				City:   parkCity,
				State:  parkState,
			},
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

package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type GameLineupTeam struct {
	Manager         Person   `json:"manager"`
	TeamSymbol      string   `json:"team_symbol"`
	FullTeamName    string   `json:"full_team_name"`
	TeamName        string   `json:"team_name"`
	TeamLocation    string   `json:"team_location"`
	StartingPitcher Person   `json:"starting_pitcher"`
	StartingLineup  []Player `json:"starting_lineup"`
}

type Umpires struct {
	HomePlate  Person `json:"home_plate"`
	FirstBase  Person `json:"first_base"`
	SecondBase Person `json:"second_base"`
	ThirdBase  Person `json:"third_base"`
	LeftField  Person `json:"left_field"`
	RightField Person `json:"right_field"`
}

type GameLineup struct {
	Date         string         `json:"date"`
	NumberOfGame string         `json:"number_of_game"`
	VisitingTeam GameLineupTeam `json:"visiting_team"`
	HomeTeam     GameLineupTeam `json:"home_team"`
	Umpires      Umpires        `json:"umpires"`
}

type LineupsResponse struct {
	Games []GameLineup `json:"games"`
}

func getGameSummaryLineups(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	date := params["date"]
	teams := strings.Split(params["teams"], "@")

	// TODO: Factor out validation for game summary endpoints
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

	var data []GameLineup

	for _, game := range games {
		visitingTeamNameData := getTeamNameData(game.VisitingTeam)
		homeTeamNameData := getTeamNameData(game.HomeTeam)

		data = append(data, GameLineup{
			Date:         game.Date.Format("2006-01-02"),
			NumberOfGame: game.NumberOfGame,
			VisitingTeam: GameLineupTeam{
				Manager: Person{
					ID:   game.VisitingManagerID,
					Name: game.VisitingManagerName,
				},
				StartingPitcher: Person{
					ID:   game.VisitingStartingPitcherID,
					Name: game.VisitingStartingPitcherName,
				},
				StartingLineup: getVisitingBattingOrder(&game),
				TeamName:       visitingTeamNameData.Name,
				FullTeamName:   visitingTeamNameData.FullName,
				TeamSymbol:     visitingTeamNameData.Symbol,
				TeamLocation:   visitingTeamNameData.Location,
			},
			HomeTeam: GameLineupTeam{
				Manager: Person{
					ID:   game.HomeManagerID,
					Name: game.HomeManagerName,
				},
				StartingPitcher: Person{
					ID:   game.HomeStartingPitcherID,
					Name: game.HomeStartingPitcherName,
				},
				StartingLineup: getHomeBattingOrder(&game),
				TeamName:       homeTeamNameData.Name,
				FullTeamName:   homeTeamNameData.FullName,
				TeamSymbol:     homeTeamNameData.Symbol,
				TeamLocation:   homeTeamNameData.Location,
			},
			Umpires: Umpires{
				HomePlate: Person{
					ID:   game.HomePlateUmpireID,
					Name: game.HomePlateUmpireName,
				},
				FirstBase: Person{
					ID:   game.FirstBaseUmpireID,
					Name: game.FirstBaseUmpireName,
				},
				SecondBase: Person{
					ID:   game.SecondBaseUmpireID,
					Name: game.SecondBaseUmpireName,
				},
				ThirdBase: Person{
					ID:   game.ThirdBaseUmpireID,
					Name: game.ThirdBaseUmpireName,
				},
				LeftField: Person{
					ID:   game.LeftFieldUmpireID,
					Name: game.LeftFieldUmpireName,
				},
				RightField: Person{
					ID:   game.RightFieldUmpireID,
					Name: game.RightFieldUmpireName,
				},
			},
		})
	}

	json.NewEncoder(w).Encode(LineupsResponse{
		Games: data,
	})
}

func getVisitingBattingOrder(game *Game) []Player {
	return []Player{
		Player{
			ID:             game.VisitingPlayer1ID,
			Name:           game.VisitingPlayer1Name,
			PositionNumber: game.VisitingPlayer1Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer1Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer1Position],
		},
		Player{
			ID:             game.VisitingPlayer2ID,
			Name:           game.VisitingPlayer2Name,
			PositionNumber: game.VisitingPlayer2Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer2Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer2Position],
		},
		Player{
			ID:             game.VisitingPlayer3ID,
			Name:           game.VisitingPlayer3Name,
			PositionNumber: game.VisitingPlayer3Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer3Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer3Position],
		},
		Player{
			ID:             game.VisitingPlayer4ID,
			Name:           game.VisitingPlayer4Name,
			PositionNumber: game.VisitingPlayer4Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer4Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer4Position],
		},
		Player{
			ID:             game.VisitingPlayer5ID,
			Name:           game.VisitingPlayer5Name,
			PositionNumber: game.VisitingPlayer5Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer5Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer5Position],
		},
		Player{
			ID:             game.VisitingPlayer6ID,
			Name:           game.VisitingPlayer6Name,
			PositionNumber: game.VisitingPlayer6Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer6Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer6Position],
		},
		Player{
			ID:             game.VisitingPlayer7ID,
			Name:           game.VisitingPlayer7Name,
			PositionNumber: game.VisitingPlayer7Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer7Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer7Position],
		},
		Player{
			ID:             game.VisitingPlayer8ID,
			Name:           game.VisitingPlayer8Name,
			PositionNumber: game.VisitingPlayer8Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer8Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer8Position],
		},
		Player{
			ID:             game.VisitingPlayer9ID,
			Name:           game.VisitingPlayer9Name,
			PositionNumber: game.VisitingPlayer9Position,
			PositionName:   PositionNamesMap[game.VisitingPlayer9Position],
			PositionSymbol: PositionSymbolsMap[game.VisitingPlayer9Position],
		},
	}
}

func getHomeBattingOrder(game *Game) []Player {
	return []Player{
		Player{
			ID:             game.HomePlayer1ID,
			Name:           game.HomePlayer1Name,
			PositionNumber: game.HomePlayer1Position,
			PositionName:   PositionNamesMap[game.HomePlayer1Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer1Position],
		},
		Player{
			ID:             game.HomePlayer2ID,
			Name:           game.HomePlayer2Name,
			PositionNumber: game.HomePlayer2Position,
			PositionName:   PositionNamesMap[game.HomePlayer2Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer2Position],
		},
		Player{
			ID:             game.HomePlayer3ID,
			Name:           game.HomePlayer3Name,
			PositionNumber: game.HomePlayer3Position,
			PositionName:   PositionNamesMap[game.HomePlayer3Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer3Position],
		},
		Player{
			ID:             game.HomePlayer4ID,
			Name:           game.HomePlayer4Name,
			PositionNumber: game.HomePlayer4Position,
			PositionName:   PositionNamesMap[game.HomePlayer4Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer4Position],
		},
		Player{
			ID:             game.HomePlayer5ID,
			Name:           game.HomePlayer5Name,
			PositionNumber: game.HomePlayer5Position,
			PositionName:   PositionNamesMap[game.HomePlayer5Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer5Position],
		},
		Player{
			ID:             game.HomePlayer6ID,
			Name:           game.HomePlayer6Name,
			PositionNumber: game.HomePlayer6Position,
			PositionName:   PositionNamesMap[game.HomePlayer6Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer6Position],
		},
		Player{
			ID:             game.HomePlayer7ID,
			Name:           game.HomePlayer7Name,
			PositionNumber: game.HomePlayer7Position,
			PositionName:   PositionNamesMap[game.HomePlayer7Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer7Position],
		},
		Player{
			ID:             game.HomePlayer8ID,
			Name:           game.HomePlayer8Name,
			PositionNumber: game.HomePlayer8Position,
			PositionName:   PositionNamesMap[game.HomePlayer8Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer8Position],
		},
		Player{
			ID:             game.HomePlayer9ID,
			Name:           game.HomePlayer9Name,
			PositionNumber: game.HomePlayer9Position,
			PositionName:   PositionNamesMap[game.HomePlayer9Position],
			PositionSymbol: PositionSymbolsMap[game.HomePlayer9Position],
		},
	}
}

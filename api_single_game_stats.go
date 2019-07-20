package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Pitching struct {
	PitchersUsed         int `json:"pitchers_used"`
	IndividualEarnedRuns int `json:"individual_earned_runs"`
	TeamEarnedRuns       int `json:"team_earned_runs"`
	WildPitches          int `json:"wild_pitches"`
	Balks                int `json:"balks"`
}

type Batting struct {
	Runs                   int `json:"runs"`
	AtBats                 int `json:"at_bats"`
	Hits                   int `json:"hits"`
	Doubles                int `json:"doubles"`
	Triples                int `json:"triples"`
	HomeRuns               int `json:"home_runs"`
	RunsBattedIn           int `json:"runs_batted_in"`
	SacrificeHits          int `json:"sacrifice_hits"`
	SacrificeFlies         int `json:"sacrifice_flies"`
	HitByPitch             int `json:"hit_by_pitch"`
	Walks                  int `json:"walks"`
	IntentionalWalks       int `json:"intentional_walks"`
	Strikeouts             int `json:"strikeouts"`
	StolenBases            int `json:"stolen_bases"`
	CaughtStealing         int `json:"caught_stealing"`
	GroundedIntoDoublePlay int `json:"grounded_into_double_play"`
	CatcherInterference    int `json:"catcher_interference"`
	LeftOnBase             int `json:"left_on_base"`
}

type Fielding struct {
	Putouts     int `json:"putouts"`
	Assists     int `json:"assists"`
	Errors      int `json:"errors"`
	PassedBalls int `json:"passed_balls"`
	DoublePlays int `json:"double_plays"`
	TriplePlays int `json:"triple_plays"`
}

type GameStatsTeam struct {
	// TODO: Move out the naming data
	TeamSymbol   string   `json:"team_symbol"`
	FullTeamName string   `json:"full_team_name"`
	TeamName     string   `json:"team_name"`
	TeamLocation string   `json:"team_location"`
	Pitching     Pitching `json:"pitching"`
	Batting      Batting  `json:"batting"`
	Fielding     Fielding `json:"fielding"`
}
type GameSummaryStats struct {
	Date         string        `json:"date"`
	NumberOfGame string        `json:"number_of_game"`
	VisitingTeam GameStatsTeam `json:"visiting_team"`
	HomeTeam     GameStatsTeam `json:"home_team"`
}

type StatsResponse struct {
	Games []GameSummaryStats `json:"games"`
}

func getGameSummaryStats(w http.ResponseWriter, req *http.Request) {
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

	var data []GameSummaryStats

	for _, game := range games {
		visitingTeamNameData := getTeamNameData(game.VisitingTeam)
		homeTeamNameData := getTeamNameData(game.HomeTeam)

		data = append(data, GameSummaryStats{
			Date:         game.Date.Format("2006-01-02"),
			NumberOfGame: game.NumberOfGame,
			VisitingTeam: GameStatsTeam{
				TeamName:     visitingTeamNameData.Name,
				FullTeamName: visitingTeamNameData.FullName,
				TeamSymbol:   visitingTeamNameData.Symbol,
				TeamLocation: visitingTeamNameData.Location,
				Pitching: Pitching{
					PitchersUsed:         game.VisitingPitchersUsed,
					IndividualEarnedRuns: game.VisitingIndividualEarnedRuns,
					TeamEarnedRuns:       game.VisitingTeamEarnedRuns,
					WildPitches:          game.VisitingWildPitches,
					Balks:                game.VisitingBalks,
				},
				Fielding: Fielding{
					Putouts:     game.VisitingPutouts,
					Assists:     game.VisitingAssists,
					Errors:      game.VisitingErrors,
					PassedBalls: game.VisitingPassedBalls,
					DoublePlays: game.VisitingDoublePlays,
					TriplePlays: game.VisitingTriplePlays,
				},
				Batting: Batting{
					Runs:                   game.VisitingTeamScore,
					AtBats:                 game.VisitingAB,
					Hits:                   game.VisitingH,
					Doubles:                game.Visiting2B,
					Triples:                game.Visiting3B,
					HomeRuns:               game.VisitingHR,
					RunsBattedIn:           game.VisitingRBI,
					SacrificeHits:          game.VisitingSH,
					SacrificeFlies:         game.VisitingSF,
					HitByPitch:             game.VisitingHBP,
					Walks:                  game.VisitingBB,
					IntentionalWalks:       game.VisitingIBB,
					Strikeouts:             game.VisitingK,
					StolenBases:            game.VisitingSB,
					CaughtStealing:         game.VisitingCS,
					GroundedIntoDoublePlay: game.VisitingGIDP,
					CatcherInterference:    game.VisitingCI,
					LeftOnBase:             game.VisitingLOB,
				},
			},
			HomeTeam: GameStatsTeam{
				TeamName:     homeTeamNameData.Name,
				FullTeamName: homeTeamNameData.FullName,
				TeamSymbol:   homeTeamNameData.Symbol,
				TeamLocation: homeTeamNameData.Location,
				Pitching: Pitching{
					PitchersUsed:         game.HomePitchersUsed,
					IndividualEarnedRuns: game.HomeIndividualEarnedRuns,
					TeamEarnedRuns:       game.HomeTeamEarnedRuns,
					WildPitches:          game.HomeWildPitches,
					Balks:                game.HomeBalks,
				},
				Fielding: Fielding{
					Putouts:     game.HomePutouts,
					Assists:     game.HomeAssists,
					Errors:      game.HomeErrors,
					PassedBalls: game.HomePassedBalls,
					DoublePlays: game.HomeDoublePlays,
					TriplePlays: game.HomeTriplePlays,
				},
				Batting: Batting{
					Runs:                   game.HomeTeamScore,
					AtBats:                 game.HomeAB,
					Hits:                   game.HomeH,
					Doubles:                game.Home2B,
					Triples:                game.Home3B,
					HomeRuns:               game.HomeHR,
					RunsBattedIn:           game.HomeRBI,
					SacrificeHits:          game.HomeSH,
					SacrificeFlies:         game.HomeSF,
					HitByPitch:             game.HomeHBP,
					Walks:                  game.HomeBB,
					IntentionalWalks:       game.HomeIBB,
					Strikeouts:             game.HomeK,
					StolenBases:            game.HomeSB,
					CaughtStealing:         game.HomeCS,
					GroundedIntoDoublePlay: game.HomeGIDP,
					CatcherInterference:    game.HomeCI,
					LeftOnBase:             game.HomeLOB,
				},
			},
		})
	}

	json.NewEncoder(w).Encode(StatsResponse{
		Games: data,
	})
}

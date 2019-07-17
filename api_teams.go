package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type TeamDataResponse struct {
	Team TeamDataSummary `json:"team"`
}

type TeamDataSummary struct {
	TeamSymbol string `json:"team_symbol"`
	Founded    int    `json:"founded"`
	League     string `json:"league"`
	Location   string `json:"location"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
}

func getTeam(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	teamSymbol := params["team"]

	teamData, ok := TEAMS[teamSymbol]

	if !ok {
		w.WriteHeader(404)

		json.NewEncoder(w).Encode(ResponseErrors{
			Errors: []Error{{Message: "There is no team with that symbol"}},
		})
		return
	}

	json.NewEncoder(w).Encode(
		TeamDataResponse{
			Team: TeamDataSummary{
				TeamSymbol: teamData.TeamSymbol,
				Founded:    teamData.Founded,
				League:     teamData.League,
				Location:   teamData.Location,
				Name:       teamData.Name,
				FullName:   fmt.Sprintf("%s %s", teamData.Location, teamData.Name),
			},
		},
	)
}

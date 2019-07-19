package main

import "fmt"

type TeamNameData struct {
	FullName string
	Symbol   string
	Name     string
	Location string
}

func getTeamNameData(teamSymbol string) *TeamNameData {
	teamData, ok := TEAMS[teamSymbol]

	if ok {
		return &TeamNameData{
			Symbol:   teamSymbol,
			Name:     teamData.Name,
			Location: teamData.Location,
			FullName: fmt.Sprintf("%s %s", teamData.Location, teamData.Name),
		}
	}

	return &TeamNameData{}
}

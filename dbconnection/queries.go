package dbconnection

import "database/sql"

const selectGameByDate = `select * from game where visiting_team = $1 and home_team = $2 and game_date = $3`
const selectAllTeams = `select * from team`
const insertTeam = `insert into team (team_symbol, founded, league, location, name) values ($1, $2, $3, $4, $5)`

var Statements = make(map[string]*sql.Stmt)

func PrepareQueries(db *sql.DB) {
	stmtSelectGameByDate, _ := db.Prepare(selectGameByDate)
	Statements["selectGameByDate"] = stmtSelectGameByDate

	stmtSelectAllTeams, _ := db.Prepare(selectAllTeams)
	Statements["selectAllTeams"] = stmtSelectAllTeams

	stmtInsertTeam, _ := db.Prepare(insertTeam)
	Statements["insertTeam"] = stmtInsertTeam
}

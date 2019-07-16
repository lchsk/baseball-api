package dbconnection

import "database/sql"

const selectGameByDate = `select * from game where visiting_team = $1 and home_team = $2 and game_date = $3 `

var Statements = make(map[string]*sql.Stmt)

func PrepareQueries(db *sql.DB) {
	stmtSelectGameByDate, err := db.Prepare(selectGameByDate)

	if err != nil {
		panic(err)
	}

	Statements["selectGameByDate"] = stmtSelectGameByDate
}

package main

import "database/sql"

const selectGameByDate = `select * from game where visiting_team = $1 and home_team = $2 and game_date = $3`
const selectAllTeams = `select * from team`
const selectAllParks = `select * from park`
const insertTeam = `insert into team (team_symbol, founded, league, location, name) values ($1, $2, $3, $4, $5)`
const insertPark = `insert into park (park_id, name, nickname, city, state, start_date, end_date, league) values ($1, $2, $3, $4, $5, $6, $7, $8)`
const queryInsertPerson = `insert into person (person_id, last_name, first_name, player_debut, manager_debut, coach_debut, umpire_debut) values ($1, $2, $3, $4, $5, $6, $7)`

var Statements = make(map[string]*sql.Stmt)

func prepareQueries(db *sql.DB) {
	stmtSelectGameByDate, _ := db.Prepare(selectGameByDate)
	Statements["selectGameByDate"] = stmtSelectGameByDate

	stmtSelectAllTeams, _ := db.Prepare(selectAllTeams)
	Statements["selectAllTeams"] = stmtSelectAllTeams

	stmtSelectAllParks, _ := db.Prepare(selectAllParks)
	Statements["selectAllParks"] = stmtSelectAllParks

	stmtInsertTeam, _ := db.Prepare(insertTeam)
	Statements["insertTeam"] = stmtInsertTeam

	stmtInsertPark, _ := db.Prepare(insertPark)
	Statements["insertPark"] = stmtInsertPark

	stmtInsertGame, _ := db.Prepare(queryInsertGame)
	Statements["insertGame"] = stmtInsertGame

	stmtInsertPerson, _ := db.Prepare(queryInsertPerson)
	Statements["insertPerson"] = stmtInsertPerson
}

const queryInsertGame = `INSERT INTO game (
	game_date,
	number_of_game,
	day_of_week,
	visiting_team,
	visiting_team_league,
	visiting_game_number,
	home_team,
	home_team_league,
	home_team_game_number,
	visiting_team_score,
	home_team_score,
	game_length_in_outs,
	day_night_indicator,
	completion_information,
	forfeit_information,
	protest_information,
	park_id,
	attendance,
	time_of_game_in_mins,
	visiting_line_score,
	home_line_score,
	visiting_ab   ,
	visiting_h    ,
	visiting_2B   ,
	visiting_3B   ,
	visiting_hr   ,
	visiting_rbi  ,
	visiting_sh   ,
	visiting_sf   ,
	visiting_hbp  ,
	visiting_bb   ,
	visiting_ibb  ,
	visiting_k    ,
	visiting_sb   ,
	visiting_cs   ,
	visiting_gidp ,
	visiting_ci   ,
	visiting_lob  ,
	visiting_pitchers_used         ,
	visiting_individual_earned_runs ,
	visiting_team_earned_runs       ,
	visiting_wild_pitches          ,
	visiting_balks                 ,
	visiting_putouts     ,
	visiting_assists     ,
	visiting_errors      ,
	visiting_passed_balls ,
	visiting_double_plays ,
	visiting_triple_plays ,
	home_ab   ,
	home_h    ,
	home_2B   ,
	home_3B   ,
	home_hr   ,
	home_rbi  ,
	home_sh   ,
	home_sf   ,
	home_hbp  ,
	home_bb   ,
	home_ibb  ,
	home_k    ,
	home_sb   ,
	home_cs   ,
	home_gidp ,
	home_ci   ,
	home_lob  ,
	home_pitchers_used         ,
	home_individual_earned_runs ,
	home_team_earned_runs       ,
	home_wild_pitches          ,
	home_balks          ,
	home_putouts     ,
	home_assists     ,
	home_errors      ,
	home_passed_balls ,
	home_double_plays ,
	home_triple_plays ,
	home_plate_umpire_id    ,
	home_plate_umpire_name  ,
	first_base_umpire_id    ,
	first_base_umpire_name  ,
	second_base_umpire_id   ,
	second_base_umpire_name ,
	third_base_umpire_id    ,
	third_base_umpire_name  ,
	left_field_umpire_id    ,
	left_field_umpire_name  ,
	right_field_umpire_id   ,
	right_field_umpire_name ,
	visiting_manager_id   ,
	visiting_manager_name ,
	home_manager_id       ,
	home_manager_name     ,
	winning_pitcher_id   ,
	winning_pitcher_name ,
	losing_pitcher_id    ,
	losing_pitcher_name  ,
	saving_pitcher_id    ,
	saving_pitcher_name  ,
	game_winning_rbi_batter_id   ,
	game_winning_rbi_batter_name ,
	visiting_starting_pitcher_id   ,
	visiting_starting_pitcher_name ,
	home_starting_pitcher_id       ,
	home_starting_pitcher_name     ,
	visiting_player1_id       ,
	visiting_player1_name     ,
	visiting_player1_position ,
	visiting_player2_id       ,
	visiting_player2_name     ,
	visiting_player2_position ,
	visiting_player3_id       ,
	visiting_player3_name     ,
	visiting_player3_position ,
	visiting_player4_id       ,
	visiting_player4_name     ,
	visiting_player4_position ,
	visiting_player5_id       ,
	visiting_player5_name     ,
	visiting_player5_position ,
	visiting_player6_id       ,
	visiting_player6_name     ,
	visiting_player6_position ,
	visiting_player7_id       ,
	visiting_player7_name     ,
	visiting_player7_position ,
	visiting_player8_id       ,
	visiting_player8_name     ,
	visiting_player8_position ,
	visiting_player9_id       ,
	visiting_player9_name     ,
	visiting_player9_position ,
	home_player1_id       ,
	home_player1_name     ,
	home_player1_position ,
	home_player2_id       ,
	home_player2_name     ,
	home_player2_position ,
	home_player3_id       ,
	home_player3_name     ,
	home_player3_position ,
	home_player4_id       ,
	home_player4_name     ,
	home_player4_position ,
	home_player5_id       ,
	home_player5_name     ,
	home_player5_position ,
	home_player6_id       ,
	home_player6_name     ,
	home_player6_position ,
	home_player7_id       ,
	home_player7_name     ,
	home_player7_position ,
	home_player8_id       ,
	home_player8_name     ,
	home_player8_position ,
	home_player9_id       ,
	home_player9_name     ,
	home_player9_position ,
	additional_information  ,
	acquisition_information) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104, $105, $106, $107, $108, $109, $110, $111, $112, $113, $114, $115, $116, $117, $118, $119, $120, $121, $122, $123, $124, $125, $126, $127, $128, $129, $130, $131, $132, $133, $134, $135, $136, $137, $138, $139, $140, $141, $142, $143, $144, $145, $146, $147, $148, $149, $150, $151, $152, $153, $154, $155, $156, $157, $158, $159, $160, $161
	)`

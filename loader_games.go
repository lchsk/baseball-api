package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Game struct {
	Date         time.Time
	DateRaw      string
	NumberOfGame string
	DayOfWeek    string

	VisitingTeam       string
	VisitingTeamLeague string
	VisitingGameNumber int

	HomeTeam           string
	HomeTeamLeague     string
	HomeTeamGameNumber int

	VisitingTeamScore int
	HomeTeamScore     int

	GameLengthInOuts  int
	DayNightIndicator string

	CompletionInformation string
	ForfeitInformation    string
	ProtestInformation    string

	ParkID            string
	Attendance        int
	TimeOfGameInMins  int
	VisitingLineScore string
	HomeLineScore     string

	// Visiting team offensive statistics
	VisitingAB   int
	VisitingH    int
	Visiting2B   int
	Visiting3B   int
	VisitingHR   int
	VisitingRBI  int
	VisitingSH   int
	VisitingSF   int
	VisitingHBP  int
	VisitingBB   int
	VisitingIBB  int
	VisitingK    int
	VisitingSB   int
	VisitingCS   int
	VisitingGIDP int
	VisitingCI   int
	VisitingLOB  int

	// Visiting team pitching statistics
	VisitingPitchersUsed         int
	VisitingIndividualEarnedRuns int
	VisitingTeamEarnedRuns       int
	VisitingWildPitches          int
	VisitingBalks                int

	// Visiting team defensive statistics
	VisitingPutouts     int
	VisitingAssists     int
	VisitingErrors      int
	VisitingPassedBalls int
	VisitingDoublePlays int
	VisitingTriplePlays int

	// Home team offensive statistics
	HomeAB   int
	HomeH    int
	Home2B   int
	Home3B   int
	HomeHR   int
	HomeRBI  int
	HomeSH   int
	HomeSF   int
	HomeHBP  int
	HomeBB   int
	HomeIBB  int
	HomeK    int
	HomeSB   int
	HomeCS   int
	HomeGIDP int
	HomeCI   int
	HomeLOB  int

	// Home team pitching statistics
	HomePitchersUsed         int
	HomeIndividualEarnedRuns int
	HomeTeamEarnedRuns       int
	HomeWildPitches          int
	HomeBalks                int

	// Home team defensive statistics
	HomePutouts     int
	HomeAssists     int
	HomeErrors      int
	HomePassedBalls int
	HomeDoublePlays int
	HomeTriplePlays int

	// Umpires
	HomePlateUmpireID    string
	HomePlateUmpireName  string
	FirstBaseUmpireID    string
	FirstBaseUmpireName  string
	SecondBaseUmpireID   string
	SecondBaseUmpireName string
	ThirdBaseUmpireID    string
	ThirdBaseUmpireName  string
	LeftFieldUmpireID    string
	LeftFieldUmpireName  string
	RightFieldUmpireID   string
	RightFieldUmpireName string

	// Managers
	VisitingManagerID   string
	VisitingManagerName string
	HomeManagerID       string
	HomeManagerName     string

	WinningPitcherID   string
	WinningPitcherName string
	LosingPitcherID    string
	LosingPitcherName  string
	SavingPitcherID    string
	SavingPitcherName  string

	GameWinningRBIBatterID   string
	GameWinningRBIBatterName string

	VisitingStartingPitcherID   string
	VisitingStartingPitcherName string
	HomeStartingPitcherID       string
	HomeStartingPitcherName     string

	// Visiting Lineups
	VisitingPlayer1ID       string
	VisitingPlayer1Name     string
	VisitingPlayer1Position int
	VisitingPlayer2ID       string
	VisitingPlayer2Name     string
	VisitingPlayer2Position int
	VisitingPlayer3ID       string
	VisitingPlayer3Name     string
	VisitingPlayer3Position int
	VisitingPlayer4ID       string
	VisitingPlayer4Name     string
	VisitingPlayer4Position int
	VisitingPlayer5ID       string
	VisitingPlayer5Name     string
	VisitingPlayer5Position int
	VisitingPlayer6ID       string
	VisitingPlayer6Name     string
	VisitingPlayer6Position int
	VisitingPlayer7ID       string
	VisitingPlayer7Name     string
	VisitingPlayer7Position int
	VisitingPlayer8ID       string
	VisitingPlayer8Name     string
	VisitingPlayer8Position int
	VisitingPlayer9ID       string
	VisitingPlayer9Name     string
	VisitingPlayer9Position int

	HomePlayer1ID       string
	HomePlayer1Name     string
	HomePlayer1Position int
	HomePlayer2ID       string
	HomePlayer2Name     string
	HomePlayer2Position int
	HomePlayer3ID       string
	HomePlayer3Name     string
	HomePlayer3Position int
	HomePlayer4ID       string
	HomePlayer4Name     string
	HomePlayer4Position int
	HomePlayer5ID       string
	HomePlayer5Name     string
	HomePlayer5Position int
	HomePlayer6ID       string
	HomePlayer6Name     string
	HomePlayer6Position int
	HomePlayer7ID       string
	HomePlayer7Name     string
	HomePlayer7Position int
	HomePlayer8ID       string
	HomePlayer8Name     string
	HomePlayer8Position int
	HomePlayer9ID       string
	HomePlayer9Name     string
	HomePlayer9Position int

	AdditionalInformation  string
	AcquisitionInformation string
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)

	if err != nil {
		return -1
	}

	return i
}

func handleNull(value string) string {
	if value == "(none)" {
		return ""
	}

	return value
}

func readLine(line []string) *Game {
	visitingGameNumber := parseInt(line[5])
	homeGameNumber := parseInt(line[8])

	visitingScore := parseInt(line[9])
	homeScore := parseInt(line[10])

	gameLengthInOuts := parseInt(line[11])

	attendance := parseInt(line[17])
	timeOfGameInMins := parseInt(line[18])

	layout := "20060102"
	var parsedDate time.Time
	parsedDate, err := time.Parse(layout, line[0])

	if err != nil {
		parsedDate = time.Unix(0, 0)
	}

	return &Game{
		Date:               parsedDate,
		DateRaw:            line[0],
		NumberOfGame:       line[1],
		DayOfWeek:          line[2],
		VisitingTeam:       line[3],
		VisitingTeamLeague: line[4],
		VisitingGameNumber: visitingGameNumber,
		HomeTeam:           line[6],
		HomeTeamLeague:     line[7],
		HomeTeamGameNumber: homeGameNumber,
		VisitingTeamScore:  visitingScore,
		HomeTeamScore:      homeScore,

		GameLengthInOuts:      gameLengthInOuts,
		DayNightIndicator:     line[12],
		CompletionInformation: line[13],
		ForfeitInformation:    line[14],
		ProtestInformation:    line[15],
		ParkID:                line[16],
		Attendance:            attendance,
		TimeOfGameInMins:      timeOfGameInMins,
		VisitingLineScore:     line[19],
		HomeLineScore:         line[20],

		VisitingAB:   parseInt(line[21]),
		VisitingH:    parseInt(line[22]),
		Visiting2B:   parseInt(line[23]),
		Visiting3B:   parseInt(line[24]),
		VisitingHR:   parseInt(line[25]),
		VisitingRBI:  parseInt(line[26]),
		VisitingSH:   parseInt(line[27]),
		VisitingSF:   parseInt(line[28]),
		VisitingHBP:  parseInt(line[29]),
		VisitingBB:   parseInt(line[30]),
		VisitingIBB:  parseInt(line[31]),
		VisitingK:    parseInt(line[32]),
		VisitingSB:   parseInt(line[33]),
		VisitingCS:   parseInt(line[34]),
		VisitingGIDP: parseInt(line[35]),
		VisitingCI:   parseInt(line[36]),
		VisitingLOB:  parseInt(line[37]),

		VisitingPitchersUsed:         parseInt(line[38]),
		VisitingIndividualEarnedRuns: parseInt(line[39]),
		VisitingTeamEarnedRuns:       parseInt(line[40]),
		VisitingWildPitches:          parseInt(line[41]),
		VisitingBalks:                parseInt(line[42]),

		VisitingPutouts:     parseInt(line[43]),
		VisitingAssists:     parseInt(line[44]),
		VisitingErrors:      parseInt(line[45]),
		VisitingPassedBalls: parseInt(line[46]),
		VisitingDoublePlays: parseInt(line[47]),
		VisitingTriplePlays: parseInt(line[48]),

		HomeAB:   parseInt(line[49]),
		HomeH:    parseInt(line[50]),
		Home2B:   parseInt(line[51]),
		Home3B:   parseInt(line[52]),
		HomeHR:   parseInt(line[53]),
		HomeRBI:  parseInt(line[54]),
		HomeSH:   parseInt(line[55]),
		HomeSF:   parseInt(line[56]),
		HomeHBP:  parseInt(line[57]),
		HomeBB:   parseInt(line[58]),
		HomeIBB:  parseInt(line[59]),
		HomeK:    parseInt(line[60]),
		HomeSB:   parseInt(line[61]),
		HomeCS:   parseInt(line[62]),
		HomeGIDP: parseInt(line[63]),
		HomeCI:   parseInt(line[64]),
		HomeLOB:  parseInt(line[65]),

		HomePitchersUsed:         parseInt(line[66]),
		HomeIndividualEarnedRuns: parseInt(line[67]),
		HomeTeamEarnedRuns:       parseInt(line[68]),
		HomeWildPitches:          parseInt(line[69]),
		HomeBalks:                parseInt(line[70]),

		HomePutouts:     parseInt(line[71]),
		HomeAssists:     parseInt(line[72]),
		HomeErrors:      parseInt(line[73]),
		HomePassedBalls: parseInt(line[74]),
		HomeDoublePlays: parseInt(line[75]),
		HomeTriplePlays: parseInt(line[76]),

		HomePlateUmpireID:    line[77],
		HomePlateUmpireName:  line[78],
		FirstBaseUmpireID:    line[79],
		FirstBaseUmpireName:  line[80],
		SecondBaseUmpireID:   line[81],
		SecondBaseUmpireName: line[82],
		ThirdBaseUmpireID:    line[83],
		ThirdBaseUmpireName:  line[84],
		LeftFieldUmpireID:    line[85],
		LeftFieldUmpireName:  handleNull(line[86]),
		RightFieldUmpireID:   line[87],
		RightFieldUmpireName: handleNull(line[88]),

		VisitingManagerID:   line[89],
		VisitingManagerName: line[90],
		HomeManagerID:       line[91],
		HomeManagerName:     line[92],

		WinningPitcherID:   line[93],
		WinningPitcherName: line[94],
		LosingPitcherID:    line[95],
		LosingPitcherName:  line[96],
		SavingPitcherID:    handleNull(line[97]),
		SavingPitcherName:  handleNull(line[98]),

		GameWinningRBIBatterID:   handleNull(line[99]),
		GameWinningRBIBatterName: handleNull(line[100]),

		VisitingStartingPitcherID:   line[101],
		VisitingStartingPitcherName: line[102],
		HomeStartingPitcherID:       line[103],
		HomeStartingPitcherName:     line[104],

		VisitingPlayer1ID:       line[105],
		VisitingPlayer1Name:     line[106],
		VisitingPlayer1Position: parseInt(line[107]),
		VisitingPlayer2ID:       line[108],
		VisitingPlayer2Name:     line[109],
		VisitingPlayer2Position: parseInt(line[110]),
		VisitingPlayer3ID:       line[111],
		VisitingPlayer3Name:     line[112],
		VisitingPlayer3Position: parseInt(line[113]),
		VisitingPlayer4ID:       line[114],
		VisitingPlayer4Name:     line[115],
		VisitingPlayer4Position: parseInt(line[116]),
		VisitingPlayer5ID:       line[117],
		VisitingPlayer5Name:     line[118],
		VisitingPlayer5Position: parseInt(line[119]),
		VisitingPlayer6ID:       line[120],
		VisitingPlayer6Name:     line[121],
		VisitingPlayer6Position: parseInt(line[122]),
		VisitingPlayer7ID:       line[123],
		VisitingPlayer7Name:     line[124],
		VisitingPlayer7Position: parseInt(line[125]),
		VisitingPlayer8ID:       line[126],
		VisitingPlayer8Name:     line[127],
		VisitingPlayer8Position: parseInt(line[128]),
		VisitingPlayer9ID:       line[129],
		VisitingPlayer9Name:     line[130],
		VisitingPlayer9Position: parseInt(line[131]),

		HomePlayer1ID:       line[132],
		HomePlayer1Name:     line[133],
		HomePlayer1Position: parseInt(line[134]),
		HomePlayer2ID:       line[135],
		HomePlayer2Name:     line[136],
		HomePlayer2Position: parseInt(line[137]),
		HomePlayer3ID:       line[138],
		HomePlayer3Name:     line[139],
		HomePlayer3Position: parseInt(line[140]),
		HomePlayer4ID:       line[141],
		HomePlayer4Name:     line[142],
		HomePlayer4Position: parseInt(line[143]),
		HomePlayer5ID:       line[144],
		HomePlayer5Name:     line[145],
		HomePlayer5Position: parseInt(line[146]),
		HomePlayer6ID:       line[147],
		HomePlayer6Name:     line[148],
		HomePlayer6Position: parseInt(line[149]),
		HomePlayer7ID:       line[150],
		HomePlayer7Name:     line[151],
		HomePlayer7Position: parseInt(line[152]),
		HomePlayer8ID:       line[153],
		HomePlayer8Name:     line[154],
		HomePlayer8Position: parseInt(line[155]),
		HomePlayer9ID:       line[156],
		HomePlayer9Name:     line[157],
		HomePlayer9Position: parseInt(line[158]),

		AdditionalInformation:  line[159],
		AcquisitionInformation: line[160],
	}
}

func getGameLogsFiles(dir string) ([]string, error) {
	var gameLogFiles []string
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return gameLogFiles, err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		ext := strings.ToLower(filepath.Ext(path))

		if ext == ".txt" {
			gameLogFiles = append(gameLogFiles, path)
		}
	}

	return gameLogFiles, nil
}

func parseGames(gameLogFiles []string) ([]*Game, error) {
	var games []*Game

	for _, path := range gameLogFiles {
		csvFile, err := os.Open(path)

		if err != nil {
			return games, err
		}

		reader := csv.NewReader(bufio.NewReader(csvFile))

		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			games = append(games, readLine(line))
		}
	}

	return games, nil
}

func loadGameLogs(dir string) {
	gameLogFiles, err := getGameLogsFiles(dir)

	if err != nil {
		panic(err)
	}

	games, err := parseGames(gameLogFiles)

	db := getDBConnection()

	log.Println("Inserting games")

	for _, game := range games {
		err := insertGame(game, db)

		if err != nil {
			log.Fatalf("Error when inserting game: %v %s", game, err)
		}
	}

	db.Close()

	log.Printf("Processed %d games", len(games))
}

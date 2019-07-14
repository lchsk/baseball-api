package loader

import (
	"encoding/csv"
	"log"
	"strings"
	"testing"
	"time"
)

const data = `"20180329","0","Thu","BOS","AL",1,"TBA","AL",1,4,6,51,"D","","","","STP01",31042,180,"030000100","00000006x",33,8,4,0,1,4,0,0,0,2,0,6,0,0,0,0,4,4,6,6,0,0,24,6,0,1,0,0,28,4,1,1,0,6,0,0,0,7,0,11,0,0,0,0,5,3,4,4,0,0,27,12,0,0,1,0,"nelsj901","Jeff Nelson","diazl901","Laz Diaz","fleta901","Andy Fletcher","gonzm901","Manny Gonzalez","","(none)","","(none)","coraa001","Alex Cora","cashk001","Kevin Cash","pruia001","Austin Pruitt","smitc004","Carson Smith","coloa001","Alex Colome","spand001","Denard Span","salec001","Chris Sale","archc001","Chris Archer","bettm001","Mookie Betts",9,"benia002","Andrew Benintendi",7,"ramih003","Hanley Ramirez",3,"martj006","J.D. Martinez",10,"bogax001","Xander Bogaerts",6,"dever001","Rafael Devers",5,"nunee002","Eduardo Nunez",4,"bradj001","Jackie Bradley",8,"vazqc001","Christian Vazquez",2,"duffm002","Matt Duffy",5,"kierk001","Kevin Kiermaier",8,"gomec002","Carlos Gomez",9,"cronc002","C.J. Cron",3,"ramow001","Wilson Ramos",2,"spand001","Denard Span",7,"hecha001","Adeiny Hechavarria",6,"robed004","Daniel Robertson",4,"refsr001","Rob Refsnyder",10,"","Y"`

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestReadLine(t *testing.T) {
	r := csv.NewReader(strings.NewReader(data))

	record, err := r.Read()

	if err != nil {
		log.Fatal(err)
	}

	game := readLine(record)

	assertEqual(t, game.Date, time.Date(2018, time.March, 29, 0, 0, 0, 0, time.UTC))
	assertEqual(t, game.DateRaw, "20180329")
	assertEqual(t, game.NumberOfGame, "0")
	assertEqual(t, game.DayOfWeek, "Thu")

	assertEqual(t, game.VisitingTeam, "BOS")
	assertEqual(t, game.VisitingTeamLeague, "AL")
	assertEqual(t, game.VisitingGameNumber, 1)

	assertEqual(t, game.HomeTeam, "TBA")
	assertEqual(t, game.HomeTeamLeague, "AL")
	assertEqual(t, game.HomeTeamGameNumber, 1)

	assertEqual(t, game.VisitingTeamScore, 4)
	assertEqual(t, game.HomeTeamScore, 6)

	assertEqual(t, game.GameLengthInOuts, 51)

	assertEqual(t, game.DayNightIndicator, "D")

	assertEqual(t, game.CompletionInformation, "")
	assertEqual(t, game.ForfeitInformation, "")
	assertEqual(t, game.ProtestInformation, "")

	assertEqual(t, game.ParkID, "STP01")
	assertEqual(t, game.Attendance, 31042)
	assertEqual(t, game.TimeOfGameInMins, 180)

	assertEqual(t, game.VisitingLineScore, "030000100")
	assertEqual(t, game.HomeLineScore, "00000006x")

	// Visiting Team Offensive Stats
	assertEqual(t, game.VisitingAB, 33)
	assertEqual(t, game.VisitingH, 8)
	assertEqual(t, game.Visiting2B, 4)
	assertEqual(t, game.Visiting3B, 0)
	assertEqual(t, game.VisitingHR, 1)
	assertEqual(t, game.VisitingRBI, 4)
	assertEqual(t, game.VisitingSH, 0)
	assertEqual(t, game.VisitingSF, 0)
	assertEqual(t, game.VisitingHBP, 0)
	assertEqual(t, game.VisitingBB, 2)
	assertEqual(t, game.VisitingIBB, 0)
	assertEqual(t, game.VisitingK, 6)
	assertEqual(t, game.VisitingSB, 0)
	assertEqual(t, game.VisitingCS, 0)
	assertEqual(t, game.VisitingGIDP, 0)
	assertEqual(t, game.VisitingCI, 0)
	assertEqual(t, game.VisitingLOB, 4)

	// Visiting Team Pitching Stats
	assertEqual(t, game.VisitingPitchersUsed, 4)
	assertEqual(t, game.VisitingIndividualEarnedRuns, 6)
	assertEqual(t, game.VisitingTeamEarnedRuns, 6)
	assertEqual(t, game.VisitingWildPitches, 0)
	assertEqual(t, game.VisitingBalks, 0)

	// Visiting Team Defensive Stats
	assertEqual(t, game.VisitingPutouts, 24)
	assertEqual(t, game.VisitingAssists, 6)
	assertEqual(t, game.VisitingErrors, 0)
	assertEqual(t, game.VisitingPassedBalls, 1)
	assertEqual(t, game.VisitingDoublePlays, 0)
	assertEqual(t, game.VisitingTriplePlays, 0)

	// Home Team Offensive Stats
	assertEqual(t, game.HomeAB, 28)
	assertEqual(t, game.HomeH, 4)
	assertEqual(t, game.Home2B, 1)
	assertEqual(t, game.Home3B, 1)
	assertEqual(t, game.HomeHR, 0)
	assertEqual(t, game.HomeRBI, 6)
	assertEqual(t, game.HomeSH, 0)
	assertEqual(t, game.HomeSF, 0)
	assertEqual(t, game.HomeHBP, 0)
	assertEqual(t, game.HomeBB, 7)
	assertEqual(t, game.HomeIBB, 0)
	assertEqual(t, game.HomeK, 11)
	assertEqual(t, game.HomeSB, 0)
	assertEqual(t, game.HomeCS, 0)
	assertEqual(t, game.HomeGIDP, 0)
	assertEqual(t, game.HomeCI, 0)
	assertEqual(t, game.HomeLOB, 5)

	// Home Team Pitching Stats
	assertEqual(t, game.HomePitchersUsed, 3)
	assertEqual(t, game.HomeIndividualEarnedRuns, 4)
	assertEqual(t, game.HomeTeamEarnedRuns, 4)
	assertEqual(t, game.HomeWildPitches, 0)
	assertEqual(t, game.HomeBalks, 0)

	// Home Team Defensive Stats
	assertEqual(t, game.HomePutouts, 27)
	assertEqual(t, game.HomeAssists, 12)
	assertEqual(t, game.HomeErrors, 0)
	assertEqual(t, game.HomePassedBalls, 0)
	assertEqual(t, game.HomeDoublePlays, 1)
	assertEqual(t, game.HomeTriplePlays, 0)

	assertEqual(t, game.HomePlateUmpireID, "nelsj901")
	assertEqual(t, game.HomePlateUmpireName, "Jeff Nelson")
	assertEqual(t, game.FirstBaseUmpireID, "diazl901")
	assertEqual(t, game.FirstBaseUmpireName, "Laz Diaz")
	assertEqual(t, game.SecondBaseUmpireID, "fleta901")
	assertEqual(t, game.SecondBaseUmpireName, "Andy Fletcher")
	assertEqual(t, game.ThirdBaseUmpireID, "gonzm901")
	assertEqual(t, game.ThirdBaseUmpireName, "Manny Gonzalez")
	assertEqual(t, game.LeftFieldUmpireID, "")
	assertEqual(t, game.LeftFieldUmpireName, "")
	assertEqual(t, game.RightFieldUmpireID, "")
	assertEqual(t, game.RightFieldUmpireName, "")

	assertEqual(t, game.VisitingManagerID, "coraa001")
	assertEqual(t, game.VisitingManagerName, "Alex Cora")
	assertEqual(t, game.HomeManagerID, "cashk001")
	assertEqual(t, game.HomeManagerName, "Kevin Cash")

	assertEqual(t, game.WinningPitcherID, "pruia001")
	assertEqual(t, game.WinningPitcherName, "Austin Pruitt")
	assertEqual(t, game.LosingPitcherID, "smitc004")
	assertEqual(t, game.LosingPitcherName, "Carson Smith")
	assertEqual(t, game.SavingPitcherID, "coloa001")
	assertEqual(t, game.SavingPitcherName, "Alex Colome")

	assertEqual(t, game.GameWinningRBIBatterID, "spand001")
	assertEqual(t, game.GameWinningRBIBatterName, "Denard Span")

	assertEqual(t, game.VisitingStartingPitcherID, "salec001")
	assertEqual(t, game.VisitingStartingPitcherName, "Chris Sale")
	assertEqual(t, game.HomeStartingPitcherID, "archc001")
	assertEqual(t, game.HomeStartingPitcherName, "Chris Archer")

	assertEqual(t, game.VisitingPlayer1ID, "bettm001")
	assertEqual(t, game.VisitingPlayer1Name, "Mookie Betts")
	assertEqual(t, game.VisitingPlayer1Position, 9)
	assertEqual(t, game.VisitingPlayer2ID, "benia002")
	assertEqual(t, game.VisitingPlayer2Name, "Andrew Benintendi")
	assertEqual(t, game.VisitingPlayer2Position, 7)
	assertEqual(t, game.VisitingPlayer3ID, "ramih003")
	assertEqual(t, game.VisitingPlayer3Name, "Hanley Ramirez")
	assertEqual(t, game.VisitingPlayer3Position, 3)
	assertEqual(t, game.VisitingPlayer4ID, "martj006")
	assertEqual(t, game.VisitingPlayer4Name, "J.D. Martinez")
	assertEqual(t, game.VisitingPlayer4Position, 10)
	assertEqual(t, game.VisitingPlayer5ID, "bogax001")
	assertEqual(t, game.VisitingPlayer5Name, "Xander Bogaerts")
	assertEqual(t, game.VisitingPlayer5Position, 6)
	assertEqual(t, game.VisitingPlayer6ID, "dever001")
	assertEqual(t, game.VisitingPlayer6Name, "Rafael Devers")
	assertEqual(t, game.VisitingPlayer6Position, 5)
	assertEqual(t, game.VisitingPlayer7ID, "nunee002")
	assertEqual(t, game.VisitingPlayer7Name, "Eduardo Nunez")
	assertEqual(t, game.VisitingPlayer7Position, 4)
	assertEqual(t, game.VisitingPlayer8ID, "bradj001")
	assertEqual(t, game.VisitingPlayer8Name, "Jackie Bradley")
	assertEqual(t, game.VisitingPlayer8Position, 8)
	assertEqual(t, game.VisitingPlayer9ID, "vazqc001")
	assertEqual(t, game.VisitingPlayer9Name, "Christian Vazquez")
	assertEqual(t, game.VisitingPlayer9Position, 2)

	assertEqual(t, game.HomePlayer1ID, "duffm002")
	assertEqual(t, game.HomePlayer1Name, "Matt Duffy")
	assertEqual(t, game.HomePlayer1Position, 5)
	assertEqual(t, game.HomePlayer2ID, "kierk001")
	assertEqual(t, game.HomePlayer2Name, "Kevin Kiermaier")
	assertEqual(t, game.HomePlayer2Position, 8)
	assertEqual(t, game.HomePlayer3ID, "gomec002")
	assertEqual(t, game.HomePlayer3Name, "Carlos Gomez")
	assertEqual(t, game.HomePlayer3Position, 9)
	assertEqual(t, game.HomePlayer4ID, "cronc002")
	assertEqual(t, game.HomePlayer4Name, "C.J. Cron")
	assertEqual(t, game.HomePlayer4Position, 3)
	assertEqual(t, game.HomePlayer5ID, "ramow001")
	assertEqual(t, game.HomePlayer5Name, "Wilson Ramos")
	assertEqual(t, game.HomePlayer5Position, 2)
	assertEqual(t, game.HomePlayer6ID, "spand001")
	assertEqual(t, game.HomePlayer6Name, "Denard Span")
	assertEqual(t, game.HomePlayer6Position, 7)
	assertEqual(t, game.HomePlayer7ID, "hecha001")
	assertEqual(t, game.HomePlayer7Name, "Adeiny Hechavarria")
	assertEqual(t, game.HomePlayer7Position, 6)
	assertEqual(t, game.HomePlayer8ID, "robed004")
	assertEqual(t, game.HomePlayer8Name, "Daniel Robertson")
	assertEqual(t, game.HomePlayer8Position, 4)
	assertEqual(t, game.HomePlayer9ID, "refsr001")
	assertEqual(t, game.HomePlayer9Name, "Rob Refsnyder")
	assertEqual(t, game.HomePlayer9Position, 10)

	assertEqual(t, game.AdditionalInformation, "")
	assertEqual(t, game.AcquisitionInformation, "Y")

}

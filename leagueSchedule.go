package main

import "fmt"

type teamStruct struct {
	name string
}

type fixtureStruct struct {
	homeTeam *teamStruct
	awayTeam *teamStruct
}

func (f fixtureStruct) printable() string {
	return "Home team: " + f.homeTeam.name + ", Away team: " + f.awayTeam.name
}

type matchdayPileStruct struct {
	next         *matchdayPileStruct
	prev         *matchdayPileStruct
	fixtures     map[int]bool
	fixedFixture *fixtureStruct
}

type provisionalMachdayFixturesStruct struct {
	next   *provisionalMachdayFixturesStruct
	prev   *provisionalMachdayFixturesStruct
	offset int
}

var team1 = teamStruct{"team1"}
var team2 = teamStruct{"team2"}
var team3 = teamStruct{"team3"}
var team4 = teamStruct{"team4"}
var team5 = teamStruct{"team5"}
var team6 = teamStruct{"team6"}
var team7 = teamStruct{"team7"}
var team8 = teamStruct{"team8"}
var team9 = teamStruct{"team9"}
var team10 = teamStruct{"team10"}
var team11 = teamStruct{"team11"}
var team12 = teamStruct{"team12"}
var team13 = teamStruct{"team13"}
var team14 = teamStruct{"team14"}
var team15 = teamStruct{"team15"}
var team16 = teamStruct{"team16"}
var team17 = teamStruct{"team17"}
var team18 = teamStruct{"team18"}
var team19 = teamStruct{"team19"}
var team20 = teamStruct{"team20"}

var allTeams = []teamStruct{team1, team2, team3, team4, team5, team6, team7, team8, team9, team10, team11, team12, team13, team14, team15, team16, team17, team18, team19, team20}
var numMatchdays = len(allTeams) - 1

func main() {
	allFixtures := createAllFixtures(allTeams)
	matchdayPileHead := createAllMatchdays(allFixtures)
	printAllMatchdays(matchdayPileHead, allFixtures)
}

func printAllMatchdays(headPile *matchdayPileStruct, allFixtures []fixtureStruct) {
	matchdayNum := 1
	for currPile := headPile; currPile != nil; currPile = currPile.next {
		printMatchdaySchedule(matchdayNum, currPile, allFixtures)
		matchdayNum = matchdayNum + 1
	}
}

func printMatchdaySchedule(matchdayNum int, pile *matchdayPileStruct, allFixtures []fixtureStruct) {
	fmt.Printf("-----Fixtures finalized for matchday %d -----\n", matchdayNum)
	for offset := range pile.fixtures {
		fmt.Println(allFixtures[offset].printable())
	}
}

func createAllMatchdays(allFixtures []fixtureStruct) *matchdayPileStruct {
	var matchdayPileHead *matchdayPileStruct
	currPile := matchdayPileHead
	for i := 1; i <= numMatchdays; i++ {
		newCurrPile := createMatchdaySchedule(i, allFixtures, currPile)
		if matchdayPileHead == nil {
			matchdayPileHead = newCurrPile
		}
		currPile = newCurrPile
	}
	return matchdayPileHead
}

func createMatchdaySchedule(matchdayNum int, allFixtures []fixtureStruct, currPile *matchdayPileStruct) *matchdayPileStruct {
	haveAllTeamsBeenScheduledToPlay := false
	newPile := createMatchdayPile()
	provisionalMachdayFixturesHead := createProvisionalMatchdayFixturesStruct()

	for !haveAllTeamsBeenScheduledToPlay {
		var homeTeamUnderConsideration *teamStruct
		teamsScheduledToPlayOnMatchday := make(map[*teamStruct]bool)

		for i := 0; i < len(allFixtures); i++ {
			if teamsScheduledToPlayOnMatchday[allFixtures[i].homeTeam] {
				continue
			}
			if homeTeamUnderConsideration == nil {
				homeTeamUnderConsideration = allFixtures[i].homeTeam
			}
			if canFixtureBeScheduledForThisMatchday(allFixtures, i, currPile, teamsScheduledToPlayOnMatchday) {
				if newPile.fixedFixture == nil {
					newPile.fixedFixture = &allFixtures[i]
				}
				handleMatchdayFixture(&newPile, &provisionalMachdayFixturesHead, allFixtures, i, teamsScheduledToPlayOnMatchday)
				homeTeamUnderConsideration = nil
			} else if hasHomeTeamFixtureNotBeenFoundForMatchday(allFixtures, i, homeTeamUnderConsideration, teamsScheduledToPlayOnMatchday) {
				offsetFromWhereToSearchForNextFixtureForHomeTeam := removePreviousFixtureScheduledForMatchday(&newPile, &provisionalMachdayFixturesHead)
				delete(teamsScheduledToPlayOnMatchday, allFixtures[offsetFromWhereToSearchForNextFixtureForHomeTeam].homeTeam)
				delete(teamsScheduledToPlayOnMatchday, allFixtures[offsetFromWhereToSearchForNextFixtureForHomeTeam].awayTeam)
				homeTeamUnderConsideration = allFixtures[offsetFromWhereToSearchForNextFixtureForHomeTeam].homeTeam
				i = offsetFromWhereToSearchForNextFixtureForHomeTeam
				continue
			}
		}
		haveAllTeamsBeenScheduledToPlay = haveAllTeamsBeenScheduledToPlayOnMatchday(teamsScheduledToPlayOnMatchday)
	}

	if currPile != nil {
		currPile.next = &newPile
		newPile.prev = currPile
	}
	currPile = &newPile
	return currPile
}

func haveAllTeamsBeenScheduledToPlayOnMatchday(teamsScheduledToPlayOnMatchday map[*teamStruct]bool) bool {
	if len(teamsScheduledToPlayOnMatchday) == len(allTeams) {
		return true
	}
	return false
}

func removePreviousFixtureScheduledForMatchday(pile *matchdayPileStruct, provHead *provisionalMachdayFixturesStruct) int {
	toBeRemoved := provHead.prev
	newPrev := toBeRemoved.prev

	newPrev.next = provHead
	provHead.prev = newPrev

	toBeRemoved.prev = nil
	toBeRemoved.next = nil

	delete(pile.fixtures, toBeRemoved.offset)
	return toBeRemoved.offset
}

func hasHomeTeamFixtureNotBeenFoundForMatchday(allFixtures []fixtureStruct, offset int, homeTeamUnderConsideration *teamStruct, teamsScheduledToPlayOnMatchday map[*teamStruct]bool) bool {
	if homeTeamUnderConsideration == nil {
		return false
	}
	if teamsScheduledToPlayOnMatchday[homeTeamUnderConsideration] {
		return false
	}
	if allFixtures[offset].homeTeam == homeTeamUnderConsideration {
		if offset == len(allFixtures)-1 {
			return true
		}
		if allFixtures[offset+1].homeTeam != homeTeamUnderConsideration {
			return true
		}
		return false
	}
	return true
}

func canFixtureBeScheduledForThisMatchday(allFixtures []fixtureStruct, offset int, currPile *matchdayPileStruct, teamsScheduledToPlayOnMatchday map[*teamStruct]bool) bool {
	if hasEitherTeamAlreadyPlayedOnMatchday(allFixtures[offset], teamsScheduledToPlayOnMatchday) {
		return false
	}
	if hasFixtureBeenScheduled(currPile, offset) {
		return false
	}
	return true
}

func hasEitherTeamAlreadyPlayedOnMatchday(fixture fixtureStruct, teamsScheduledToPlayOnMatchday map[*teamStruct]bool) bool {
	if teamsScheduledToPlayOnMatchday[fixture.homeTeam] || teamsScheduledToPlayOnMatchday[fixture.awayTeam] {
		return true
	}
	return false
}

func hasFixtureBeenScheduled(pile *matchdayPileStruct, fixtureOffset int) bool {
	for curr := pile; curr != nil; curr = curr.prev {
		if curr.fixtures[fixtureOffset] == true {
			return true
		}
	}
	return false
}

func handleMatchdayFixture(pile *matchdayPileStruct, provisionalMachdayFixturesHead *provisionalMachdayFixturesStruct, allFixtures []fixtureStruct, offset int, teamsScheduledToPlayOnMatchday map[*teamStruct]bool) {
	teamsScheduledToPlayOnMatchday[allFixtures[offset].homeTeam] = true
	teamsScheduledToPlayOnMatchday[allFixtures[offset].awayTeam] = true
	pile.fixtures[offset] = true

	if provisionalMachdayFixturesHead.offset == -1 {
		provisionalMachdayFixturesHead.offset = offset
	} else {
		prov := createProvisionalMatchdayFixturesStruct()
		prov.offset = offset

		// circular linked list... adding element at the end
		prevLastNode := provisionalMachdayFixturesHead.prev
		prov.next = provisionalMachdayFixturesHead
		provisionalMachdayFixturesHead.prev = &prov

		if prevLastNode != nil {
			prevLastNode.next = &prov
			prov.prev = prevLastNode
		} else {
			// first node after the head node is being added
			provisionalMachdayFixturesHead.next = &prov
			prov.prev = provisionalMachdayFixturesHead
		}
	}
}

func createAllFixtures(teams []teamStruct) []fixtureStruct {
	allFixtures := make([]fixtureStruct, 0)
	for i := 0; i < len(teams)-1; i++ {
		for j := i + 1; j < len(teams); j++ {
			allFixtures = append(allFixtures, createFixture(&teams[i], &teams[j]))
		}
	}
	return allFixtures
}

func createFixture(team1, team2 *teamStruct) fixtureStruct {
	return fixtureStruct{team1, team2}
}

func createMatchdayPile() matchdayPileStruct {
	return matchdayPileStruct{nil, nil, make(map[int]bool), nil}
}

func createProvisionalMatchdayFixturesStruct() provisionalMachdayFixturesStruct {
	return provisionalMachdayFixturesStruct{nil, nil, -1}
}

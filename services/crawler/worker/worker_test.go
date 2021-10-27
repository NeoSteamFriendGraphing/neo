package worker

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/iamcathal/neo/services/crawler/datastructures"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestVerifyFormatOfSteamIDsVerifiesTwoValidSteamIDs(t *testing.T) {
	expectedSteamIDs := []string{"12345678901234456", "72348978301996243"}
	inputData := datastructures.CrawlUsersInput{
		FirstSteamID:  expectedSteamIDs[0],
		SecondSteamID: expectedSteamIDs[1],
	}

	receivedValidSteamIDs, err := VerifyFormatOfSteamIDs(inputData)

	assert.Nil(t, err)
	assert.Equal(t, expectedSteamIDs, receivedValidSteamIDs, "expect two valid format steamIDs are returned")
}

func TestVerifyFormatOfSteamIDsReturnsNothingForTwoInvalidFormatSteamIDs(t *testing.T) {
	expectedSteamIDs := []string{"12345634456", "0"}
	inputData := datastructures.CrawlUsersInput{
		FirstSteamID:  expectedSteamIDs[0],
		SecondSteamID: expectedSteamIDs[1],
	}

	receivedValidSteamIDs, err := VerifyFormatOfSteamIDs(inputData)

	assert.Nil(t, err)
	assert.Len(t, receivedValidSteamIDs, 0, "expect to receive back no steamIDs for two invalid steamID inputs")
}

func TestExctractSteamIDsfromFriendsList(t *testing.T) {
	expectedList := []string{"1234", "5436", "6718"}
	friends := datastructures.Friendslist{
		Friends: []datastructures.Friend{
			{
				Steamid: "1234",
			},
			{
				Steamid: "5436",
			},
			{
				Steamid: "6718",
			},
		},
	}

	realList := extractSteamIDsfromFriendsList(friends)

	assert.Equal(t, expectedList, realList)
}

func TestBreakSteamIDsIntoListsOf100OrLessWith100IDs(t *testing.T) {
	idList := []string{}
	for i := 0; i < 100; i++ {
		idList = append(idList, strconv.Itoa(i))
	}
	URLFormattedSteamIDs := strings.Join(idList, ",")
	expectedSteamIDList := []string{URLFormattedSteamIDs}

	realSteamIDList := breakIntoStacksOf100OrLessSteamIDs(idList)

	assert.Equal(t, expectedSteamIDList, realSteamIDList)
}

func TestBreakSteamIDsIntoListsOf100OrLessWith120IDs(t *testing.T) {
	idList := []string{}
	for i := 0; i < 120; i++ {
		idList = append(idList, strconv.Itoa(i))
	}
	firstBatchOfURLFormattedSteamIDs := strings.Join(idList[:100], ",")
	remainderBatchOfURLFormattedSteamIDs := strings.Join(idList[100:], ",")

	expectedSteamIDList := []string{firstBatchOfURLFormattedSteamIDs, remainderBatchOfURLFormattedSteamIDs}

	realSteamIDList := breakIntoStacksOf100OrLessSteamIDs(idList)

	assert.Equal(t, expectedSteamIDList, realSteamIDList)
}

func TestBreakSteamIDsIntoListsOf100OrLess20IDs(t *testing.T) {
	idList := []string{}
	for i := 0; i < 20; i++ {
		idList = append(idList, strconv.Itoa(i))
	}
	URLFormattedSteamIDs := strings.Join(idList, ",")
	expectedSteamIDList := []string{URLFormattedSteamIDs}

	realSteamIDList := breakIntoStacksOf100OrLessSteamIDs(idList)

	assert.Equal(t, expectedSteamIDList, realSteamIDList)
}

func TestBreakSteamIDsIntoListsOf100OrLessWith1911IDs(t *testing.T) {
	idList := []string{}
	for i := 0; i < 1911; i++ {
		idList = append(idList, strconv.Itoa(i))
	}

	realSteamIDList := breakIntoStacksOf100OrLessSteamIDs(idList)

	assert.Len(t, realSteamIDList, 20)
}

func TestGetUsersProfileSummaryFromSliceReturnsTheSearchedForProfile(t *testing.T) {
	expectedUserProfile := datastructures.Player{
		Steamid:  "54290543656",
		Realname: "Eddie Durcan",
	}
	exampleSummaries := []datastructures.Player{
		{
			Steamid:  "213023525435",
			Realname: "Buzz Mc Donell",
		},
		expectedUserProfile,
		{
			Steamid:  "5647568578975",
			Realname: "The Boogenhagen",
		},
	}

	found, realUserProfile := getUsersProfileSummaryFromSlice(expectedUserProfile.Steamid, exampleSummaries)

	assert.True(t, found)
	assert.Equal(t, expectedUserProfile, realUserProfile)
}

func TestGetUsersProfileSummaryFromSliceReturnsFalseWhenNotFound(t *testing.T) {
	nonExistantSteamID := "45356346547567"
	exampleSummaries := []datastructures.Player{
		{
			Steamid:  "213023525435",
			Realname: "Buzz Mc Donell",
		},
		{
			Steamid:  "5647568578975",
			Realname: "The Boogenhagen",
		},
	}

	found, realUserProfile := getUsersProfileSummaryFromSlice(nonExistantSteamID, exampleSummaries)

	assert.False(t, found)
	assert.Empty(t, realUserProfile)
}

func TestGetSteamIDsFromPlayersReturnsAllSteamIDs(t *testing.T) {
	examplePlayers := []datastructures.Player{
		{
			Steamid:  "213023525435",
			Realname: "Buzz Mc Donell",
		},
		{
			Steamid:  "54290543656",
			Realname: "Eddie Durcan",
		},
		{
			Steamid:  "5647568578975",
			Realname: "The Boogenhagen",
		},
	}
	expectedSteamIDList := []string{examplePlayers[0].Steamid, examplePlayers[1].Steamid, examplePlayers[2].Steamid}

	realSteamIDs := getSteamIDsFromPlayers(examplePlayers)

	assert.Equal(t, expectedSteamIDList, realSteamIDs)
}

func TestGetSteamIDsFromPlayersFromAnEmptySliceReturnsNothing(t *testing.T) {
	examplePlayers := []datastructures.Player{}

	realSteamIDs := getSteamIDsFromPlayers(examplePlayers)

	assert.Empty(t, realSteamIDs)
}

func TestGetPublicProfilesReturnsOnlyPublicProfiles(t *testing.T) {
	expectedPublicProfile := datastructures.Player{
		Steamid:                  "54290543656",
		Realname:                 "Eddie Durcan",
		Communityvisibilitystate: 3,
	}
	examplePlayers := []datastructures.Player{
		{
			Steamid:                  "213023525435",
			Realname:                 "Buzz Mc Donell",
			Communityvisibilitystate: 2,
		},
		expectedPublicProfile,
		{
			Steamid:                  "5647568578975",
			Realname:                 "The Boogenhagen",
			Communityvisibilitystate: 1,
		},
	}

	realPublicProfiles := getPublicProfiles(examplePlayers)

	assert.Equal(t, expectedPublicProfile, realPublicProfiles[0])
}

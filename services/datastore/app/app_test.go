package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/IamCathal/neo/services/datastore/configuration"
	"github.com/IamCathal/neo/services/datastore/controller"
	"github.com/neosteamfriendgraphing/common"
	"github.com/neosteamfriendgraphing/common/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	testSaveUserDTO dtos.SaveUserDTO
)

func TestMain(m *testing.M) {
	initTestData()

	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"/dev/null"}
	logger, err := c.Build()
	if err != nil {
		log.Fatal(err)
	}
	configuration.Logger = logger

	code := m.Run()

	os.Exit(code)
}

func TestSaveUserToDBCallsMongoDBOnce(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	mockController.On("InsertOne",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(nil, nil)
	configuration.DBClient = &mongo.Client{}

	err := SaveUserToDB(mockController, testSaveUserDTO.User)

	assert.Nil(t, err)
	mockController.AssertNumberOfCalls(t, "InsertOne", 1)
}

func TestSaveUserToDBCallsReturnsErrorWhenMongoDoes(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	expectedError := errors.New("expected error response")
	mockController.On("InsertOne",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(nil, expectedError)
	configuration.DBClient = &mongo.Client{}

	err := SaveUserToDB(mockController, testSaveUserDTO.User)

	assert.EqualError(t, err, expectedError.Error())
	mockController.AssertNumberOfCalls(t, "InsertOne", 1)
}

func TestSaveCrawlingStatsToDBForExistingUserAtMaxLevelOnlyCallsUpdate(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	maxLevelTestSaveUserDTO := testSaveUserDTO
	maxLevelTestSaveUserDTO.CurrentLevel = maxLevelTestSaveUserDTO.MaxLevel

	mockController.On("UpdateCrawlingStatus",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(true, nil)
	configuration.DBClient = &mongo.Client{}

	crawlingStatus := common.CrawlingStatus{
		OriginalCrawlTarget: maxLevelTestSaveUserDTO.User.AccDetails.SteamID,
		MaxLevel:            maxLevelTestSaveUserDTO.MaxLevel,
		CrawlID:             maxLevelTestSaveUserDTO.CrawlID,
		TotalUsersToCrawl:   len(maxLevelTestSaveUserDTO.User.FriendIDs),
	}
	err := SaveCrawlingStatsToDB(mockController, maxLevelTestSaveUserDTO.MaxLevel, crawlingStatus)

	assert.Nil(t, err)
	mockController.AssertNumberOfCalls(t, "UpdateCrawlingStatus", 1)
	mockController.AssertNotCalled(t, "InsertOne")
}

func TestSaveCrawlingStatsToDBCallsUpdateAndThenInsertForNewUser(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	configuration.DBClient = &mongo.Client{}

	// Return document does not exist when trying to update it
	mockController.On("UpdateCrawlingStatus",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(false, nil)

	// Return valid for insertion of new record
	mockController.On("InsertOne",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(nil, nil)

	crawlingStatus := common.CrawlingStatus{
		OriginalCrawlTarget: testSaveUserDTO.User.AccDetails.SteamID,
		MaxLevel:            testSaveUserDTO.MaxLevel,
		CrawlID:             testSaveUserDTO.CrawlID,
		TotalUsersToCrawl:   len(testSaveUserDTO.User.FriendIDs),
	}
	err := SaveCrawlingStatsToDB(mockController, 1, crawlingStatus)

	assert.Nil(t, err)
	mockController.AssertNumberOfCalls(t, "UpdateCrawlingStatus", 1)
	mockController.AssertNumberOfCalls(t, "InsertOne", 1)
}

func TestSaveCrawlingStatsToDBReturnsNilWhenFailsToIncrementUsersCrawledForUserOnMaxLevel(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	configuration.DBClient = &mongo.Client{}
	maxLevelTestSaveUserDTO := testSaveUserDTO
	maxLevelTestSaveUserDTO.CurrentLevel = maxLevelTestSaveUserDTO.MaxLevel

	// Return document does not exist when trying to update it
	mockController.On("UpdateCrawlingStatus",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(false, nil).Once()

	// Return an error when this max level user cannot be updated
	mockController.On("UpdateCrawlingStatus",
		mock.Anything,
		mock.Anything,
		maxLevelTestSaveUserDTO,
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).Return(false, nil).Once()

	crawlingStatus := common.CrawlingStatus{
		OriginalCrawlTarget: testSaveUserDTO.User.AccDetails.SteamID,
		MaxLevel:            testSaveUserDTO.MaxLevel,
		CrawlID:             testSaveUserDTO.CrawlID,
		TotalUsersToCrawl:   len(testSaveUserDTO.User.FriendIDs),
	}
	err := SaveCrawlingStatsToDB(mockController, testSaveUserDTO.MaxLevel, crawlingStatus)

	assert.Nil(t, err)
	mockController.AssertNumberOfCalls(t, "UpdateCrawlingStatus", 1)
	mockController.AssertNotCalled(t, "InsertOne")
}

func TestGetUser(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	mockController.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(testSaveUserDTO.User, nil)

	user, err := GetUserFromDB(mockController, testSaveUserDTO.User.AccDetails.SteamID)

	assert.NoError(t, err)
	assert.Equal(t, user, testSaveUserDTO.User)
}

func TestGetUserReturnsAnErrorAndEmptyUserWhenMongoReturnsAnError(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	expectedError := fmt.Errorf("error message")
	mockController.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(common.UserDocument{}, expectedError)

	user, err := GetUserFromDB(mockController, testSaveUserDTO.User.AccDetails.SteamID)

	assert.EqualError(t, err, expectedError.Error())
	assert.Equal(t, user, common.UserDocument{})
}

func TestGetCrawlingStatsFromDBFromCrawlID(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	configuration.DBClient = &mongo.Client{}

	crawlID := "crawlID"
	expectedCrawlingStatus := common.CrawlingStatus{
		TimeStarted: time.Now().Unix(),
		CrawlID:     crawlID,
	}
	mockController.On("GetCrawlingStatusFromDBFromCrawlID", mock.Anything, crawlID).Return(expectedCrawlingStatus, nil)

	crawlingStatus, err := GetCrawlingStatsFromDBFromCrawlID(mockController, crawlID)

	assert.Nil(t, err)
	assert.Equal(t, expectedCrawlingStatus, crawlingStatus)
	mockController.AssertNumberOfCalls(t, "GetCrawlingStatusFromDBFromCrawlID", 1)
}

func TestGetCrawlingStatsFromDBReturnsAnErrorWhenControllerMethodDoes(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	configuration.DBClient = &mongo.Client{}

	crawlID := "crawlID"
	expectedError := errors.New("expected error")
	mockController.On("GetCrawlingStatusFromDBFromCrawlID", mock.Anything, crawlID).Return(common.CrawlingStatus{}, expectedError)

	crawlingStatus, err := GetCrawlingStatsFromDBFromCrawlID(mockController, crawlID)

	assert.Empty(t, crawlingStatus)
	assert.Equal(t, expectedError, err)
	mockController.AssertNumberOfCalls(t, "GetCrawlingStatusFromDBFromCrawlID", 1)
}

func TestIsCurrentlyBeingCrawledReturnsTrueForAnActiveCrawl(t *testing.T) {
	mockController := &controller.MockCntrInterface{}

	activeCrawlingStatus := common.CrawlingStatus{
		OriginalCrawlTarget: "steamID",
		TotalUsersToCrawl:   140,
		UsersCrawled:        95,
	}
	mockController.On("GetCrawlingStatusFromDBFromCrawlID", mock.Anything, mock.Anything).Return(activeCrawlingStatus, nil)

	isActive, username, err := IsCurrentlyBeingCrawled(mockController, "crawlID")

	assert.True(t, isActive)
	assert.Equal(t, activeCrawlingStatus.OriginalCrawlTarget, username)
	assert.Nil(t, err)
}

func TestIsCurrentlyBeingCrawledReturnsAnErrorWhenAnErrorISReturnedSearchingForCrawlingStatus(t *testing.T) {
	mockController := &controller.MockCntrInterface{}

	randomError := errors.New("random error")
	mockController.On("GetCrawlingStatusFromDBFromCrawlID", mock.Anything, mock.Anything).Return(common.CrawlingStatus{}, randomError)

	isActive, username, err := IsCurrentlyBeingCrawled(mockController, "crawlID")

	assert.False(t, isActive)
	assert.Equal(t, "", username)
	assert.NotNil(t, err)
}

func TestIsCurrentlyBeingCrawledReturnsAnEmptyUsernameAndNoErrorWhenTheCrawlIsNotFinished(t *testing.T) {
	mockController := &controller.MockCntrInterface{}
	activeCrawlingStatus := common.CrawlingStatus{
		OriginalCrawlTarget: "steamID",
		TotalUsersToCrawl:   140,
		UsersCrawled:        140,
	}

	mockController.On("GetCrawlingStatusFromDBFromCrawlID", mock.Anything, mock.Anything).Return(activeCrawlingStatus, nil)

	isActive, username, err := IsCurrentlyBeingCrawled(mockController, "crawlID")

	assert.False(t, isActive)
	assert.Equal(t, "", username)
	assert.Nil(t, err)
}

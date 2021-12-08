// Code generated by mockery v2.9.4. DO NOT EDIT.

package controller

import (
	common "github.com/neosteamfriendgraphing/common"
	amqp "github.com/streadway/amqp"

	dtos "github.com/neosteamfriendgraphing/common/dtos"

	mock "github.com/stretchr/testify/mock"
)

// CntrInterface is an autogenerated mock type for the CntrInterface type
type MockCntrInterface struct {
	mock.Mock
}

// CallGetFriends provides a mock function with given fields: steamID
func (_m *MockCntrInterface) CallGetFriends(steamID string) ([]string, error) {
	ret := _m.Called(steamID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(steamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(steamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CallGetOwnedGames provides a mock function with given fields: steamID
func (_m *MockCntrInterface) CallGetOwnedGames(steamID string) (common.GamesOwnedResponse, error) {
	ret := _m.Called(steamID)

	var r0 common.GamesOwnedResponse
	if rf, ok := ret.Get(0).(func(string) common.GamesOwnedResponse); ok {
		r0 = rf(steamID)
	} else {
		r0 = ret.Get(0).(common.GamesOwnedResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(steamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CallGetPlayerSummaries provides a mock function with given fields: steamIDList
func (_m *MockCntrInterface) CallGetPlayerSummaries(steamIDList string) ([]common.Player, error) {
	ret := _m.Called(steamIDList)

	var r0 []common.Player
	if rf, ok := ret.Get(0).(func(string) []common.Player); ok {
		r0 = rf(steamIDList)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Player)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(steamIDList)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConsumeFromJobsQueue provides a mock function with given fields:
func (_m *MockCntrInterface) ConsumeFromJobsQueue() (<-chan amqp.Delivery, error) {
	ret := _m.Called()

	var r0 <-chan amqp.Delivery
	if rf, ok := ret.Get(0).(func() <-chan amqp.Delivery); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan amqp.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCrawlingStatsFromDataStore provides a mock function with given fields: crawlID
func (_m *MockCntrInterface) GetCrawlingStatsFromDataStore(crawlID string) (common.CrawlingStatus, error) {
	ret := _m.Called(crawlID)

	var r0 common.CrawlingStatus
	if rf, ok := ret.Get(0).(func(string) common.CrawlingStatus); ok {
		r0 = rf(crawlID)
	} else {
		r0 = ret.Get(0).(common.CrawlingStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(crawlID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGraphableDataFromDataStore provides a mock function with given fields: steamID
func (_m *MockCntrInterface) GetGraphableDataFromDataStore(steamID string) (dtos.GetGraphableDataForUserDTO, error) {
	ret := _m.Called(steamID)

	var r0 dtos.GetGraphableDataForUserDTO
	if rf, ok := ret.Get(0).(func(string) dtos.GetGraphableDataForUserDTO); ok {
		r0 = rf(steamID)
	} else {
		r0 = ret.Get(0).(dtos.GetGraphableDataForUserDTO)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(steamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserFromDataStore provides a mock function with given fields: steamID
func (_m *MockCntrInterface) GetUserFromDataStore(steamID string) (common.UserDocument, error) {
	ret := _m.Called(steamID)

	var r0 common.UserDocument
	if rf, ok := ret.Get(0).(func(string) common.UserDocument); ok {
		r0 = rf(steamID)
	} else {
		r0 = ret.Get(0).(common.UserDocument)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(steamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsernamesForSteamIDs provides a mock function with given fields: steamIDs
func (_m *MockCntrInterface) GetUsernamesForSteamIDs(steamIDs []string) (map[string]string, error) {
	ret := _m.Called(steamIDs)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func([]string) map[string]string); ok {
		r0 = rf(steamIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(steamIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishToJobsQueue provides a mock function with given fields: jobJSON
func (_m *MockCntrInterface) PublishToJobsQueue(jobJSON []byte) error {
	ret := _m.Called(jobJSON)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(jobJSON)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveCrawlingStatsToDataStore provides a mock function with given fields: currentLevel, crawlingStatus
func (_m *MockCntrInterface) SaveCrawlingStatsToDataStore(currentLevel int, crawlingStatus common.CrawlingStatus) (bool, error) {
	ret := _m.Called(currentLevel, crawlingStatus)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int, common.CrawlingStatus) bool); ok {
		r0 = rf(currentLevel, crawlingStatus)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, common.CrawlingStatus) error); ok {
		r1 = rf(currentLevel, crawlingStatus)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUserToDataStore provides a mock function with given fields: _a0
func (_m *MockCntrInterface) SaveUserToDataStore(_a0 dtos.SaveUserDTO) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(dtos.SaveUserDTO) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dtos.SaveUserDTO) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

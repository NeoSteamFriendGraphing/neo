package datastructures

import "github.com/neosteamfriendgraphing/common"

type GetProcessedGraphDataDTO struct {
	Status        string                `json:"status"`
	UserGraphData common.UsersGraphData `json:"usergraphdata"`
}

type AddUserEvent struct {
	SteamID     string `json:"steamid"`
	PersonaName string `json:"personaname"`
	ProfileURL  string `json:"profileurl"`
	Avatar      string `json:"avatar"`
	CountryCode string `json:"countrycode"`
	CrawlTime   int64  `json:"crawltime"`
}

type ShortestDistanceInfo struct {
	CrawlIDs         []string              `json:"crawlids"`
	FirstUser        common.UserDocument   `json:"firstuser"`
	SecondUser       common.UserDocument   `json:"seconduser"`
	ShortestDistance []common.UserDocument `json:"shortestdistance"`
	TotalNetworkSpan int                   `json:"totalnetworkspan"`
	TimeStarted      int64                 `json:"timestarted"`
}

type GetShortestDistanceInfoDataInputDTO struct {
	CrawlIDs []string `json:"crawlids"`
}

type FinishedCrawlWithItsUser struct {
	CrawlingStatus common.CrawlingStatus `json:"crawlingstatus"`
	User           common.UserDocument   `json:"user"`
}
type GetFinishedCrawlsDTO struct {
	Status                     string                     `json:"status"`
	AllFinishedCrawlsWithUsers []FinishedCrawlWithItsUser `json:"crawls"`
}

type GetFinishedShortestDistanceCrawlsDTO struct {
	Status         string                 `json:"status"`
	CrawlingStatus []ShortestDistanceInfo `json:"crawlingstatus"`
}

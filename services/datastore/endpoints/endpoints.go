package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IamCathal/neo/services/datastore/app"
	"github.com/IamCathal/neo/services/datastore/configuration"
	"github.com/IamCathal/neo/services/datastore/controller"
	"github.com/IamCathal/neo/services/datastore/datastructures"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/neosteamfriendgraphing/common"
	"github.com/neosteamfriendgraphing/common/dtos"
	"github.com/neosteamfriendgraphing/common/util"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Endpoints struct {
	Cntr controller.CntrInterface
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
// Taken from https://blog.questionable.services/article/guide-logging-middleware-go/
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (endpoints *Endpoints) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/status", endpoints.Status).Methods("POST")
	apiRouter.HandleFunc("/saveuser", endpoints.SaveUser).Methods("POST")
	apiRouter.HandleFunc("/insertgame", endpoints.InsertGame).Methods("POST")
	apiRouter.HandleFunc("/getuser/{steamid}", endpoints.GetUser).Methods("GET")
	apiRouter.HandleFunc("/getdetailsforgames", endpoints.GetDetailsForGames).Methods("POST")
	apiRouter.HandleFunc("/savecrawlingstats", endpoints.SaveCrawlingStatsToDB).Methods("POST")
	apiRouter.HandleFunc("/hasbeencrawledbefore", endpoints.HasBeenCrawledBefore).Methods("POST")
	apiRouter.HandleFunc("/getcrawlingstatus/{crawlid}", endpoints.GetCrawlingStatus).Methods("GET")
	apiRouter.HandleFunc("/getgraphabledata/{steamid}", endpoints.GetGraphableData).Methods("GET")
	apiRouter.HandleFunc("/getusernamesfromsteamids", endpoints.GetUsernamesFromSteamIDs).Methods("POST")
	apiRouter.HandleFunc("/saveprocessedgraphdata/{crawlid}", endpoints.SaveProcessedGraphData).Methods("POST")
	apiRouter.HandleFunc("/getprocessedgraphdata/{crawlid}", endpoints.GetProcessedGraphData).Methods("POST")
	apiRouter.Use(endpoints.LoggingMiddleware)

	wsRouter := r.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("/newuserstream", endpoints.NewUserStream).Methods("GET")

	return r
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (endpoints *Endpoints) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				vars := mux.Vars(r)
				w.WriteHeader(http.StatusInternalServerError)
				response := struct {
					Error string `json:"error"`
				}{
					fmt.Sprintf("Give the code monkeys this ID: '%s'", vars["requestID"]),
				}
				json.NewEncoder(w).Encode(response)

				requestStartTime, timeParseErr := strconv.ParseInt(vars["requestStartTime"], 10, 64)
				if timeParseErr != nil {
					configuration.Logger.Fatal(fmt.Sprintf("%v", err),
						zap.String("requestID", vars["requestID"]),
						zap.Int("status", http.StatusInternalServerError),
						zap.Int64("duration", util.GetCurrentTimeInMs()-requestStartTime),
						zap.String("path", r.URL.EscapedPath()),
					)
					panic(timeParseErr)
				}

				configuration.Logger.Error(fmt.Sprintf("%v", err),
					zap.String("requestID", vars["requestID"]),
					zap.Int("status", http.StatusInternalServerError),
					zap.Int64("duration", util.GetCurrentTimeInMs()-requestStartTime),
					zap.String("path", r.URL.EscapedPath()),
				)
			}
		}()

		vars := mux.Vars(r)

		identifier := ksuid.New()
		vars["requestID"] = identifier.String()

		requestStartTime := time.Now().UnixNano() / int64(time.Millisecond)
		vars["requestStartTime"] = strconv.Itoa(int(requestStartTime))

		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		configuration.Logger.Info("served content",
			zap.String("requestID", vars["requestID"]),
			zap.Int("status", wrapped.status),
			zap.Int64("duration", util.GetCurrentTimeInMs()-requestStartTime),
			zap.String("path", r.URL.EscapedPath()),
		)

		urlPathBasic := ""
		urlPath := strings.Split(r.URL.EscapedPath(), "/")

		if len(urlPath) > 1 {
			urlPathBasic = urlPath[1]
		} else {
			urlPathBasic = "/"
		}
		// TODO change from blocking to async
		writeAPI := configuration.InfluxDBClient.WriteAPIBlocking(os.Getenv("ORG"), os.Getenv("ENDPOINT_LATENCIES_BUCKET"))
		point := influxdb2.NewPointWithMeasurement("endpointLatencies").
			AddTag("path", urlPathBasic).
			AddTag("service", "datastore").
			AddField("latency", util.GetCurrentTimeInMs()-requestStartTime).
			SetTime(time.Now())
		err := writeAPI.WritePoint(context.Background(), point)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func (endpoints *Endpoints) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authenticate JWT
		next.ServeHTTP(w, r)
	})
}

func (endpoints *Endpoints) SaveUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	saveUserDTO := datastructures.SaveUserDTO{}

	err := json.NewDecoder(r.Body).Decode(&saveUserDTO)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		LogBasicErr(err, r, http.StatusBadRequest)
		return
	}

	crawlingStats := datastructures.CrawlingStatus{
		CrawlID:             saveUserDTO.CrawlID,
		OriginalCrawlTarget: saveUserDTO.OriginalCrawlTarget,
		MaxLevel:            saveUserDTO.MaxLevel,
		TotalUsersToCrawl:   len(saveUserDTO.User.FriendIDs),
	}

	err = app.SaveCrawlingStatsToDB(endpoints.Cntr, saveUserDTO.CurrentLevel, crawlingStats)
	if err != nil {
		configuration.Logger.Sugar().Errorf("failed to save crawling stats to DB: %+v", err)
		util.SendBasicInvalidResponse(w, r, "cannot save crawling stats", vars, http.StatusBadRequest)
		return
	}

	err = app.SaveUserToDB(endpoints.Cntr, saveUserDTO.User)
	if err != nil {
		configuration.Logger.Sugar().Errorf("failed to save user to DB: %+v", err)
		util.SendBasicInvalidResponse(w, r, "cannot save user", vars, http.StatusBadRequest)
		return
	}

	// Then save the games to the game table

	configuration.Logger.Sugar().Infof("successfully saved user %s to db", saveUserDTO.User.AccDetails.SteamID)

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		"success",
		"very good",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) InsertGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bareGameInfo := datastructures.BareGameInfo{}

	err := json.NewDecoder(r.Body).Decode(&bareGameInfo)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		return
	}

	success, err := endpoints.Cntr.InsertGame(context.TODO(), bareGameInfo)
	if err != nil || success != true {
		util.SendBasicInvalidResponse(w, r, "couldn't insert game", vars, http.StatusBadRequest)
		LogBasicInfo("couldn't insert game", r, http.StatusBadRequest)
		return
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		"success",
		"very good",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) SaveCrawlingStatsToDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	crawlingStatusInput := datastructures.SaveCrawlingStatsDTO{}

	err := json.NewDecoder(r.Body).Decode(&crawlingStatusInput)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		LogBasicErr(err, r, http.StatusBadRequest)
		return
	}

	err = app.SaveCrawlingStatsToDB(endpoints.Cntr, crawlingStatusInput.CurrentLevel, crawlingStatusInput.CrawlingStatus)
	if err != nil {
		LogBasicErr(err, r, http.StatusBadRequest)
		util.SendBasicInvalidResponse(w, r, "cannot save crawling stats", vars, http.StatusBadRequest)
		return
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		"success",
		"very good",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) HasBeenCrawledBefore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	crawlDetails := datastructures.HasBeenCrawledBeforeInputDTO{}

	err := json.NewDecoder(r.Body).Decode(&crawlDetails)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		return
	}

	if isValidFormat := util.IsValidFormatSteamID(crawlDetails.SteamID); !isValidFormat {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		return
	}

	hasBeenCrawled, err := endpoints.Cntr.HasUserBeenCrawledBeforeAtLevel(context.TODO(), crawlDetails.Level, crawlDetails.SteamID)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Could not lookup crawling status", vars, http.StatusBadRequest)
		configuration.Logger.Sugar().Fatalf("couldn't lookup has user been crawled before: %+v", err)
		panic(err)
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		"success",
		"",
	}

	if hasBeenCrawled {
		response.Message = "does exist"
	} else {
		response.Message = "does not exist"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Validate steamid
	if isValid := util.IsValidFormatSteamID(vars["steamid"]); !isValid {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		LogBasicInfo("invalid steamID given", r, http.StatusBadRequest)
		return
	}

	user, err := app.GetUserFromDB(endpoints.Cntr, vars["steamid"])
	if err == mongo.ErrNoDocuments {
		util.SendBasicInvalidResponse(w, r, "user does not exist", vars, http.StatusNotFound)
		LogBasicInfo("user was not found in DB", r, http.StatusNotFound)
		return
	}

	if err != nil {
		util.SendBasicInvalidResponse(w, r, "couldn't get user", vars, http.StatusBadRequest)
		LogBasicInfo(fmt.Sprintf("couldn't get user: %+v", err), r, http.StatusBadRequest)
		return
	}
	response := struct {
		Status string              `json:"status"`
		User   common.UserDocument `json:"user"`
	}{
		"success",
		user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (endpoints *Endpoints) GetDetailsForGames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	gamesInput := datastructures.GetDetailsForGamesDTO{}

	err := json.NewDecoder(r.Body).Decode(&gamesInput)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		return
	}

	if len(gamesInput.GameIDs) == 0 || len(gamesInput.GameIDs) > 20 {
		util.SendBasicInvalidResponse(w, r, "Can only request 1-20 games in a request", vars, http.StatusBadRequest)
		return
	}

	gameDetails, err := endpoints.Cntr.GetDetailsForGames(context.TODO(), gamesInput.GameIDs)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Error retrieving games", vars, http.StatusBadRequest)
		configuration.Logger.Sugar().Errorf("error retrieving games: %+v", err)
		return
	}
	if len(gameDetails) == 0 {
		gameDetails = []datastructures.BareGameInfo{}
	}

	response := struct {
		Status string                        `json:"status"`
		Games  []datastructures.BareGameInfo `json:"games"`
	}{
		"success",
		gameDetails,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) GetCrawlingStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_, err := ksuid.Parse(vars["crawlid"])
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "invalid crawlid", vars, http.StatusNotFound)
		LogBasicInfo("invalid crawlid given", r, http.StatusNotFound)
		return
	}
	crawlingStatus, err := app.GetCrawlingStatsFromDBFromCrawlID(endpoints.Cntr, vars["crawlid"])
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "couldn't get crawling status", vars, http.StatusNotFound)
		LogBasicInfo("couldn't get crawling status", r, http.StatusNotFound)
		return
	}
	response := datastructures.GetCrawlingStatusDTO{
		Status:         "success",
		CrawlingStatus: crawlingStatus,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) GetGraphableData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Validate steamid
	if isValid := util.IsValidFormatSteamID(vars["steamid"]); !isValid {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		return
	}

	user, err := endpoints.Cntr.GetUser(context.TODO(), vars["steamid"])
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "couldn't get user", vars, http.StatusNotFound)
		return
	}

	graphableDataForUser := dtos.GetGraphableDataForUserDTO{
		Username:  user.AccDetails.Personaname,
		SteamID:   user.AccDetails.SteamID,
		FriendIDs: user.FriendIDs,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(graphableDataForUser)
}

func (endpoints *Endpoints) GetUsernamesFromSteamIDs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	steamIDsInput := dtos.GetUsernamesFromSteamIDsInputDTO{}

	err := json.NewDecoder(r.Body).Decode(&steamIDsInput)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		LogBasicErr(err, r, http.StatusBadRequest)
		return
	}

	// Validate all given steamids
	for _, steamID := range steamIDsInput.SteamIDs {
		if isValid := util.IsValidFormatSteamID(steamID); !isValid {
			util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
			return
		}
	}

	steamIDsToUsernames, err := endpoints.Cntr.GetUsernames(context.TODO(), steamIDsInput.SteamIDs)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "couldn't get usernames from steamIDs", vars, http.StatusNotFound)
		configuration.Logger.Error(err.Error())
		return
	}
	configuration.Logger.Sugar().Infof("retrieved %d usernames from steamIDs", len(steamIDsToUsernames))

	returnJSON := dtos.GetUsernamesFromSteamIDsDTO{}
	for key, val := range steamIDsToUsernames {
		currentMapping := dtos.SteamIDAndUsername{
			SteamID:  key,
			Username: val,
		}
		returnJSON.SteamIDAndUsername = append(returnJSON.SteamIDAndUsername, currentMapping)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(returnJSON)
}

func (endpoints *Endpoints) SaveProcessedGraphData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_, err := ksuid.Parse(vars["crawlid"])
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "invalid input", vars, http.StatusBadRequest)
		LogBasicInfo("invalid crawlid given", r, http.StatusNotFound)
		return
	}

	graphData := datastructures.UsersGraphData{}
	err = json.NewDecoder(r.Body).Decode(&graphData)
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Invalid input", vars, http.StatusBadRequest)
		return
	}

	success, err := endpoints.Cntr.SaveProcessedGraphData(vars["crawlid"], graphData)
	if success == false || err != nil {
		util.SendBasicInvalidResponse(w, r, "could not retrieve graph data", vars, http.StatusBadRequest)
		configuration.Logger.Sugar().Errorf("could not retrieve graph data: %+v", err)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		"success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) GetProcessedGraphData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_, err := ksuid.Parse(vars["crawlid"])
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "invalid input", vars, http.StatusBadRequest)
		LogBasicInfo("invalid crawlid given", r, http.StatusNotFound)
		return
	}

	usersProcessedGraphData, err := endpoints.Cntr.GetProcessedGraphData(vars["crawlid"])
	if err != nil {
		util.SendBasicInvalidResponse(w, r, "Couldn't get graph data", vars, http.StatusBadRequest)
		errMsg := fmt.Errorf("failed to get processed graph data: %+v", err)
		configuration.Logger.Error(errMsg.Error())
		return
	}
	response := datastructures.GetProcessedGraphDataDTO{
		Status:        "success",
		UserGraphData: usersProcessedGraphData,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (endpoints *Endpoints) Status(w http.ResponseWriter, r *http.Request) {
	req := common.UptimeResponse{
		Uptime: time.Since(configuration.ApplicationStartUpTime),
		Status: "operational",
	}
	jsonObj, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonObj))
}

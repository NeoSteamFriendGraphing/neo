package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IamCathal/neo/services/frontend/endpoints"
	"github.com/IamCathal/neo/services/frontend/util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logConfig := util.LoadLoggingConfig()
	endpoints := &endpoints.Endpoints{
		Logger:                 util.InitLogger(logConfig),
		ApplicationStartUpTime: time.Now(),
	}

	router := endpoints.SetupRouter()
	router.Handle("/static", http.NotFoundHandler())
	fs := http.FileServer(http.Dir(os.Getenv("STATIC_CONTENT_DIR")))
	router.PathPrefix("/").Handler(http.StripPrefix("/static", endpoints.DisallowFileBrowsing(fs)))

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("API_PORT")),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	endpoints.Logger.Info(fmt.Sprintf("frontend start up and serving requsts on %s:%s", util.GetLocalIPAddress(), os.Getenv("API_PORT")))
	log.Fatal(srv.ListenAndServe())
}

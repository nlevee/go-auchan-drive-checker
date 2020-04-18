package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nlevee/go-auchan-drive-checker/pkg/auchan"
	"github.com/nlevee/go-auchan-drive-checker/pkg/drivestate"
)

func handle(currentState *drivestate.DriveState) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "` + (*currentState).Dispo + `"}`)); err != nil {
			log.Fatal(err)
		}
		log.Printf("current state is : %v", (*currentState).IsActive)
	}
}

var tick = time.NewTicker(2 * time.Minute)
var done = make(chan bool)
var currentState = make(map[string]*drivestate.DriveState)

func addDriveHandler(driveId string) {
	config := auchan.NewConfig(driveId)
	currentState[driveId] = config.State
	go auchan.GetDriveState(config, tick, done)
}

func main() {
	auchanDriveId := flag.String("id", "", "The drive Id")
	listenHost := flag.String("host", "0.0.0.0", "Start a server and listen on this host")
	listenPort := flag.String("port", "", "Start a server and listen on this port")
	flag.Parse()

	if *auchanDriveId == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *listenPort != "" && *listenHost != "" {
		addDriveHandler(*auchanDriveId)
		r := mux.NewRouter()
		r.HandleFunc("/", handle(currentState[*auchanDriveId])).Methods(http.MethodGet)
		log.Fatal(http.ListenAndServe(*listenHost+":"+*listenPort, r))
	} else {
		config := auchan.NewConfig(*auchanDriveId)
		auchan.GetDriveState(config, tick, done)
	}
}

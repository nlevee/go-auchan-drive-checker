package main

import (
	"log"
	"net/http"
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
	currentState[driveId] = &drivestate.DriveState{
		IsActive: false,
		Dispo:    "",
	}
	config := auchan.DriveConfig{
		DriveId: driveId,
		State:   currentState[driveId],
	}
	go auchan.GetDriveState(config, tick, done)
}

func main() {
	auchanDriveId := "930"
	addDriveHandler(auchanDriveId)

	r := mux.NewRouter()
	r.HandleFunc("/", handle(currentState[auchanDriveId])).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8089", r))
}

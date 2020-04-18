// Package auchan provides function to get auchan drive disponibility
package auchan

import (
	"log"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/nlevee/go-auchan-drive-checker/pkg/drivestate"
	"github.com/nlevee/go-auchan-drive-checker/pkg/utils"
)

const (
	auchanDriveUrl = "https://www.auchandrive.fr/drive/mag/anything-"
)

type DriveConfig struct {
	DriveId string
	State   *drivestate.DriveState
}

// Create a new Drive config with driveId
func NewConfig(driveId string) DriveConfig {
	state := &drivestate.DriveState{
		IsActive: false,
		Dispo:    "",
	}
	return DriveConfig{
		DriveId: driveId,
		State:   state,
	}
}

func loadDriveState(config DriveConfig) (hasChanged bool, err error) {
	driveUrl := auchanDriveUrl + config.DriveId
	currentState := config.State

	log.Printf("Request uri : %v", driveUrl)

	doc, err := utils.LoadURL(driveUrl)
	if err != nil {
		return false, err
	}
	node := htmlquery.FindOne(doc, `//div[@class='next-slot__text-slot']`)
	if node != nil {
		if newDispo := htmlquery.InnerText(node); (*currentState).Dispo != newDispo {
			(*currentState).IsActive = true
			(*currentState).Dispo = newDispo
			log.Printf("Nouveau créneau %v", currentState)
			return true, nil
		}
	} else if (*currentState).IsActive {
		(*currentState).IsActive = false
		log.Printf("Aucun créneau pour le moment")
		return true, nil
	}
	return false, nil
}

func GetDriveState(config DriveConfig, tick *time.Ticker, done chan bool) {

	log.Printf("Démarrage du check de créneau Auchan Drive %v", config.DriveId)

	// premier appel sans attendre le premier tick
	if _, err := loadDriveState(config); err != nil {
		log.Print(err)
	}

	for {
		select {
		case <-tick.C:
			// a chaque tick du timer on lance une recherche de state
			if _, err := loadDriveState(config); err != nil {
				log.Print(err)
			}
		case <-done:
			log.Printf("Ticker stopped")
			tick.Stop()
			return
		}
	}
}

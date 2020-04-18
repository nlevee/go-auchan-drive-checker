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

func loadDriveState(driveUrl string, currentState *drivestate.DriveState) (hasChanged bool, err error) {
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
	driveUrl := auchanDriveUrl + config.DriveId
	currentState := config.State

	log.Printf("Démarrage du check de créneau Auchan Drive %v", config.DriveId)
	log.Printf("Request uri : %v", driveUrl)

	// premier appel sans attendre le premier tick
	if _, err := loadDriveState(driveUrl, currentState); err != nil {
		log.Print(err)
	}

	for {
		select {
		case <-tick.C:
			// a chaque tick du timer on lance une recherche de state
			if _, err := loadDriveState(driveUrl, currentState); err != nil {
				log.Print(err)
			}
		case <-done:
			log.Printf("Ticker stopped")
			tick.Stop()
			return
		}

	}
}

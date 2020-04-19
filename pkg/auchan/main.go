// Package auchan provides function to get auchan drive disponibility
package auchan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/nlevee/go-auchan-drive-checker/pkg/drivestate"
	"github.com/nlevee/go-auchan-drive-checker/pkg/utils"
)

const (
	auchanDriveURL = "https://www.auchandrive.fr/drive/mag/anything-"
)

type DriveConfig struct {
	DriveID string
	State   *drivestate.DriveState
}

type DriveStore struct {
	DriveID string
	Name    string
}

// NewConfig Create a new Drive config with driveId
func NewConfig(driveID string) DriveConfig {
	state := &drivestate.DriveState{
		IsActive: false,
		Dispo:    "",
	}
	return DriveConfig{
		DriveID: driveID,
		State:   state,
	}
}

type store struct {
	Features []struct {
		Properties struct {
			Store_ID string
			Name     string
		}
	}
}

// GetStoreByPostalCode fetch stores by postal code
func GetStoreByPostalCode(postalCode string) ([]DriveStore, error) {
	stores := []DriveStore{}

	cities, err := utils.GetCitiesByPostalCode(postalCode)
	if err != nil || len(cities) == 0 {
		return stores, err
	}

	woosKey := os.Getenv("WOOS_KEY")
	woosDomain := os.Getenv("WOOS_DOMAIN")
	if woosKey == "" || woosDomain == "" {
		log.Fatal("env var 'WOOS_KEY' and 'WOOS_DOMAIN' are required")
	}
	url := "https://api.woosmap.com/stores/search?key=" + woosKey + "&max_distance=20000&query=type:DRIVE"

	city := cities[0]

	requrl := url + "&lat=" + fmt.Sprintf("%f", city.Lat) + "&lng=" + fmt.Sprintf("%f", city.Lon)
	log.Print(requrl)
	req, err := http.NewRequest("GET", requrl, bytes.NewReader([]byte{}))
	if err != nil {
		log.Print(err)
		return stores, err
	}
	req.Header.Add("origin", "https://"+os.Getenv("WOOS_DOMAIN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(err, resp.Status, string(dump))
		return stores, err
	}

	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return stores, err
	}

	storeFound := store{}
	json.Unmarshal(bodyContent, &storeFound)

	for _, v := range storeFound.Features {
		stores = append(stores, DriveStore{
			DriveID: v.Properties.Store_ID,
			Name:    v.Properties.Name,
		})
	}

	return stores, nil
}

// GetStoreIDByPostalCode fetch storeIDs by postal code
func GetStoreIDByPostalCode(postalCode string) ([]string, error) {
	storeIds := []string{}

	stores, err := GetStoreByPostalCode(postalCode)
	if err != nil {
		return storeIds, err
	}

	for _, v := range stores {
		storeIds = append(storeIds, v.DriveID)
	}

	return storeIds, nil
}

func loadDriveState(config DriveConfig) (hasChanged bool, err error) {
	driveURL := auchanDriveURL + config.DriveID
	currentState := config.State

	log.Printf("Request uri : %v", driveURL)

	doc, err := utils.LoadHTMLURL(driveURL)
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
		log.Printf("Aucun créneau pour le moment")
		(*currentState).IsActive = false
		return true, nil
	} else {
		log.Printf("Aucun créneau pour le moment")
	}
	return false, nil
}

// LoadIntervalDriveState fetch each tick the drive state config
func LoadIntervalDriveState(config DriveConfig, tick *time.Ticker, done chan bool) {
	log.Printf("Démarrage du check de créneau Auchan Drive %v", config.DriveID)

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

// GetDriveState get the state of a drive
func GetDriveState(driveID string) *drivestate.DriveState {
	return drivestate.GetDriveState(driveID)
}

// NewDriveHandler add a new drive handler
func NewDriveHandler(driveID string) {
	config := NewConfig(driveID)
	drivestate.NewDriveState(driveID, config.State)

	tick := time.NewTicker(2 * time.Minute)
	done := make(chan bool)

	LoadIntervalDriveState(config, tick, done)
}

package main

import (
	"flag"
	"log"
	"os"

	"github.com/nlevee/go-auchan-drive-checker/internal/api"
	"github.com/nlevee/go-auchan-drive-checker/pkg/auchan"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	auchanDriveID := flag.String("id", "", "The drive Id")
	postalCode := flag.String("cp", "", "The Postal Code")
	listenHost := flag.String("host", "0.0.0.0", "Start a server and listen on this host")
	listenPort := flag.String("port", "", "Start a server and listen on this port")
	flag.Parse()

	// recherche du driveId si code postal
	if *auchanDriveID == "" && *postalCode != "" {
		storeIDs, _ := auchan.GetStoreIDByPostalCode(*postalCode)
		if len(storeIDs) > 0 {
			auchanDriveID = &storeIDs[0]
		} else {
			log.Fatal("no stores found")
		}
	}

	if *listenPort != "" && *listenHost != "" {
		if *auchanDriveID != "" {
			go auchan.NewDriveHandler(*auchanDriveID)
		}
		api.StartServer(*listenHost, *listenPort)
	} else {
		if *auchanDriveID == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		auchan.NewDriveHandler(*auchanDriveID)
	}
}

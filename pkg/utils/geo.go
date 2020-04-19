package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Cities expose a city structure
type Cities struct {
	Name string
	Lat  float64
	Lon  float64
}

type communes []struct {
	Code   string
	Nom    string
	Centre struct {
		Coordinates []float64
	}
}

// GetCitiesByPostalCode fetch cities center point by postal code
func GetCitiesByPostalCode(postalCode string) ([]Cities, error) {
	cities := []Cities{}

	url := "https://geo.api.gouv.fr/communes?codePostal=" + postalCode + "&fields=centre"

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return cities, err
	}
	defer resp.Body.Close()

	bodyContent, _ := ioutil.ReadAll(resp.Body)
	var foundCities communes
	json.Unmarshal(bodyContent, &foundCities)

	for _, v := range foundCities {
		cities = append(cities, Cities{
			Name: v.Nom,
			Lon:  v.Centre.Coordinates[0],
			Lat:  v.Centre.Coordinates[1],
		})
	}

	return cities, nil
}

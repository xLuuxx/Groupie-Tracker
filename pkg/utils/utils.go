package utils

import (
	"Groupie-Tracker/internal/API"
	"Groupie-Tracker/internal/structure"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func IsInList(slice []string, element string) bool {
	/*
		This function checks if an element is in a slice.
		params: slice []string - the slice to check, element string - the element to check
		return: bool - true if the element is in the slice, false otherwise
	*/
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func GenerateGeoJson(locationsCity []string, locationsCoordinates [][]float32, locationsDates []string) (string, error) {
	/*
		This function generates a GeoJSON string from the locations, coordinates and dates.
		params: locationsCity []string - the locations, locationsCoordinates [][]float32 - the coordinates, locationsDates []string - the dates
		return: string - the GeoJSON string, error
	*/
	geoJSON := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": []map[string]interface{}{},
	}

	for i, location := range locationsCity {
		parsedDate, err := time.Parse("02-01-2006", locationsDates[i])
		if err != nil {
			return "", err
		}

		formattedDate := fmt.Sprintf("%02d %s %d", parsedDate.Day(), parsedDate.Month().String()[:3], parsedDate.Year())

		if !IsInList(locationsCity[:i], location) {
			geoJSON["features"] = append(geoJSON["features"].([]map[string]interface{}), map[string]interface{}{
				"type": "Feature",
				"geometry": map[string]interface{}{
					"type":        "Point",
					"coordinates": locationsCoordinates[i],
				},
				"properties": map[string]interface{}{
					"title":       location,
					"description": formattedDate,
				},
			})
		} else {
			for j, feature := range geoJSON["features"].([]map[string]interface{}) {
				if feature["properties"].(map[string]interface{})["title"] == location {
					description := feature["properties"].(map[string]interface{})["description"].(string)
					geoJSON["features"].([]map[string]interface{})[j]["properties"].(map[string]interface{})["description"] = description + ", " + formattedDate
				}
			}
		}
	}

	jsonData, err := json.Marshal(geoJSON)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func GeoJsonByArtistID(id int) (string, error) {
	/*
		This function gets the GeoJSON data for an artist by ID.
		params: id int - the ID of the artist
		return: string - the GeoJSON string, error
	*/
	locations, dates, err := API.GetArtistRelations(id)
	if err != nil {
		return "", err
	}

	var locationsCity []string
	var locationsCoordinates [][]float32

	for _, location := range locations {
		name, coordinates, err := API.GetDataCitybyAPI(location)
		if err != nil {
			return "", err
		}
		locationsCity = append(locationsCity, name)
		locationsCoordinates = append(locationsCoordinates, coordinates)
	}

	geoJson, err := GenerateGeoJson(locationsCity, locationsCoordinates, dates)

	return geoJson, err
}

func GetSpotifyURLByID(id string) string {
	// Read the file
	data, err := os.ReadFile("data/SpotifyID.json")
	if err != nil {
		return ""
	}

	// Unmarshal the JSON data into a slice of SpotifyID structs
	var spotifyIDs []structure.SpotifyID
	err = json.Unmarshal(data, &spotifyIDs)
	if err != nil {
		return ""
	}

	// Search for the entry with the matching ID
	for _, item := range spotifyIDs {
		if item.ID == id {
			return item.SpotifyURL
		}
	}

	return ""
}

func ValidInput(input string) bool {
	// Check if the input is valid
	for _, r := range input {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

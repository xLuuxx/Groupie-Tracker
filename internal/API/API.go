package API

import (
	"Groupie-Tracker/internal/structure"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

var Artists []structure.ArtistAPI
var artistTime time.Time

func GetArtists() (*[]structure.ArtistAPI, error) {
	/*
		This function updates the artists slice by calling the GetArtists function.
		return: a slice of ArtistAPI structs, error
	*/

	// If the artists slice is empty or if the last update was more than 30 minutes ago
	if Artists == nil || time.Since(artistTime) > 30*time.Minute {
		return UpdateArtists()
	}
	return &Artists, nil
}

func UpdateArtists() (*[]structure.ArtistAPI, error) {
	/*
		This function makes a GET request to the groupietrackers API to get all artists.
		return: a slice of ArtistAPI structs, error
	*/
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Artists); err != nil {
		return nil, err
	}
	artistTime = time.Now()
	return &Artists, nil
}

func GetArtistByID(id int) (*structure.ArtistAPI, error) {
	/*
		This function makes a GET request to the groupietrackers API to get an artist by ID.
		params: id int - the ID of the artist
		return: a pointer to an ArtistAPI struct, error
	*/

	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}

	return &(*artists)[id-1], nil

}

func GetArtistRelations(id int) ([]string, []string, error) {
	/*
		This function makes a GET request to the groupietrackers API to get the relations of an artist by ID.
		params: id int - the ID of the artist
		return: two slices of strings - the location and dates, error
	*/
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var relationData structure.RelationAPI
	if err := json.NewDecoder(resp.Body).Decode(&relationData); err != nil {
		return nil, nil, err
	}

	var locations []string
	var dates []string
	for location, dateList := range relationData.DatesLocations {
		for _, date := range dateList {
			locations = append(locations, location)
			dates = append(dates, date)
		}
	}
	return locations, dates, nil
}

func GetDataCitybyAPI(city string) (string, []float32, error) {
	/*
		This function makes a GET request to the mapbox API to get the city data by city name.
		params: city string - the name of the city
		return: a string - the name of the city, a slice of float32 - the coordinates, error
	*/
	resp, err := http.Get("https://api.mapbox.com/search/geocode/v6/forward?q=" + city + "&limit=1&access_token=pk.eyJ1Ijoib2loYSIsImEiOiJjbTYweGEycHEwaTJsMmxzNTBjanU0OWVvIn0.IZDNi1lmxBbldfGE7VhSrA")
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var mapboxResp structure.MapboxResponse
	if err := json.NewDecoder(resp.Body).Decode(&mapboxResp); err != nil {
		return "", nil, err
	}

	if len(mapboxResp.Features) == 0 {
		return "", nil, errors.New("no features found")
	}

	name := mapboxResp.Features[0].Properties.Name
	coordinates := []float32{
		mapboxResp.Features[0].Properties.Coordinates.Longitude,
		mapboxResp.Features[0].Properties.Coordinates.Latitude,
	}
	country := mapboxResp.Features[0].Properties.Context.Country.Name

	return name + " " + country, coordinates, nil
}

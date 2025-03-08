package back

import (
	"Groupie-Tracker/internal/API"
	"Groupie-Tracker/internal/structure"
	"Groupie-Tracker/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func artist(w http.ResponseWriter, r *http.Request) {
	/*
		This function is the handler for the /artist route.
		It gets the artist ID from the URL query, gets the artist info and relations from the API,
		formats the data and renders the artist template.
	*/
	idData := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx := context.WithValue(r.Context(), "status", http.StatusInternalServerError)
		errorTemplate(w, r.WithContext(ctx))
		return
	}

	if id < 1 || id > len(API.Artists) {
		w.WriteHeader(http.StatusInternalServerError)
		ctx := context.WithValue(r.Context(), "status", http.StatusInternalServerError)
		errorTemplate(w, r.WithContext(ctx))
		return
	}

	artistInfo, err := API.GetArtistByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx := context.WithValue(r.Context(), "status", http.StatusInternalServerError)
		errorTemplate(w, r.WithContext(ctx))
		return
	}

	locations, _, err := API.GetArtistRelations(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx := context.WithValue(r.Context(), "status", http.StatusInternalServerError)
		errorTemplate(w, r.WithContext(ctx))
		return
	}

	for i, location := range locations {
		locations[i] = formatLocation(location)
	}

	artistData := structure.ArtistData{
		ID:           artistInfo.ID,
		Image:        artistInfo.Image,
		Name:         artistInfo.Name,
		Members:      artistInfo.Members,
		CreationDate: artistInfo.CreationDate,
		FirstAlbum:   artistInfo.FirstAlbum,
		Locations:    locations,
		GeoJson:      getGeoJson(id),
		SpotifyID:    utils.GetSpotifyURLByID(strconv.Itoa(artistInfo.ID)),
	}

	renderTemplate(w, "artist", artistData)
}

func formatLocation(location string) string {
	/*
		This function formats the location string to be displayed in the artist template.
	*/
	return strings.Title(strings.ReplaceAll(strings.ReplaceAll(location, "_", " "), "-", " "))
}

func getGeoJson(id int) string {
	/*
		This function gets the geoJson data for the artist by ID.
	*/
	geoJson, err := utils.GeoJsonByArtistID(id)
	if err != nil {
		fmt.Println(err) //TODO: log errorTemplate
	}
	return geoJson
}

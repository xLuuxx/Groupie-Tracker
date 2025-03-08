package back

import (
	"Groupie-Tracker/internal/API"
	"Groupie-Tracker/internal/structure"
	"Groupie-Tracker/pkg/utils"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func search(w http.ResponseWriter, r *http.Request) {
	/*
		This function is the handler for the /search route.
		It gets the search parameters from the URL query, filters the artists based on the parameters,
		and renders the search results template.
	*/
	Name := r.URL.Query().Get("name")
	GroupCreation := r.URL.Query().Get("Group-creation")
	GroupLocation := r.URL.Query().Get("Group-location")
	GroupAlbum := r.URL.Query().Get("Group-album")
	var GroupMember []string

	for i := 1; i < 9; i++ {
		if r.URL.Query().Get("Group-member-"+strconv.Itoa(i)) != "" {
			GroupMember = append(GroupMember, r.URL.Query().Get("Group-member-"+strconv.Itoa(i)))
		}
	}

	if !utils.ValidInput(Name) || !utils.ValidInput(GroupCreation) || !utils.ValidInput(GroupLocation) || !utils.ValidInput(GroupAlbum) {
		data := structure.PageData{
			Title:   "Invalid input",
			Message: "Invalid input, please enter only letters and numbers.",
		}
		renderTemplate(w, "error", data)
		return
	}

	artists, err := API.GetArtists()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx := context.WithValue(r.Context(), "status", http.StatusInternalServerError)
		errorTemplate(w, r.WithContext(ctx))
		return
	}
	var filteredArtists []structure.ArtistData
	for _, artistFiltered := range *artists {
		if Name != "" && !strings.Contains(strings.ToLower(artistFiltered.Name), strings.ToLower(Name)) {
			continue
		}

		if GroupCreation != "1987" && GroupCreation != strconv.Itoa(artistFiltered.CreationDate) {
			continue
		}

		if len(GroupMember) > 0 && !utils.IsInList(GroupMember, strconv.Itoa(len(artistFiltered.Members))) {
			continue
		}

		if GroupAlbum != "" && !strings.Contains(strings.ToLower(artistFiltered.FirstAlbum), strings.ToLower(GroupAlbum)) {
			continue
		}

		locationConcert, _, _ := API.GetArtistRelations(artistFiltered.ID)
		locationMatch := false
		for _, location := range locationConcert {
			if GroupLocation != "" && strings.Contains(strings.ToLower(location), strings.ToLower(GroupLocation)) {
				locationMatch = true
				break
			}
		}
		if GroupLocation != "" && !locationMatch {
			continue
		}

		filteredArtists = append(filteredArtists, structure.ArtistData{
			ID:           artistFiltered.ID,
			Image:        artistFiltered.Image,
			Name:         artistFiltered.Name,
			Members:      artistFiltered.Members,
			CreationDate: artistFiltered.CreationDate,
			FirstAlbum:   artistFiltered.FirstAlbum,
			Locations:    []string{artistFiltered.Locations},
		})
	}
	data := structure.PageData{
		Title:       "Search Results",
		ArtistsData: filteredArtists,
		Message:     "",
	}

	if len(filteredArtists) == 0 {
		data.Message = "No results found for your search."
	}

	renderTemplate(w, "homepage", data)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	/*
		This function is the handler for the / and /homepage routes.
		It gets the list of artists from the API and renders the homepage template.
	*/
	items := make([]int, 52)
	for i := 0; i < 52; i++ {
		items[i] = i + 1
	}

	artists, err := API.GetArtists()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx := context.WithValue(r.Context(), "status", http.StatusInternalServerError)
		errorTemplate(w, r.WithContext(ctx))
		return
	}

	var listArtists []structure.ArtistData
	for _, artists := range *artists {
		listArtists = append(listArtists, structure.ArtistData{
			ID:           artists.ID,
			Image:        artists.Image,
			Name:         artists.Name,
			Members:      artists.Members,
			CreationDate: artists.CreationDate,
			FirstAlbum:   artists.FirstAlbum,
			Locations:    []string{artists.Locations},
		})
	}

	data := structure.PageData{
		Title:       "HOMEPAGE",
		Items:       items,
		ArtistsData: listArtists,
	}
	renderTemplate(w, "homepage", data)
}

func errorTemplate(w http.ResponseWriter, r *http.Request) {
	/*
		This function is the handler for the /error route.
		It gets the status code from the context and renders the error template.
	*/
	statusCode := http.StatusInternalServerError
	if status := r.Context().Value("status"); status != nil {
		statusCode = status.(int)
	}

	var message string
	switch statusCode {
	case http.StatusNotFound:
		message = "The page you are looking for does not exist."
	case http.StatusBadRequest:
		message = "Invalid input, please enter only letters and numbers."
	default:
		message = "An unexpected error occurred."
	}

	data := structure.PageData{
		Title:   fmt.Sprintf("Error %d", statusCode),
		Message: message,
	}

	renderTemplate(w, "error", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	/*
		This function is used to render the templates using the data sent
	*/
	t, err := template.ParseFiles("web/template/" + tmpl + ".gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

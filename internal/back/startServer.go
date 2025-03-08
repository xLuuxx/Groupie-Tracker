package back

import (
	"Groupie-Tracker/internal/API"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func StartServer() {
	/*
		This function starts the server and initializes the routes.
		also updates the artists slice by calling the API.GetArtists function.
		and handles the error if the API call fails.
	*/
	_, err := API.GetArtists()
	if err != nil {
		log.Fatal("Error getting artists:", err)
		return
	}

	http.HandleFunc("/", homepage)
	http.HandleFunc("/homepage", homepage)
	http.HandleFunc("/artist", artist)
	http.HandleFunc("/search", search)
	http.HandleFunc("/error", errorTemplate)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	fmt.Println("(http://localhost:8080) - Server started on port 8080")
	err = http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/" && path != "/homepage" && path != "/artist" && path != "/search" &&
			path != "/error" && !strings.HasPrefix(path, "/static/") {
			w.WriteHeader(http.StatusNotFound)
			ctx := context.WithValue(r.Context(), "status", http.StatusNotFound)
			errorTemplate(w, r.WithContext(ctx))
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

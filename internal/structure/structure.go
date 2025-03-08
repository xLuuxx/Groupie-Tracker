package structure

type PageData struct {
	Title        string
	Items        []int
	ItemsPerRow  int
	ArtistsData  []ArtistData
	Message      string
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
}

type ArtistAPI struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type ArtistData struct {
	ID           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    []string
	GeoJson      string
	SpotifyID    string
}

type RelationAPI struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MapboxResponse struct {
	/*
		This struct is used to decode the response from the mapbox API.
		Oui c'était chiant à faire
	*/
	Features []struct {
		Properties struct {
			Name        string `json:"name"`
			Coordinates struct {
				Longitude float32 `json:"longitude"`
				Latitude  float32 `json:"latitude"`
			} `json:"coordinates"`
			Context struct {
				Country struct {
					Name string `json:"name"`
				} `json:"country"`
			} `json:"context"`
		} `json:"properties"`
	} `json:"features"`
}

type SpotifyID struct {
	ID         string `json:"id"`
	SpotifyURL string `json:"spotifyUrl"`
}

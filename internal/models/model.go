package models

type Artists struct {
	ID            int                 `json:"id"`
	Image         string              `json:"image"`
	Name          string              `json:"name"`
	Members       []string            `json:"members"`
	CreationDate  int                 `json:"creationDate"`
	FirstAlbum    string              `json:"firstAlbum"`
	DateLocations map[string][]string `json:"datesLocations"`
}

type Relations struct {
	Index []struct {
		ID            int                 `json:"id"`
		DateLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type Error struct {
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}

type JWToken struct {
	Token string `json:"token"`
}

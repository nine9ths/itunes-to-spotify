package models

type Song struct {
	Name string
	Artist string
}

type SpotifySearch struct {
	Tracks SpotifyTrack `json:"tracks"`
}

type SpotifyTrack struct {
	Items []SpotifyItem `json:"items"`
}

type SpotifyItem struct {
	Name string `json:"name"`
	DurationMs int `json:"duration_ms"`
	URI string `json:"uri"`
}

type SpotifySearchSimple struct {
	URI string `json:"uri"`
}

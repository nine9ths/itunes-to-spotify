package spotify

import (
	"github.com/skyerus/itunes-to-spotify/pkg/models"
	"io"
)

type Service interface {
	ReadSongs(r io.Reader) ([]models.Song, error)
	SearchSong(song models.Song) (result models.SpotifySearchSimple, success bool, err error)
	SearchSongs(songs []models.Song) ([]models.SpotifySearchSimple, []models.Song, error)
	AddResultsToSpotifyPlaylist([]models.SpotifySearchSimple) error
	AddNonexistentToFile([]models.Song) error
	CreateSpotifyPlaylist(name string) error
}

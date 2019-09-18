package spotify

import (
	"github.com/skyerus/itunes-to-spotify/pkg/models"
	"io"
)

type Service interface {
	ReadSongs(r io.Reader) ([]models.Song, error)
	SearchSong(song models.Song) (result string, success bool, err error)
	SearchSongs(songs []models.Song) ([]string, []models.Song, error)
	AddResultsToSpotifyPlaylist(playlistObj models.SpotifyPlaylistObject, results []string) error
	AddNonexistentToFile(songs []models.Song, path string) error
	GetSpotifyUserObject() (models.SpotifyUserObject, error)
	CreateSpotifyPlaylist(name string, userObj models.SpotifyUserObject) (models.SpotifyPlaylistObject, error)
	AddToSpotifyPlaylist(playlistObj models.SpotifyPlaylistObject, results models.SpotifyTrackURIs) error
}

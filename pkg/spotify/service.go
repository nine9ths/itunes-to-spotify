package spotify

import (
	"github.com/skyerus/itunes-to-spotify/pkg/models"
	"io"
)

type Service interface {
	ReadSongs(r io.Reader) ([]models.Song, error)
	SearchSong(song models.Song) (models.SpotifySearchSimple, error)
}

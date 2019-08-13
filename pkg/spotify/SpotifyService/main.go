package SpotifyService

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/skyerus/itunes-to-spotify/pkg/models"
	"github.com/skyerus/itunes-to-spotify/pkg/spotify"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const SpotifyBaseUrl = "https://api.spotify.com/v1"
const SpotifySearchEndpoint = "/search"
const SpotifyPlaylistsEndpoint = "/playlists"

type spotifyService struct {
	Token string
}

func NewSpotifyService(token string) spotify.Service {
	return &spotifyService{Token: token}
}

func (s spotifyService) ReadSongs(inputFile io.Reader) ([]models.Song, error) {
	var songs []models.Song

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "	")
		if len(s) < 2 {
			return songs, errors.New("file incorrectly formatted")
		}

		var song models.Song
		song.Name = s[0]
		song.Artist = s[1]

		songs = append(songs, song)
	}

	return songs, nil
}

func (s spotifyService) SearchSong(song models.Song) (result models.SpotifySearchSimple, success bool, err error) {
	var simpleSearchResponse models.SpotifySearchSimple
	client := &http.Client{}

	request, err := http.NewRequest("GET", SpotifyBaseUrl + SpotifySearchEndpoint + "?q=" + song.Artist + "%20" + song.Name, nil)
	if err != nil {
		return simpleSearchResponse, false, err
	}

	request.Header.Set("Authorization", "Bearer " + s.Token)

	response, err := client.Do(request)
	if err != nil {
		return simpleSearchResponse, false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return simpleSearchResponse, false, errors.New("access token is unauthorized, it might have expired")
	}

	if response.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return simpleSearchResponse, false, err
		}
		return simpleSearchResponse, false, errors.New(string(bodyBytes))
	}

	var spotifySearch models.SpotifySearch

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return simpleSearchResponse, false, err
	}
	err = json.Unmarshal(body, &spotifySearch)
	if err != nil {
		return simpleSearchResponse, false, err
	}

	if len(spotifySearch.Tracks.Items) < 1 {
		return simpleSearchResponse, false, nil
	}

	simpleSearchResponse.URI = spotifySearch.Tracks.Items[0].URI

	return simpleSearchResponse, true, nil
}

func (s spotifyService) SearchSongs(songs []models.Song) ([]models.SpotifySearchSimple, []models.Song, error) {
	results := make([]models.SpotifySearchSimple, 0, len(songs))
	nonexistentSongs := make([]models.Song, 0, len(songs))

	for i := 0; i < len(songs); i++ {
		result, success, err := s.SearchSong(songs[i])
		if err != nil {
			return results, nonexistentSongs, err
		}
		if !success {
			nonexistentSongs[i] = songs[i]
			continue
		}
		results[i] = result
	}

	return results, nonexistentSongs, nil
}

func (s spotifyService) CreateSpotifyPlaylist(name string) error {
	var playlistPayload models.SpotifyPlaylist
	playlistPayload.Name = name

	client := &http.Client{}

	bodyBytes, err := json.Marshal(playlistPayload)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(bodyBytes)

	request, err := http.NewRequest("POST", SpotifyBaseUrl + SpotifyPlaylistsEndpoint, b)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer " + s.Token)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("access token is unauthorized, it might have expired")
	}

	if response.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bodyBytes))
	}
}

func (s spotifyService) AddResultsToSpotifyPlaylist([]models.SpotifySearchSimple) error {

}

func (s spotifyService) AddNonexistentToFile([]models.Song) error {
	panic("implement me")
}

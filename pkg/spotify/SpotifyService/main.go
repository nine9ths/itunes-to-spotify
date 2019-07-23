package SpotifyService

import (
	"bufio"
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

func (s spotifyService) SearchSong(song models.Song) (models.SpotifySearchSimple, error) {
	var simpleSearchResponse models.SpotifySearchSimple
	client := &http.Client{}

	request, err := http.NewRequest("GET", SpotifyBaseUrl + SpotifySearchEndpoint + "?q=" + song.Artist + "%20" + song.Name, nil)
	if err != nil {
		return simpleSearchResponse, err
	}

	request.Header.Set("Authorization", "Bearer " + s.Token)

	response, err := client.Do(request)
	if err != nil {
		return simpleSearchResponse, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return simpleSearchResponse, errors.New("access token is unauthorized, it might have expired")
	}

	if response.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return simpleSearchResponse, err
		}
		return simpleSearchResponse, errors.New(string(bodyBytes))
	}

	var spotifySearch models.SpotifySearch

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return simpleSearchResponse, err
	}
	err = json.Unmarshal(body, &spotifySearch)
	if err != nil {
		return simpleSearchResponse, err
	}

	if len(spotifySearch.Tracks.Items) < 1 {
		return simpleSearchResponse, errors.New()
	}
}

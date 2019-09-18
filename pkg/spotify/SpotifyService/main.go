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
	"math"
	"net/http"
	"os"
	"strings"
)

const SpotifyBaseUrl = "https://api.spotify.com/v1"
const SpotifySearchEndpoint = "/search"
const SpotifyPlaylistsEndpoint = "/playlists"
const SpotifyTracksEndpoint = "/tracks"

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
		s := strings.Split(scanner.Text(), "\t")
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

func (s spotifyService) SearchSong(song models.Song) (result string, success bool, err error) {
	var simpleSearchResponse string
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

	simpleSearchResponse = spotifySearch.Tracks.Items[0].URI

	return simpleSearchResponse, true, nil
}

func (s spotifyService) SearchSongs(songs []models.Song) ([]string, []models.Song, error) {
	results := make([]string, 0, len(songs))
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

func (s spotifyService) CreateSpotifyPlaylist(name string) (models.SpotifyPlaylistObject, error) {
	var playlistPayload models.SpotifyPlaylist
	var playlistObject models.SpotifyPlaylistObject
	playlistPayload.Name = name

	client := &http.Client{}

	bodyBytes, err := json.Marshal(playlistPayload)
	if err != nil {
		return playlistObject, err
	}
	b := bytes.NewBuffer(bodyBytes)

	request, err := http.NewRequest("POST", SpotifyBaseUrl + SpotifyPlaylistsEndpoint, b)
	if err != nil {
		return playlistObject, err
	}

	request.Header.Set("Authorization", "Bearer " + s.Token)

	response, err := client.Do(request)
	if err != nil {
		return playlistObject, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return playlistObject, errors.New("access token is unauthorized, it might have expired")
	}

	if response.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return playlistObject, err
		}
		return playlistObject, errors.New(string(bodyBytes))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return playlistObject, err
	}
	err = json.Unmarshal(body, &playlistObject)
	if err != nil {
		return playlistObject, err
	}

	return playlistObject, nil
}

func (s spotifyService) AddResultsToSpotifyPlaylist(playlistObj models.SpotifyPlaylistObject, results []string) error {
	var spotifyTrackURIs models.SpotifyTrackURIs
	numOfSongs := len(results)
	if numOfSongs > 100 {
		iterations := int(math.Round(float64(numOfSongs / 100)))
		for i := 0; i < iterations - 1; i++ {
			spotifyTrackURIs.URIs = results[i * 100:(i + 1) * 100]
			err := s.AddToSpotifyPlaylist(playlistObj, spotifyTrackURIs)
			if err != nil {
				return err
			}
		}
		spotifyTrackURIs.URIs = results[(iterations - 1) * 100:]
		return s.AddToSpotifyPlaylist(playlistObj, spotifyTrackURIs)
	}

	spotifyTrackURIs.URIs = results
	return s.AddToSpotifyPlaylist(playlistObj, spotifyTrackURIs)
}

func (s spotifyService) AddToSpotifyPlaylist(playlistObj models.SpotifyPlaylistObject, results models.SpotifyTrackURIs) error {
	client := &http.Client{}

	bodyBytes, err := json.Marshal(results)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(bodyBytes)

	request, err := http.NewRequest("POST", SpotifyBaseUrl + SpotifyPlaylistsEndpoint + "/" + playlistObj.ID + SpotifyTracksEndpoint, b)
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

	return nil
}

func (s spotifyService) AddNonexistentToFile(songs []models.Song, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte("Name\tArtist\tComposer\tAlbum\tGrouping\tWork\tMovement Number\tMovement Count\tMovement Name\tGenre\tSize\tTime\tDisc Number\tDisc Count\tTrack Number\tTrack Count\tYear\tDate Modified\tDate Added\tBit Rate\tSample Rate\tVolume Adjustment\tKind\tEqualiser\tComments\tPlays\tLast Played\tSkips\tLast Skipped\tMy Rating\tLocation\n"))
	if err != nil {
		return err
	}

	for i := 0; i < len(songs); i++ {
		_, err = f.Write([]byte(songs[i].Name + "\t" + songs[i].Artist + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

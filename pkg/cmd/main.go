package cmd

import (
	"flag"
	"fmt"
	"github.com/skyerus/itunes-to-spotify/pkg/spotify/SpotifyService"
	"os"
)

func main() {
	pathToInputFile := flag.String("path", "", "Please provide the path to the input file")
	playlistName := flag.String("playlist", "", "Please provide the name of the playlist to create")
	token := flag.String("token", "", "Please provide a Spotify access token with the required scopes (check the README.md)")
	flag.Parse()

	if *pathToInputFile == "" {
		fmt.Println("Please provide the path to the input file e.g. ./itunes-to-spotify -path=/test.txt")
		return
	}

	if *playlistName == "" {
		fmt.Println("Please provide the name of the playlist to create")
		return
	}

	if *token == "" {
		fmt.Println("Please provide a Spotify access token with the required scopes (check the README.md)")
		return
	}

	if _, err := os.Stat(*pathToInputFile); os.IsNotExist(err) {
		fmt.Println("No file exists at the provided path")
		return
	}

	inputFile, err := os.Open(*pathToInputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	spotifyService := SpotifyService.NewSpotifyService(*token)

	songs, err := spotifyService.ReadSongs(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = spotifyService.CreateSpotifyPlaylist(*playlistName)
	if err != nil {
		fmt.Println(err)
		return
	}
}

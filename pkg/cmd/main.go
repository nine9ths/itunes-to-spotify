package cmd

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	pathToInputFile := flag.String("path", "", "Please provide the path to the input file")
	flag.Parse()

	if *pathToInputFile == "" {
		fmt.Println("Please provide the path to the input file e.g. ./itunes-to-spotify -path=/test.txt")
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
}

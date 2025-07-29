package main

import (
	"log"
	"os"

	"github.com/beeploop/foorrent/torrent"
)

func main() {
	inputFile := os.Args[1]
	outputPath := os.Args[2]

	content, err := torrent.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening torrent file: %s\n", err.Error())
	}

	if err := torrent.Download(content, outputPath); err != nil {
		log.Fatalf("Error downloading torrent file: %s\n", err.Error())
	}

	log.Println("Download complete")
}

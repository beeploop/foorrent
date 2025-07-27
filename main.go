package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beeploop/foorrent/torrent"
)

func main() {
	inputFile := os.Args[1]

	content, err := torrent.Open(inputFile)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(content.Announce)
	fmt.Println(content.Comment)
	t := time.Unix(int64(content.CreationDate), 0)
	fmt.Println(t)
}

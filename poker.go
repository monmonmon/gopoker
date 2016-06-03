package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	numPlayers := flag.Int("n", 0, "Number of players (2 or bigger)")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	if *numPlayers < 2 {
		fmt.Println("Too small number of players")
		flag.PrintDefaults()
		os.Exit(1)
	}

	game := NewGame(*numPlayers)
	for {
		if cont := game.Start(); !cont {
			break
		}
	}
}

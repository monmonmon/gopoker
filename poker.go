package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

type Card struct {
	Suit   string
	Number int
}

type Player struct {
	Num   int
	Cards [2]Card
}

func main() {
	numberOfPlayers := flag.Int("n", 0, "Number of players (2 or bigger)")
	flag.Parse()

	if *numberOfPlayers < 2 {
		fmt.Println("Too small number of players")
		flag.PrintDefaults()
		os.Exit(1)
	}

	players := make([]Player, *numberOfPlayers)
	for i := 0; i < *numberOfPlayers; i++ {
		players[i] = Player{Num: i}
	}

	// initialize cards
	suits := [4]string{"Spade", "Diamond", "Club", "Heart"}
	cardNumbers := [13]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	cards := make([]Card, 13*4+2)

	for i := 0; i < len(suits); i++ {
		for j := 0; j < len(cardNumbers); j++ {
			cards[i*13+j] = Card{Suit: suits[i], Number: cardNumbers[j]}
		}
	}
	cards[52] = Card{Suit: "Joker"}
	cards[53] = Card{Suit: "Joker"}

	// shuffle cards
	shuffledCards := make([]Card, len(cards))
	perm := rand.Perm(len(cards))
	for i, v := range perm {
		shuffledCards[v] = cards[i]
	}

	fmt.Println(shuffledCards)
}

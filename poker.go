package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

type Game struct {
	Players        []Player
	Suits          [4]string
	CardNumbers    [13]int
	Cards          [4*13 + 2]Card
	CommunityCards [5]Card
}

type Player struct {
	Num   int
	Cards [2]Card
}

type Card struct {
	Suit   string
	Number int
}

func NewGame(n int) *Game {
	players := make([]Player, n)
	for i := 0; i < n; i++ {
		players[i] = Player{Num: i}
	}

	return &Game{
		Players:        players,
		Suits:          [4]string{"Spade", "Diamond", "Club", "Heart"},
		CardNumbers:    [13]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		Cards:          [4*13 + 2]Card{},
		CommunityCards: [5]Card{},
	}
}

func (g *Game) Start() (cont bool) {
	g.ShuffleCards()
	fmt.Println(g.Cards)
	return false

	//for {
	//	player = game.NextPlayer()
	//	player.Play()
	//	if game.finished() {
	//		break
	//	}
	//}

	//cont = g.AskContinue()
	//return cont
}

func (g *Game) ShuffleCards() {
	// initialize cards
	cards := [4*13 + 2]Card{}
	for i := 0; i < len(g.Suits); i++ {
		for j := 0; j < len(g.CardNumbers); j++ {
			cards[i*13+j] = Card{Suit: g.Suits[i], Number: g.CardNumbers[j]}
		}
	}
	cards[52] = Card{Suit: "Joker"}
	cards[53] = Card{Suit: "Joker"}
	// shuffle
	perm := rand.Perm(len(cards))
	for i, v := range perm {
		g.Cards[v] = cards[i]
	}
}

func main() {
	numberOfPlayers := flag.Int("n", 0, "Number of players (2 or bigger)")
	flag.Parse()

	if *numberOfPlayers < 2 {
		fmt.Println("Too small number of players")
		flag.PrintDefaults()
		os.Exit(1)
	}

	game := NewGame(*numberOfPlayers)
	for {
		if cont := game.Start(); !cont {
			break
		}
	}
}

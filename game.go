package main

import (
	"fmt"
	"math/rand"
)

type Game struct {
	players             []*Player
	Suits               [4]string
	CardNumbers         [13]uint
	Cards               [4*13 + 2]Card
	CommunityCards      [5]Card
	BetAmount           uint
	Pot                 uint
	cardPointer         uint
	dealerIndex         uint
	playerIndex         uint
	currentBettingRound uint
}

func NewGame(n int) *Game {
	players := make([]*Player, n)
	for i := 0; i < n; i++ {
		players[i] = &Player{Num: uint(i), Chip: 100}
	}
	g := &Game{
		players:        players,
		Suits:          [4]string{"Spade", "Diamond", "Club", "Heart"},
		CardNumbers:    [13]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		Cards:          [4*13 + 2]Card{},
		CommunityCards: [5]Card{},
		BetAmount:      10,
	}
	for i := 0; i < n; i++ {
		players[i].game = g
	}
	return g
}

func (g *Game) ChooseFirstDealer() {
	g.dealerIndex = uint(rand.Intn(len(g.players)))
	g.playerIndex = g.dealerIndex
}

func (g *Game) ChooseNextDealer() {
	g.dealerIndex = (g.dealerIndex + 1) % uint(len(g.players))
}

func (g *Game) CurrentDealer() *Player {
	return g.players[g.dealerIndex]
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
	g.cardPointer = 0
}

func (g *Game) NextCard() (c Card) {
	c = g.Cards[g.cardPointer]
	g.cardPointer++
	return
}

func (g *Game) DealCards() {
	num := len(g.players)
	for i := 0; i < num*2; i++ {
		g.players[i%num].Cards[i/num] = g.NextCard()
	}
}

func (g *Game) BigBlindAmount() uint {
	return g.BetAmount
}

func (g *Game) SmallBlindAmount() uint {
	return g.BetAmount / 2
}

func (g *Game) NextPlayer() *Player {
	g.playerIndex = (g.playerIndex + 1) % uint(len(g.players))
	return g.players[g.playerIndex]
}

func (g *Game) CurrentPlayer() *Player {
	return g.players[g.playerIndex]
}

func (g *Game) BlindBets() {
	p := g.NextPlayer()
	p.Bet(g.SmallBlindAmount())
	p = g.NextPlayer()
	p.Bet(g.BigBlindAmount())
}

func (g *Game) Start() bool {
	g.ChooseFirstDealer()
	g.ShuffleCards()
	g.DealCards()
	g.BlindBets()

	// テスト出力
	dealer := g.CurrentDealer()
	fmt.Println("current dealer:", dealer.Num)
	fmt.Println("cards:", g.Cards)
	for i := 0; i < len(g.players); i++ {
		fmt.Println("player cards:", i, g.players[i].Cards)
	}
	for i := 0; i < len(g.players); i++ {
		fmt.Println("player bet:", i, g.players[i].BetAmount)
	}

	//for {
	//	for {
	//		player = g.NextPlayer()
	//		player.Action()
	//		if g.BettingRoundFinished() {
	//			break
	//		}
	//	}
	//	if !g.NextBettingRound() {
	//		break
	//	}
	//}

	//g.Showdown()

	return false
}

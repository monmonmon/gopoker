package main

import (
	"fmt"
	"math/rand"
)

var Suits = [4]string{"Spade", "Diamond", "Club", "Heart"}
var CardNumbers = [13]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

type Game struct {
	players             []*Player
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
		players[i] = &Player{Num: uint(i), ChipAmount: 100}
	}
	g := &Game{
		players:        players,
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
	for i := 0; i < len(Suits); i++ {
		for j := 0; j < len(CardNumbers); j++ {
			cards[i*13+j] = Card{Suit: Suits[i], Number: CardNumbers[j]}
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

func (g *Game) NextPlayer() (p *Player) {
	for i := 0; i < len(g.players); i++ {
		g.playerIndex = (g.playerIndex + 1) % uint(len(g.players))
		p = g.players[g.playerIndex]
		if !p.folded && p.ChipAmount > 0 {
			return
		}
	}
	return nil
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

func (g *Game) CurrentBetAmountPerPerson() uint {
	max := uint(0)
	for i := 0; i < len(g.players); i++ {
		player := g.players[i]
		if max < player.BetAmount {
			max = player.BetAmount
		}
	}
	return max
}

func (g *Game) Start() bool {
	g.ChooseFirstDealer()
	g.ShuffleCards()
	g.DealCards()
	g.BlindBets()

	// テスト出力
	dealer := g.CurrentDealer()
	fmt.Println("cards:", g.Cards)
	fmt.Println("current dealer:", dealer.Num)
	for i := 0; i < len(g.players); i++ {
		fmt.Println("player cards:", i, g.players[i].Cards)
	}
	for i := 0; i < len(g.players); i++ {
		fmt.Println("player bet:", i, g.players[i].BetAmount)
	}
	fmt.Println("current player:", g.CurrentPlayer().Num)

	player := g.NextPlayer()
	player.Action()

	//for {
	//	// Betting Round
	//	for {
	//		player := g.NextPlayer()
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

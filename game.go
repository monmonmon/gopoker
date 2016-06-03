package main

import (
	"fmt"
	"math/rand"
)

var Suits = [4]string{"Spade", "Diamond", "Club", "Heart"}
var CardNumbers = [13]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

type Game struct {
	Players             Players
	Cards               [4*13 + 2]Card
	CommunityCards      []Card
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
		Players:   players,
		Cards:     [4*13 + 2]Card{},
		BetAmount: 10,
	}
	for i := 0; i < n; i++ {
		players[i].game = g
	}
	return g
}

func (g *Game) ChooseFirstDealer() {
	g.dealerIndex = uint(rand.Intn(len(g.Players)))
	g.playerIndex = g.dealerIndex
}

func (g *Game) ChooseNextDealer() {
	g.dealerIndex = (g.dealerIndex + 1) % uint(len(g.Players))
}

func (g *Game) CurrentDealer() *Player {
	return g.Players[g.dealerIndex]
}

func (g *Game) ShuffleCards() {
	// initialize cards
	cards := [4*13 + 2]Card{}
	for i, suit := range Suits {
		for j, num := range CardNumbers {
			cards[i*13+j] = Card{Suit: suit, Number: num}
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
	num := len(g.Players)
	for i := 0; i < num*2; i++ {
		g.Players[i%num].Cards[i/num] = g.NextCard()
	}
}

func (g *Game) BigBlindAmount() uint {
	return g.BetAmount
}

func (g *Game) SmallBlindAmount() uint {
	return g.BetAmount / 2
}

func (g *Game) NextPlayer() (p *Player) {
	for _, _ = range g.Players {
		g.playerIndex = (g.playerIndex + 1) % uint(len(g.Players))
		p = g.Players[g.playerIndex]
		if !p.folded && p.ChipAmount > 0 {
			return
		}
	}
	return nil
}

func (g *Game) CurrentPlayer() *Player {
	return g.Players[g.playerIndex]
}

func (g *Game) BlindBets() {
	p := g.NextPlayer()
	p.Bet(g.SmallBlindAmount())
	p = g.NextPlayer()
	p.Bet(g.BigBlindAmount())
}

func (g *Game) CurrentBetAmountPerPerson() uint {
	max := uint(0)
	for _, p := range g.Players {
		if max < p.BetAmount {
			max = p.BetAmount
		}
	}
	return max
}

func (g *Game) BettingRoundFinished() (yes bool) {
	amount := g.Players.Playing()[0].BetAmount
	for _, p := range g.Players.Playing() {
		if !p.playedRound || p.BetAmount != amount {
			return false
		}
	}
	return true
}

func (g *Game) NextBettingRound() bool {
	if g.currentBettingRound >= 3 {
		return false
	} else {
		g.currentBettingRound += 1
		g.MoveChipsToPot()
		g.DealCommunityCards()
		g.playerIndex = g.dealerIndex
		for _, p := range g.Players {
			p.playedRound = false
		}
		g.PrintStatus()
		return true
	}
}

func (g *Game) DealCommunityCards() {
	n := 0
	switch g.currentBettingRound {
	case 1:
		n = 3
	case 2, 3:
		n = 1
	default:
		panic("hoe-")
	}
	for i := 0; i < n; i++ {
		g.CommunityCards = append(g.CommunityCards, g.NextCard())
	}
}

func (g *Game) MoveChipsToPot() {
	for _, p := range g.Players {
		g.Pot += p.BetAmount
		p.BetAmount = 0
	}
}

func (g *Game) CurrentBettingRound() string {
	switch g.currentBettingRound {
	case 0:
		return "1st Betting Round"
	case 1:
		return "2nd Betting Round"
	case 2:
		return "3rd Betting Round"
	case 3:
		return "4th Betting Round"
	default:
		panic("hoe-")
	}
}

func (g *Game) PrintStatus() {
	fmt.Println(g.CurrentBettingRound())
	fmt.Println("Community Cards:", g.CommunityCards)
	fmt.Println("Pot:", g.Pot)
	for _, p := range g.Players {
		fmt.Printf("Player %d: %d Bets / %d Chips", p.Num, p.BetAmount, p.ChipAmount)
		if p.folded {
			fmt.Print(" (folded)")
		}
		fmt.Println()
	}
}

func (g *Game) Start() bool {
	g.ChooseFirstDealer()
	g.ShuffleCards()
	g.DealCards()
	g.BlindBets()

	//// テスト出力
	//dealer := g.CurrentDealer()
	//fmt.Println("cards:", g.Cards)
	//fmt.Println("current dealer:", dealer.Num)
	//for _, p := range g.Players.Playing() {
	//	fmt.Println("player cards:", p.Num, p.Cards)
	//}
	//for _, p := range g.Players.Playing() {
	//	fmt.Println("player bet:", p.Num, p.BetAmount)
	//}

	player := g.NextPlayer()
	player.Action()

	for {
		// Betting Round
		for {
			player := g.NextPlayer()
			player.Action()
			if g.BettingRoundFinished() {
				break
			}
		}
		if !g.NextBettingRound() {
			break
		}
	}

	//g.Showdown()

	return false
}

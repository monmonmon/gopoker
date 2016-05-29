package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	CheckOrCall = iota + 1
	Raise
	Fold
)

type Player struct {
	Num         uint
	Cards       [2]Card
	ChipAmount  uint
	BetAmount   uint
	game        *Game
	folded      bool
	playedRound bool
}

type Players []*Player

func (pp Players) Playing() (ret Players) {
	for _, p := range pp {
		if !p.folded && p.ChipAmount > 0 {
			ret = append(ret, p)
		}
	}
	return
}

func (p *Player) AskForNumber(prompt string) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[Player %d] %s: ", p.Num, prompt)
		text, _ := reader.ReadString('\n')
		if num, err := strconv.Atoi(strings.Trim(text, "\n")); err == nil {
			return num
		}
		if len(text) == 0 {
			fmt.Println()
		}
	}
}

func (p *Player) AskForYesOrNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[Player %d] %s (y/N): ", p.Num, prompt)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		switch text {
		case "y", "Y":
			return true
		case "n", "N", "":
			return false
		}
	}
}

func (p *Player) Bet(amount uint) (ok bool) {
	if p.ChipAmount > amount {
		p.ChipAmount -= amount
		p.BetAmount += amount
		return true
	} else {
		return false
	}
}

func (p *Player) CallAmount() uint {
	return p.game.CurrentBetAmountPerPerson() - p.BetAmount
}

func (p *Player) CanCheck() bool {
	return p.CallAmount() == 0
}

func (p *Player) Action() {
	for {
		prompt := "1:Call, 2:Raise, 3:Fold"
		if p.CanCheck() {
			prompt = "1:Check, 2:Raise, 3:Fold"
		}
		num := p.AskForNumber(prompt)
		ok := false
		switch num {
		case CheckOrCall:
			if p.CanCheck() {
				ok = p.Check()
			} else if p.CallAmount() > 0 {
				ok = p.Call()
			} else {
				panic("hoe-")
			}
		case Raise:
			ok = p.Raise()
		case Fold:
			ok = p.Fold()
		default:
			ok = false
		}
		if ok {
			break
		}
	}
	p.playedRound = true
}

func (p *Player) Check() (ok bool) {
	fmt.Println("Check")
	if p.CanCheck() {
		return true
	} else {
		panic("hoe-")
	}
}

func (p *Player) Call() (ok bool) {
	fmt.Println("Call")
	betAmount := p.CallAmount()
	if betAmount > 0 {
		prompt := fmt.Sprintf("Bet %d chips more to call?", betAmount)
		if p.AskForYesOrNo(prompt) {
			p.Bet(betAmount)
			return true
		} else {
			return false
		}
	} else {
		panic("hoe-")
	}
}

func (p *Player) Raise() (ok bool) {
	fmt.Println("Raise")
	amount := p.AskForNumber("Raise amount?")
	if amount <= 0 || p.ChipAmount < uint(amount) {
		return false
	} else {
		if ok := p.Bet(uint(amount)); ok {
			return true
		} else {
			return false
		}
	}
}

func (p *Player) Fold() (ok bool) {
	fmt.Println("Fold")
	p.folded = true
	return true
}

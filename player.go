package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	CheckOrCall = iota + 1
	Raise
	Fold
)

var playerColors []color.Attribute = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgBlack,
}

type Player struct {
	Num         uint
	Cards       [2]Card
	ChipAmount  uint
	BetAmount   uint
	game        *Game
	folded      bool
	playedRound bool
	w           *color.Color
}

func NewPlayer(num int, chip int) *Player {
	c := playerColors[num%len(playerColors)]
	w := color.New(c)
	return &Player{Num: uint(num), ChipAmount: uint(chip), w: w}
}

func (p *Player) AskForNumber(prompt string) (ret int, ok bool) {
	reader := bufio.NewReader(os.Stdin)
	p.w.Printf("[Player %d] %s: ", p.Num, prompt)
	text, _ := reader.ReadString('\n')
	if len(text) == 0 {
		p.w.Println()
	}
	if num, err := strconv.Atoi(strings.Trim(text, "\n")); err == nil {
		return num, true
	} else {
		return 0, false
	}
}

func (p *Player) AskForYesOrNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		p.w.Printf("[Player %d] %s (y/N): ", p.Num, prompt)
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

func (p *Player) MinRaiseAmount() (ret uint) {
	ret = p.game.CurrentBetAmountPerPerson() - p.BetAmount
	if ret < 0 {
		ret = 0
	}
	return
}

func (p *Player) CanCheck() bool {
	return p.CallAmount() == 0
}

func (p *Player) Action() {
	for {
		var prompt string
		if p.CanCheck() {
			prompt = "1:Check, 2:Raise, 3:Fold"
		} else {
			prompt = "1:Call, 2:Raise, 3:Fold"
		}
		num, ok := p.AskForNumber(prompt)
		if !ok {
			continue
		}
		ok = false
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
	p.w.Println("Check")
	if p.CanCheck() {
		return true
	} else {
		panic("hoe-")
	}
}

func (p *Player) Call() (ok bool) {
	p.w.Println("Call")
	betAmount := p.CallAmount()
	if betAmount > 0 {
		return p.Bet(betAmount)
	} else {
		panic("hoe-")
	}
}

func (p *Player) Raise() (ok bool) {
	p.w.Println("Raise")
	minAmount := p.MinRaiseAmount()
	amount, ok := p.AskForNumber("Raise amount?")
	if !ok || uint(amount) <= minAmount || p.ChipAmount < uint(amount) {
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
	p.w.Println("Fold")
	p.folded = true
	return true
}

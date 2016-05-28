package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Player struct {
	Num       uint
	Cards     [2]Card
	Chip      uint
	BetAmount uint
	game      *Game
}

func (p *Player) Bet(amount uint) (ok bool) {
	if p.Chip > amount {
		p.Chip -= amount
		p.BetAmount += amount
		return true
	} else {
		return false
	}
}

func (p *Player) Action() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[Player %d] 1. Check 2. Raise 3. Fold: ")
		text, _ := reader.ReadString('\n')
		num, err := strconv.Atoi(strings.Trim(text, "\n"))
		fmt.Println(">>", num)
		if err != nil {
			continue
		}
		switch num {
		case 1:
			fmt.Println("Check")
		case 2:
			fmt.Println("Raise")
		case 3:
			fmt.Println("Fold")
		default:
			continue
		}
		break
	}
	return true
}

package main

import "fmt"

type Card struct {
	Suit   string
	Number uint
}

func (c *Card) String() string {
	return fmt.Sprintf("[%s %d]", c.Suit, c.Number)
}

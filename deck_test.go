package gojack

import "testing"

func TestNewDeck(t *testing.T) {
	for i := 1; i <= 8; i++ {
		d := NewDeck(i)
		d.Shuffle()
		for d.CardsLeft() > 0 {
			d.Deal()
		}
	}
}

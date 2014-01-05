package gojack

import "testing"

func TestNewGame(t *testing.T) {
    g := NewGame(2)
    g.AddPlayer("Me")
    g.AddPlayer("You")
    g.InitialDeal()
    if g.ScoreDealer() > 21 || g.ScoreHand("You") > 21 {
        t.Errorf("Score cannot be greater than 21 initially")
    }
    g.Deal("Me")
    g.Deal("You")
    g.DecidePlay()
    if g.IsEnd() {
        t.Errorf("The game just started, can't have eneded")
    }
}

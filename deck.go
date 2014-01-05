package gojack

import (
    "math/rand"
    "time"
)

const DECK_TOTAL int = 4 * (1*2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10 + 10 + 10 + 10)

type Deck struct {
    cards     []string
    num_decks int
    total_val int
}

// Given a number of decks to combine, num_deck
//  Return a pointer to the unshuffled deck
func NewDeck(num_decks int) *Deck {
    d := &Deck{make([]string, 0), num_decks, DECK_TOTAL * num_decks}
    d.BuildDeck()
    return d
}

// Initialize the cards of the deck
func (d *Deck) BuildDeck() {
    values := [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
    suits := [4]string{"C", "S", "D", "H"}
    num_decks := d.num_decks
    d.cards = make([]string, 0, 52*num_decks)
    for num_decks > 0 {
        for _, v := range values {
            for _, s := range suits {
                // each card is a length 2 string
                d.cards = append(d.cards, v+s)
            }
        }
        num_decks--
    }
}

// Shuffle the cards of the deck
func (d *Deck) Shuffle() {
    rand.Seed(time.Now().UnixNano())
    // shuffle the deck in place using Fisher-Yates algorithm
    for i := range d.cards {
        j := rand.Intn(i + 1)
        d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
    }
}

// Return the next card from the top of the deck
func (d *Deck) Deal() string {
    if len(d.cards) > 0 {
        c := d.cards[len(d.cards)-1]
        d.cards = d.cards[:len(d.cards)-1]
        d.total_val -= ScoreCard(c)
        return c
    }
    return ""
}

// Remove the first occurrence of card from the deck
func (d *Deck) RemoveCard(card string) {
    for i := len(d.cards) - 1; i >= 0; i-- {
        if d.cards[i] == card {
            d.cards = append(d.cards[:i], d.cards[i+1:]...)
            d.total_val -= ScoreCard(card)
            break
        }
    }
}

// Return the number of cards left
func (d *Deck) CardsLeft() int {
    return len(d.cards)
}

// Return the total value of all remaining cards in the deck
func (d *Deck) TotalVal() int {
    return d.total_val
}

// Return the score of a non-ace card
func ScoreCard(card string) int {
    score := card[0] - 48
    if score > 9 {
        score = 10
    }
    return int(score)
}

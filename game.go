package gojack

import "fmt"

const (
	HIT = iota
	STAND
	SURRENDER
	SPLIT
	BUST
	CONTINUE
	WIN
	LOSS
	END
)

type Game struct {
	Dealer      map[string][]string
	Players     []map[string]interface{}
	Last_Move   string
	deck        *Deck
	order       map[string]int
	status      map[string]int
	turn_number int
}

// Return a pointer to a new game established with num_decks number of decks
func NewGame(num_decks int) *Game {
	g := &Game{make(map[string][]string), make([]map[string]interface{}, 0),
		"dealer", NewDeck(num_decks), make(map[string]int), make(map[string]int), 1}
	g.deck.Shuffle()
	return g
}

// Add the player name to the game
func (g *Game) AddPlayer(name string) {
	g.Players = append(g.Players, make(map[string]interface{}))
	g.Players[len(g.Players)-1]["name"] = name
	g.Players[len(g.Players)-1]["cards"] = make([]string, 0)
	g.order[name] = len(g.Players) - 1
}

// Perform the initial dealing of the game
func (g *Game) InitialDeal() {
	g.status["dealer"] = CONTINUE
	// dealing order: players, dealer, players, dealer
	for _, player := range g.Players {
		g.status[player["name"].(string)] = CONTINUE
		player["cards"] = append(player["cards"].([]string), g.deck.Deal())
	}
	g.Dealer["cards"] = append(g.Dealer["cards"], g.deck.Deal())
	for _, player := range g.Players {
		player["cards"] = append(player["cards"].([]string), g.deck.Deal())
	}
	g.Dealer["cards"] = append(g.Dealer["cards"], g.deck.Deal())
	g.Last_Move = "dealer"
}

// Deal the next card from the deck to player
func (g *Game) Deal(player string) {
	if player == "dealer" {
		g.Dealer["cards"] = append(g.Dealer["cards"], g.deck.Deal())
	} else {
		p := g.Players[g.order[player]]
		p["cards"] = append(p["cards"].([]string), g.deck.Deal())
	}
}

// Return the value of the dealer's hand
func (g *Game) ScoreDealer() int {
	score := 0
	num_aces := 0
	for _, card := range g.Dealer["cards"] {
		if card[0] != 'A' {
			score += ScoreCard(card)
		} else if num_aces == 0 {
			score += 11
		} else {
			num_aces++
			score++
		}
	}
	return score
}

// Return the value of the player's hand
func (g *Game) ScoreHand(player string) int {
	score := 0
	num_aces := 0
	for _, card := range g.Players[g.order[player]]["cards"].([]string) {
		if card[0] != 'A' {
			score += ScoreCard(card)
		} else {
			num_aces++
		}
	}
	// decide what score to make the aces
	if 21-(score+num_aces) >= 10 {
		score += num_aces + 10
	} else {
		score += num_aces
	}
	return score
}

// Return the decision code of play of the current player
//  Decision codes are STAND and HIT
func (g *Game) DecidePlay() int {
	var current string
	if g.Last_Move == "dealer" {
		current = g.Players[0]["name"].(string)
	} else if g.order[g.Last_Move] == len(g.Players)-1 {
		// the dealer always hits
		return HIT
	} else {
		current = g.Players[g.order[g.Last_Move]+1]["name"].(string)
	}
	// we 'cheat' by counting cards
	s := g.ScoreHand(current)
	if s+int(float64(g.deck.TotalVal())/float64(g.deck.CardsLeft())) > 21 {
		return STAND
	}
	return HIT
}

// Return the status code of the player
//  Status codes are LOSS, WIN (non-dealer only), BUST, CONTINUE, and END (dealer only)
func (g *Game) GetStatus(player string) int {
	return g.status[player]
}

// Update the statuses of players in the game
func (g *Game) UpdateStatuses() {
	dealer := g.ScoreDealer()
	if dealer >= 17 && dealer <= 21 {
		g.status["dealer"] = END
	} else if dealer > 21 {
		g.status["dealer"] = BUST
	} else {
		g.status["dealer"] = CONTINUE
	}
	for player := range g.order {
		score := g.ScoreHand(player)
		if score > 21 {
			g.status[player] = BUST
		} else if score > dealer && g.status["dealer"] == END {
			g.status[player] = WIN
		} else if score <= dealer && g.status["dealer"] == CONTINUE {
			g.status[player] = LOSS
		} else if g.status["dealer"] == BUST {
			g.status[player] = WIN
		} else {
			g.status[player] = CONTINUE
		}
	}
}

// Return whether the game has ended
func (g *Game) IsEnd() bool {
	g.UpdateStatuses()
	return g.status["dealer"] == END || g.status["dealer"] == BUST
}

// Return a string representation of status code
func (g *Game) reprStatus(status int) string {
	switch status {
	case BUST:
		return "busts"
	case CONTINUE:
		return "continues play"
	case WIN:
		return "wins"
	case LOSS:
		return "loses"
	case END:
		return "ends the game"
	}
	return "<invalid status>"
}

// Return a string representation of the current state of the game
func (g *Game) String() string {
	var out string
	if g.IsEnd() {
		out += "The game has ended.\n"
	} else {
		out += "The game continues.\n"
	}
	for player := range g.order {
		out += fmt.Sprintf("%s %s with score of %d.\n",
			player, g.reprStatus(g.status[player]), g.ScoreHand(player))
	}
	out += fmt.Sprintf("%s %s with score of %d.\n",
		"dealer", g.reprStatus(g.status["dealer"]), g.ScoreDealer())
	return out
}

// Update the turn number and go to the next player
func (g *Game) IncrementMove() {
	g.turn_number++
	if g.Last_Move == "dealer" {
		g.Last_Move = g.Players[0]["name"].(string)
	} else {
		if g.order[g.Last_Move] == len(g.Players)-1 {
			g.Last_Move = "dealer"
		} else {
			g.Last_Move = g.Players[g.order[g.Last_Move]+1]["name"].(string)
		}
	}
}

// Return the current player in the game
func (g *Game) CurrentPlayer() string {
	if g.Last_Move == "dealer" {
		return g.Players[0]["name"].(string)
	} else {
		if g.order[g.Last_Move] == len(g.Players)-1 {
			return "dealer"
		} else {
			return g.Players[g.order[g.Last_Move]+1]["name"].(string)
		}
	}
	return ""
}

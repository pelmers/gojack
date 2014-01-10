package gojack

import (
	"encoding/json"
	"io/ioutil"
)

// Return a pointer to a game object loaded from JSON-encoded filename
func LoadGame(filename string) *Game {
	i := NewGame(1)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return i
	}
	json.Unmarshal([]byte(b), &i)
	for k, player := range i.Players {
		i.order[player["name"].(string)] = k
		player["cards2"] = make([]string, 0)
		for _, card := range player["cards"].([]interface{}) {
			player["cards2"] = append(player["cards2"].([]string), card.(string))
			i.deck.RemoveCard(card.(string))
		}
		player["cards"] = player["cards2"]
		player["cards2"] = nil
	}
	for _, card := range i.Dealer["cards"] {
		i.deck.RemoveCard(card)
	}
	return i
}

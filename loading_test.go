package gojack

import (
    "fmt"
    "testing"
)

func TestLoadGame(t *testing.T) {
    g := LoadGame("ex_input.json")
    fmt.Println(g)
}

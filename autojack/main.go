package main

import (
    "fmt"
    g "gojack"
    "os"
    "strconv"
)

const WIDTH = 80

func progress_bar(width, percent int) {
    var i int
    width -= 9
    filled := (width * percent) / 100
    fmt.Print("\r[ ")
    for i = 0; i < filled; i++ {
        fmt.Print("#")
    }
    for ; i < width; i++ {
        fmt.Print("-")
    }
    fmt.Printf(" ] %d%%", percent)
}

func main() {
    usage := "autojack -d [number of decks] -p [number of players] -g [number of games]"
    // defaults in case any are not specified in args
    players, decks, games, progress, new_progress := 4, 2, 500, 0, 0
    win_stats := make(map[string]int)
    bust_stats := make(map[string]int)
    playerlist := []string{"dealer"}
    for i := 0; i < len(os.Args); i++ {
        switch os.Args[i] {
        case "-d":
            decks, _ = strconv.Atoi(os.Args[i+1])
        case "-p":
            players, _ = strconv.Atoi(os.Args[i+1])
        case "-g":
            games, _ = strconv.Atoi(os.Args[i+1])
        case "-h":
            fmt.Println(usage)
            return
        }
    }
    for p := players; p > 0; p-- {
        playerlist = append(playerlist, strconv.Itoa(p))
    }
    fmt.Printf("Simulating %d games with %d players and %d decks ...\n",
        games, players, decks)
    for i := 0; i < games; i++ {
        new_progress = int(float64(i) * 100 / float64(games))
        if new_progress > progress {
            progress = new_progress
            progress_bar(WIDTH, progress)
        }
        game := g.NewGame(decks)
        for p := 0; p < players; p++ {
            game.AddPlayer(playerlist[p+1])
        }
        game.InitialDeal()
        // main game loop
        for !(game.IsEnd()) {
            if game.DecidePlay() == g.HIT {
                game.Deal(game.CurrentPlayer())
            }
            game.IncrementMove()
        }
        for _, p := range playerlist {
            if game.GetStatus(p) == g.WIN {
                win_stats[p]++
            } else if game.GetStatus(p) == g.BUST {
                bust_stats[p]++
            }
        }
    }
    fmt.Printf("\nDealer busted %d games\n", bust_stats["dealer"])
    for _, p := range playerlist[1:] {
        fmt.Printf("Player %s won %d games and busted %d games.\n",
            p, win_stats[p], bust_stats[p])
    }
}

package main

import (
	"fmt"
	. "github.com/cconger/loveletter/lib"
	"math/rand"
)

// Building loveletter which is a simple card game.
// Each player has a 1 card hand which is their role.  Each turn a player
// draws one card and can play either of the two in their hand.

// Value, Name, Count
// 1, Soldier, 5
// 2, Clown, 2
// 3, Knight, 2
// 4, Priestess, 2
// 5, Wizard, 2
// 6, General, 1
// 7, Minister, 1
// 8, Princess, 1

func main() {
	game := CreateGame(3)
	fmt.Println(game)
	game.StartGame()
	fmt.Println(game)

	for game.Gamestate == "Playing" {
		player := game.GetActivePlayer()
		cardIdx := rand.Intn(2)
		var target *Player

		for target == nil || target == player || !target.Alive {
			target = game.Players[rand.Intn(len(game.Players))]
		}

		cardVal := rand.Intn(8) + 1

		game.PlayCard(player, player.Hand[cardIdx].Value, target, cardVal)

		fmt.Println(game)
	}

	winner := game.FindWinningPlayer()
	fmt.Println(winner.Name, "won!")
}

package loveletter

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Game struct {
	Players      []*Player
	Deck         Deck
	HiddenCard   *Card
	Gamestate    string
	actingPlayer int
}

func CreateGame(numOfPlayers int) (game Game) {
	for i := 0; i < numOfPlayers; i++ {
		player := CreatePlayer("Player" + strconv.Itoa(i))
		game.Players = append(game.Players, player)
	}

	game.Deck = CreateDeck()
	game.Gamestate = "Created"

	return game
}

func (game *Game) StartGame() {
	rand.Seed(time.Now().Unix())
	fmt.Println(game.Deck.Cards)
	for i := 0; i < 10; i++ {
		game.Deck.shuffle()
	}
	fmt.Println(game.Deck.Cards)

	game.HiddenCard = game.Deck.drawCard()

	for _, player := range game.Players {
		player.DealCard(game.Deck.drawCard())
	}

	game.actingPlayer = rand.Intn(len(game.Players))
	game.Gamestate = "Playing"
	game.DrawCard(game.Players[game.actingPlayer])
}

/* PlayCard is the main method to the game.  On your turn you will draw a new
* card which will be your hand.  You can then play either one.
*  Player is the player who is acting,
*  cardType is the desired card to play.
*  target is the target player for your action (can be null)
*  argument is only used in the case of the soldier in which you're guessing
*  another players card.  argument is your guess.
 */
func (game *Game) PlayCard(player *Player, cardType int, target *Player, argument int) {

	fmt.Println(player.Name, "plays", CardNameForValue(cardType), "on", target.Name)

	if player != game.Players[game.actingPlayer] {
		fmt.Println("ERROR: It is not this players turn...")
		return
	}

	if !player.HasCard(cardType) {
		fmt.Println("ERROR: This player doesn't have that card to play...")
		return
	}

	// Remove the card from hand
	player.ExpendCard(cardType)
	player.SetProtected(false)

	switch cardType {
	case 1:
		if target.IsProtected() || !target.Alive {
			fmt.Println(target.Name, "is an invalid target")
			break
		}
		if target.HasCard(argument) && !target.IsProtected() {
			target.KillPlayer()
		}
	case 2:
		fmt.Println(target.Name, "has card", target.Hand[0])
	case 3:
		if target.IsProtected() || !target.Alive {
			fmt.Println(target.Name, "is an invalid target")
			break
		}
		if target.Hand[0].Value < player.Hand[0].Value {
			fmt.Println(target.Name, "had", target.Hand[0], "which is beaten by your", player.Hand[0])
			target.KillPlayer()
		} else if target.Hand[0].Value > player.Hand[0].Value {
			fmt.Println(player.Name, "had", player.Hand[0], "which is beaten by the", target.Hand[0], "of", target.Name)
			player.KillPlayer()
		} else {
			fmt.Println(player.Name, "tied with the", player.Hand[0], "of", target.Name)
		}
	case 4:
		player.SetProtected(true)
	case 5:
		if target.IsProtected() || !target.Alive {
			fmt.Println(target.Name, "is an invalid target")
			break
		}
		game.DiscardCard(target)
		game.DrawCard(target)
	case 6:
		if target.IsProtected() || !target.Alive {
			fmt.Println(target.Name, "is an invalid target")
			break
		}
		target.Hand[0], player.Hand[0] = player.Hand[0], target.Hand[0]
	case 7:
		//Nothing Happens
	case 8:
		player.KillPlayer()
	}

	if game.IsGameOver() {
		game.Gamestate = "GameOver"
		fmt.Println("Game Over!")
		return
	}

	game.AdvancePlayer()
	game.DrawCard(nil)
}

func (game *Game) DiscardCard(player *Player) {
	if player.HasCard(8) {
		player.KillPlayer()
	} else {
		player.DiscardCard()
	}
}

func (game *Game) IsGameOver() bool {
	alivePlayers := 0
	for _, player := range game.Players {
		if player.Alive {
			alivePlayers++
		}
	}
	if alivePlayers <= 1 {
		return true
	}

	return false
}

func (game *Game) FindWinningPlayer() *Player {
	var highestRankingPlayer *Player
	highestRank := 0
	playersAtRank := 0
	for _, player := range game.Players {
		if player.Alive {
			playerRank := player.GetRank()
			if playerRank > highestRank {
				highestRank = playerRank
				highestRankingPlayer = player
				playersAtRank = 1
			} else if playerRank == highestRank {
				playersAtRank++
			}
		}
	}

	if playersAtRank > 1 {
		//Go into a tiebreaker based on expended cards
	}

	return highestRankingPlayer
}

func (game *Game) AdvancePlayer() {
	game.actingPlayer = (game.actingPlayer + 1) % len(game.Players)
	player := game.Players[game.actingPlayer]
	if !player.Alive {
		game.AdvancePlayer()
	}
}

func (game *Game) GetActivePlayer() *Player {
	return game.Players[game.actingPlayer]
}

func (game *Game) DrawCard(player *Player) {
	if player == nil {
		player = game.Players[game.actingPlayer]
	}
	nextCard := game.Deck.drawCard()
	if nextCard == nil {
		fmt.Println("Ran out of cards!")
		winner := game.FindWinningPlayer()
		fmt.Println(winner.Name, "won by having the highest rank")
		game.Gamestate = "GameOver"
		return
	}
	player.DrawCard(nextCard)

	if player.HasCard(7) && player.CardSum() >= 12 {
		fmt.Println(player.Name, "was betrayed by their minister")
		player.KillPlayer()
		game.AdvancePlayer()
		game.DrawCard(nil)
	}
}

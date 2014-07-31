package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

type Card struct {
	Value int
}

func CardNameForValue(value int) string {
	cardValues := map[int]string{
		1: "Solider",
		2: "Clown",
		3: "Knight",
		4: "Priestess",
		5: "Wizard",
		6: "General",
		7: "Minister",
		8: "Princess",
	}

	return cardValues[value]
}

func (card Card) String() string {
	return CardNameForValue(card.Value)
}

type Player struct {
	Name          string
	Hand          []*Card
	ExpendedCards []*Card
	Alive         bool
	Protected     bool
}

func CreatePlayer(name string) *Player {
	return &Player{name, make([]*Card, 2), make([]*Card, 10), true, false}
}

func (player *Player) DealCard(card *Card) {
	player.Hand[0] = card
}


func (player *Player) GetRank() int {
  rank := 1
  for _, card := range player.Hand {
    if card != nil && card.Value > rank {
      rank = card.Value
    }
  }
  return rank
}

func (player *Player) DrawCard(card *Card) {
  for idx, val := range player.Hand {
    if val == nil {
      player.Hand[idx] = card
      return
    }
  }
}

func (player *Player) ExpendCard(cardType int) {
	for idx, card := range player.Hand {
		if card != nil && card.Value == cardType {
			player.discardCardInIndex(idx)
		}
	}
}

func (player *Player) DiscardCard() {
  player.discardCardInIndex(0)
}

func (player *Player) discardCardInIndex(idx int) {
	player.ExpendedCards = append(player.ExpendedCards, player.Hand[idx])
	player.Hand[idx] = nil

	if idx == 0 {
		player.Hand[0] = player.Hand[1]
		player.Hand[1] = nil
	}
}

func (player Player) String() string {
  if (!player.Alive) {
    return fmt.Sprintf("%s: <DEAD>", player.Name)
  }
	return fmt.Sprintf("%s: %v", player.Name, player.Hand)
}

func (player *Player) IsProtected() bool {
	return player.Protected
}

func (player *Player) SetProtected(protected bool) {
	player.Protected = protected
}

func (player *Player) KillPlayer() {
	player.Alive = false
	player.ExpendedCards = append(player.ExpendedCards, player.Hand...)
	player.Hand[0] = nil
	player.Hand[1] = nil
	fmt.Println(player.Name, "has been killed")
}

func (player *Player) HasCard(cardType int) bool {
	for _, card := range player.Hand {
		if card != nil && card.Value == cardType {
			return true
		}
	}
	return false
}

func (player *Player) CardSum() int {
	sum := 0
	for _, card := range player.Hand {
		if card != nil {
			sum += card.Value
		}
	}
	return sum
}

type Deck struct {
	Cards []*Card
}

func CreateDeck() (deck Deck) {
	deck.Cards = make([]*Card, 16)

	cardTypes := map[int]int{
		1: 5,
		2: 2,
		3: 2,
		4: 2,
		5: 2,
		6: 1,
		7: 1,
		8: 1,
	}

	idx := 0
	for val, count := range cardTypes {
		for i := 0; i < count; i++ {
			deck.Cards[idx] = &Card{val}
			idx++
		}
	}

	return deck
}

func (deck Deck) String() string {
	return strconv.Itoa(len(deck.Cards))
}

func (deck *Deck) addCard(card *Card) {
	deck.Cards = append(deck.Cards, card)
}

func (deck *Deck) drawCard() *Card {
  if len(deck.Cards) > 0 {
    topCard := deck.Cards[len(deck.Cards)-1]
    deck.Cards = deck.Cards[:len(deck.Cards)-1]

    return topCard
  }
  return nil
}

func (deck *Deck) shuffle() {
	for i := range deck.Cards {
		j := rand.Intn(i + 1)
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	}
}

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

/* playCard is the main method to the game.  On your turn you will draw a new
* card which will be your hand.  You can then play either one.
*  Player is the player who is acting,
*  cardType is the desired card to play.
*  target is the target player for your action (can be null)
*  argument is only used in the case of the soldier in which you're guessing
*  another players card.  argument is your guess.
 */
func (game *Game) playCard(player *Player, cardType int, target *Player, argument int) {

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

func (game *Game) TurnEnd() {
  //Check win conditions here
  game.AdvancePlayer()
  game.DrawCard(nil)
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

func main() {
	game := CreateGame(3)
	fmt.Println(game)
	game.StartGame()
	fmt.Println(game)

	for game.Gamestate == "Playing" {
		player := game.Players[game.actingPlayer]
		cardIdx := rand.Intn(2)
		target := game.Players[(game.actingPlayer+1)%len(game.Players)]
		cardVal := rand.Intn(8) + 1

		game.playCard(player, player.Hand[cardIdx].Value, target, cardVal)

		fmt.Println(game)
	}
}

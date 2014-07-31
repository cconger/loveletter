package loveletter

import (
	"fmt"
)

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
	if !player.Alive {
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

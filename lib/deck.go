package loveletter

import (
	"math/rand"
	"strconv"
)

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

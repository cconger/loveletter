package loveletter

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

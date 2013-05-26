package poker

import (
	"math/rand"
	"sort"
	"time"
)

type PlayerOutcome struct {
	Player int
	Level  HandLevel
	Cards  []Card
}

type HandSorter struct {
	Outcomes []PlayerOutcome
}

func (hs HandSorter) Len() int {
	return len(hs.Outcomes)
}

func (hs HandSorter) Swap(i, j int) {
	hs.Outcomes[i], hs.Outcomes[j] = hs.Outcomes[j], hs.Outcomes[i]
}

func (hs HandSorter) Less(i, j int) bool {
	return Beats(hs.Outcomes[i].Level, hs.Outcomes[j].Level) && !Beats(hs.Outcomes[j].Level, hs.Outcomes[i].Level)
}

type Pack struct {
	Cards   [52]Card
	randGen *rand.Rand
}

func (p *Pack) initialise() {
	p.randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

	i := 0
	for s := 0; s < 4; s++ {
		for r := 0; r < 13; r++ {
			p.Cards[i] = Card{Rank(r), Suit(s)}
			i++
		}
	}
}

// Shuffle the pack
func (p *Pack) Shuffle() {
	for i := 0; i < 52; i++ {
		j := p.randGen.Intn(52-i) + i
		p.Cards[i], p.Cards[j] = p.Cards[j], p.Cards[i]
	}
}

func (p *Pack) PlayHoldem(players int) (onTable []Card, playerCards [][]Card, handSorter HandSorter) {
	onTable = p.Cards[0:5]

	playerCards = make([][]Card, players)
	outcomes := make([]PlayerOutcome, players)
	for player := 0; player < players; player++ {
		playerCards[player] = p.Cards[5+2*player : 7+2*player]
		level, cards := Classify(playerCards[player], onTable)
		outcomes[player] = PlayerOutcome{player + 1, level, cards}
	}

	handSorter = HandSorter{outcomes}
	sort.Sort(handSorter)

	return onTable, playerCards, handSorter
}

func NewPack() Pack {
	var result Pack
	result.initialise()
	return result
}

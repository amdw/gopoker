package poker

import (
	"math/rand"
	"sort"
	"time"
)

type HandSorter struct {
	Hands   []HandLevel
	Players []int
}

func (hs HandSorter) Len() int {
	return len(hs.Hands)
}

func (hs HandSorter) Swap(i, j int) {
	hs.Hands[i], hs.Hands[j] = hs.Hands[j], hs.Hands[i]
	hs.Players[i], hs.Players[j] = hs.Players[j], hs.Players[i]
}

func (hs HandSorter) Less(i, j int) bool {
	return Beats(hs.Hands[i], hs.Hands[j]) && !Beats(hs.Hands[j], hs.Hands[i])
}

func NewHandSorter(levels []HandLevel) HandSorter {
	players := make([]int, len(levels))
	for i := range players {
		players[i] = i + 1
	}
	return HandSorter{levels, players}
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
	hands := make([]HandLevel, players)
	for player := 0; player < players; player++ {
		playerCards[player] = p.Cards[5+2*player : 7+2*player]
		hands[player] = Classify(playerCards[player], onTable)
	}

	handSorter = NewHandSorter(hands)
	sort.Sort(handSorter)

	return onTable, playerCards, handSorter
}

func NewPack() Pack {
	var result Pack
	result.initialise()
	return result
}

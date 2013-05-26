/*
Copyright 2013 Andrew Medworth

This file is part of Gopoker, a set of miscellaneous poker-related functions
written in the Go programming language (http://golang.org).

Gopoker is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Gopoker is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with Gopoker.  If not, see <http://www.gnu.org/licenses/>.
*/
package poker

import (
	"fmt"
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

// Shuffle the pack, but fix certain cards in place. For use in simulations.
func (p *Pack) shuffleFixing(tableCards, yourCards []Card) {
	if len(tableCards) > 5 {
		panic(fmt.Sprintf("Maximum of 5 table cards supported, found %v", len(tableCards)))
	}
	// Remove fixed cards from the pack, shuffle those, and fill them in
	nonFixedCards := make([]Card, 52-len(yourCards)-len(tableCards))

	indexOf := func(cards []Card, card Card) int {
		result := -1
		for i, c := range cards {
			if c == card {
				result = i
				break
			}
		}
		return result
	}

	i := 0
	for _, c := range p.Cards {
		if indexOf(yourCards, c) == -1 && indexOf(tableCards, c) == -1 {
			nonFixedCards[i] = c
			i++
		}
	}
	for i = 0; i < len(nonFixedCards); i++ {
		j := p.randGen.Intn(len(nonFixedCards)-i) + i
		nonFixedCards[i], nonFixedCards[j] = nonFixedCards[j], nonFixedCards[i]
	}

	copy(p.Cards[0:len(tableCards)], tableCards)

	// If not all table cards were supplied, fill in the gaps from the non-fixed cards
	for i = 0; i < 5-len(tableCards); i++ {
		p.Cards[i+len(tableCards)] = nonFixedCards[i]
	}
	copy(p.Cards[5:5+len(yourCards)], yourCards)
	copy(p.Cards[5+len(yourCards):52], nonFixedCards[i:len(nonFixedCards)])
}

func (p *Pack) getHoldemOutcomes(players int) (onTable []Card, playerCards [][]Card, outcomes []PlayerOutcome) {
	onTable = p.Cards[0:5]

	playerCards = make([][]Card, players)
	outcomes = make([]PlayerOutcome, players)
	for player := 0; player < players; player++ {
		playerCards[player] = p.Cards[5+2*player : 7+2*player]
		level, cards := Classify(playerCards[player], onTable)
		outcomes[player] = PlayerOutcome{player + 1, level, cards}
	}
	return onTable, playerCards, outcomes
}

// Play out one hand of Texas Hold'em and return the full outcome (positions of all cards, and outcome for all players).
func (p *Pack) PlayHoldem(players int) (onTable []Card, playerCards [][]Card, handSorter HandSorter) {
	onTable, playerCards, outcomes := p.getHoldemOutcomes(players)

	handSorter = HandSorter{outcomes}
	sort.Sort(handSorter)

	return onTable, playerCards, handSorter
}

// Play out one hand of Texas Hold'em and return whether or not player 1 won, plus player 1's hand level, plus one other player's hand level.
func (p *Pack) SimulateOneHoldemHand(players int) (won bool, ourLevel, opponentLevel HandLevel) {
	_, _, outcomes := p.getHoldemOutcomes(players)

	won = true
	for i := 1; i < len(outcomes); i++ {
		if Beats(outcomes[i].Level, outcomes[0].Level) {
			won = false
			break
		}
	}

	return won, outcomes[0].Level, outcomes[1].Level
}

func NewPack() Pack {
	var result Pack
	result.initialise()
	return result
}

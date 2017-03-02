/*
Copyright 2013, 2015, 2017 Andrew Medworth

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
package holdem

import (
	"errors"
	"fmt"
	"github.com/amdw/gopoker/poker"
	"math/rand"
	"time"
)

func SimulateHoldem(tableCards, yourCards []poker.Card, players, handsToPlay int) *poker.Simulator {
	s := poker.Simulator{}
	// Very crude attempt to detect situation where exhaustive enumeration is cheaper than simulation
	if len(tableCards) == 5 && len(yourCards) == 2 && players == 2 && handsToPlay > 990 {
		enumerateHoldem(&s, tableCards, yourCards, players)
		return &s
	}
	s.Reset(players, handsToPlay)
	p := poker.NewPack()
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < handsToPlay; i++ {
		shuffleFixing(&p, tableCards, yourCards, randGen)
		handOutcome := SimulateOneHoldemHand(&p, players, randGen)
		s.ProcessHand(handOutcome)
	}
	return &s
}

func enumerateHoldem(s *poker.Simulator, tableCards, yourCards []poker.Card, players int) {
	s.Reset(players, 0)
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	// For now we only enumerate the case where we have only one opponent and a full set of table cards.
	remainingPack := make([]poker.Card, 45)
	i := 0
	isUsed := func(card poker.Card) bool {
		for _, c := range tableCards {
			if c == card {
				return true
			}
		}
		for _, c := range yourCards {
			if c == card {
				return true
			}
		}
		return false
	}
	for _, card := range poker.NewPack().Cards {
		if !isUsed(card) {
			remainingPack[i] = card
			i++
		}
	}
	opponentHands := poker.AllCardCombinations(remainingPack, 2)
	for _, opponentHand := range opponentHands {
		playerCards := [][]poker.Card{yourCards, opponentHand}
		outcomes := DealOutcomes(tableCards, playerCards)
		handOutcome := calcHandOutcome(outcomes, randGen)
		s.ProcessHand(handOutcome)
		s.HandCount++
	}
}

// Shuffle the pack, but fix certain cards in place. For use in simulations.
// It is assumed that there are no duplicate cards in (tableCards+yourCards).
func shuffleFixing(p *poker.Pack, tableCards, yourCards []poker.Card, randGen *rand.Rand) {
	if len(tableCards) > 5 || len(yourCards) > 2 {
		panic(fmt.Sprintf("Maximum of 5 table cards and 2 hole cards supported, found %v and %v", len(tableCards), len(yourCards)))
	}

	indexOf := func(cards [52]poker.Card, card poker.Card) int {
		result := -1
		for i, c := range cards {
			if c == card {
				result = i
				break
			}
		}
		return result
	}

	// Just shuffle the pack and then swap the fixed cards into place from wherever they are in the deck
	p.Shuffle(randGen)
	for i := 0; i < len(tableCards); i++ {
		swapIdx := indexOf(p.Cards, tableCards[i])
		p.Cards[i], p.Cards[swapIdx] = p.Cards[swapIdx], p.Cards[i]
	}
	for i := 0; i < len(yourCards); i++ {
		swapIdx := indexOf(p.Cards, yourCards[i])
		targetIdx := i + 5
		p.Cards[targetIdx], p.Cards[swapIdx] = p.Cards[swapIdx], p.Cards[targetIdx]
	}
}

type StartingPair struct {
	Rank1, Rank2 poker.Rank
	SameSuit     bool
}

func (pair StartingPair) Validate() error {
	if pair.Rank1 == pair.Rank2 && pair.SameSuit {
		return errors.New(fmt.Sprintf("Pair of %vs cannot be the same suit!", pair.Rank1))
	}
	return nil
}

func (pair StartingPair) SampleCards() (poker.Card, poker.Card) {
	err := pair.Validate()
	if err != nil {
		panic(err)
	}
	// Just pick arbitrary suits, either the same or different
	card1 := poker.Card{pair.Rank1, poker.Club}
	card2 := poker.Card{pair.Rank2, poker.Heart}
	if pair.SameSuit {
		card2.Suit = poker.Club
	}
	return card1, card2
}

func (pair StartingPair) RunSimulation(players, handsToPlay int) *poker.Simulator {
	card1, card2 := pair.SampleCards()
	return SimulateHoldem([]poker.Card{}, []poker.Card{card1, card2}, players, handsToPlay)
}

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
	"fmt"
	"github.com/amdw/gopoker/poker"
)

func classify(tableCards, holeCards []poker.Card) (poker.HandLevel, []poker.Card) {
	allCards := make([]poker.Card, 7)
	copy(allCards, holeCards)
	copy(allCards[2:], tableCards)

	// Construct all possible hands and find the best one
	allPossibleHands := poker.AllCardCombinations(allCards, 5)

	bestHand := allPossibleHands[0]
	bestRank := poker.ClassifyHand(bestHand)

	for i := 1; i < len(allPossibleHands); i++ {
		hand := allPossibleHands[i]
		rank := poker.ClassifyHand(hand)
		if poker.Beats(rank, bestRank) {
			bestHand = hand
			bestRank = rank
		}
	}

	return bestRank, bestHand
}

type PlayerOutcome struct {
	Player         int
	Level          poker.HandLevel
	Cards          []poker.Card
	Won            bool
	PotFractionWon float64
}

func Deal(p *poker.Pack, players int) (onTable []poker.Card, playerCards [][]poker.Card) {
	if players < 1 {
		panic(fmt.Sprintf("At least one player required, found %v", players))
	}

	onTable = p.Cards[0:5]

	playerCards = make([][]poker.Card, players)
	for player := 0; player < players; player++ {
		playerCards[player] = p.Cards[5+2*player : 7+2*player]
	}

	return onTable, playerCards
}

// Assess the hand each player holds and return a sorted list of outcomes
// by hand strength (descending) then player number (ascending).
// Player numbers are in ascending order of playerCards entries, starting with 1.
func DealOutcomes(onTable []poker.Card, playerCards [][]poker.Card) []PlayerOutcome {
	outcomes := make([]PlayerOutcome, len(playerCards))
	var bestLevel poker.HandLevel
	for playerIdx, hand := range playerCards {
		level, cards := classify(onTable, hand)
		outcomes[playerIdx] = PlayerOutcome{playerIdx + 1, level, cards, false, 0}
		if playerIdx == 0 || poker.Beats(level, bestLevel) {
			bestLevel = level
		}
	}

	winners := make([]int, 0, len(outcomes))
	for i, outcome := range outcomes {
		if !poker.Beats(bestLevel, outcome.Level) {
			winners = append(winners, i)
		}
	}
	for _, i := range winners {
		outcomes[i].Won = true
		outcomes[i].PotFractionWon = 1.0 / float64(len(winners))
	}
	return outcomes
}

/*
Copyright 2017 Andrew Medworth

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
package omaha8

import (
	"github.com/amdw/gopoker/poker"
)

func Deal(pack *poker.Pack, players int) (tableCards []poker.Card, playerCards [][]poker.Card) {
	tableCards = pack.Cards[0:5]
	playerCards = make([][]poker.Card, players)
	for i := 0; i < players; i++ {
		playerCards[i] = pack.Cards[5+(i*4) : 9+(i*4)]
	}
	return tableCards, playerCards
}

func makePossibleCombinations(tableCards, holeCards []poker.Card) [][]poker.Card {
	// We need exactly two cards from the four hole cards and exactly three from the table
	possibleTableCards := poker.AllCardCombinations(tableCards, 3)
	possibleHoleCards := poker.AllCardCombinations(holeCards, 2)

	result := make([][]poker.Card, 0, len(possibleTableCards)*len(possibleHoleCards))

	for _, table := range possibleTableCards {
		for _, hole := range possibleHoleCards {
			combination := make([]poker.Card, 5)
			copy(combination, table)
			copy(combination[3:], hole)
			result = append(result, combination)
		}
	}

	return result
}

// Only low hands which are 8-high or better qualify for consideration as the best low hand
func lowLevelQualifies(level poker.HandLevel) bool {
	return level.Class == poker.HighCard && poker.IsRankLess(level.Tiebreaks[0], poker.Nine, true)
}

type Omaha8Level struct {
	HighLevel, LowLevel poker.HandLevel
	HighHand, LowHand   []poker.Card
	LowLevelQualifies   bool
}

// Identify an Omaha/8 hand
func classify(tableCards, holeCards []poker.Card) Omaha8Level {
	possibleCombinations := makePossibleCombinations(tableCards, holeCards)

	bestHighHand := possibleCombinations[0]
	bestHighLevel := poker.ClassifyHand(bestHighHand)
	bestLowHand := bestHighHand
	bestLowLevel := poker.ClassifyAceToFiveLow(bestLowHand)

	for i := 1; i < len(possibleCombinations); i++ {
		highLevel := poker.ClassifyHand(possibleCombinations[i])
		if poker.Beats(highLevel, bestHighLevel) {
			bestHighHand = possibleCombinations[i]
			bestHighLevel = highLevel
		}

		lowLevel := poker.ClassifyAceToFiveLow(possibleCombinations[i])
		if poker.BeatsAceToFiveLow(lowLevel, bestLowLevel) {
			bestLowHand = possibleCombinations[i]
			bestLowLevel = lowLevel
		}
	}

	return Omaha8Level{bestHighLevel, bestLowLevel, bestHighHand, bestLowHand, lowLevelQualifies(bestLowLevel)}
}

type PlayerOutcome struct {
	Player         int
	Level          Omaha8Level
	PotFractionWon float64
}

func PlayerOutcomes(tableCards []poker.Card, playerCards [][]poker.Card) []PlayerOutcome {
	levels := make([]Omaha8Level, len(playerCards))
	for i, hand := range playerCards {
		levels[i] = classify(tableCards, hand)
	}

	// Now work out how the pot should be divided
	bestHighLevel := levels[0].HighLevel
	bestLowLevel := levels[0].LowLevel
	for i := 1; i < len(levels); i++ {
		if poker.Beats(levels[i].HighLevel, bestHighLevel) {
			bestHighLevel = levels[i].HighLevel
		}
		if poker.BeatsAceToFiveLow(levels[i].LowLevel, bestLowLevel) {
			bestLowLevel = levels[i].LowLevel
		}
	}

	highWinners := make([]int, 0, len(levels))
	lowWinners := make([]int, 0, len(levels))
	lowQualified := lowLevelQualifies(bestLowLevel)
	result := make([]PlayerOutcome, len(levels))
	for i, level := range levels {
		outcome := PlayerOutcome{i + 1, level, 0}
		if !poker.Beats(bestHighLevel, level.HighLevel) {
			highWinners = append(highWinners, i)
		}
		if lowQualified && !poker.BeatsAceToFiveLow(bestLowLevel, level.LowLevel) {
			lowWinners = append(lowWinners, i)
		}
		result[i] = outcome
	}

	// If there were no eligible low hands then the whole pot goes to the best high hand
	highMultiple := 0.5
	if len(lowWinners) == 0 {
		highMultiple = 1.0
	}
	for _, highWinnerIdx := range highWinners {
		result[highWinnerIdx].PotFractionWon += highMultiple / float64(len(highWinners))
	}
	for _, lowWinnerIdx := range lowWinners {
		result[lowWinnerIdx].PotFractionWon += 0.5 / float64(len(lowWinners))
	}

	return result
}

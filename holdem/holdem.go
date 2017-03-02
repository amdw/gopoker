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
	"math/rand"
	"sort"
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
	Player int
	Level  poker.HandLevel
	Cards  []poker.Card
}

func sortOutcomes(outcomes []PlayerOutcome) {
	sort.Slice(outcomes, func(i, j int) bool {
		iBeatsJ := poker.Beats(outcomes[i].Level, outcomes[j].Level)
		jBeatsI := poker.Beats(outcomes[j].Level, outcomes[i].Level)
		if iBeatsJ && !jBeatsI {
			return true
		}
		if jBeatsI && !iBeatsJ {
			return false
		}
		return outcomes[i].Player < outcomes[j].Player
	})
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
	for playerIdx, hand := range playerCards {
		level, cards := classify(onTable, hand)
		outcomes[playerIdx] = PlayerOutcome{playerIdx + 1, level, cards}
	}
	sortOutcomes(outcomes)
	return outcomes
}

func calcPotFraction(won bool, potSplit int) float64 {
	if won {
		return 1.0 / float64(potSplit)
	} else {
		return 0
	}
}

func calcHandOutcome(outcomes []PlayerOutcome, randGen *rand.Rand) *poker.HandOutcome {
	// This works even if player 1 was equal first, since equal hands are sorted by player
	won := outcomes[0].Player == 1

	var ourOutcome PlayerOutcome
	opponentOutcomes := make([]PlayerOutcome, len(outcomes)-1)
	potSplit := 0
	i := 0
	for _, o := range outcomes {
		if o.Player == 1 {
			ourOutcome = o
		} else {
			opponentOutcomes[i] = o
			i++
		}
		if !poker.Beats(outcomes[0].Level, o.Level) {
			// The best hand doesn't beat this hand so it must be a winner
			potSplit++
		}
	}

	bestOpponentWon := !poker.Beats(ourOutcome.Level, opponentOutcomes[0].Level)

	randomOpponentIdx := randGen.Intn(len(opponentOutcomes))
	randomOpponentLevel := opponentOutcomes[randomOpponentIdx].Level
	randomOpponentWon := !poker.Beats(outcomes[0].Level, randomOpponentLevel)

	potFractionWon := calcPotFraction(won, potSplit)
	bestOpponentPotFraction := calcPotFraction(bestOpponentWon, potSplit)
	randomOpponentPotFraction := calcPotFraction(randomOpponentWon, potSplit)

	return &poker.HandOutcome{won, bestOpponentWon, randomOpponentWon,
		potFractionWon, bestOpponentPotFraction, randomOpponentPotFraction,
		ourOutcome.Level, opponentOutcomes[0].Level, randomOpponentLevel}
}

// Play out one hand of Texas Hold'em and return whether or not player 1 won,
// plus player 1's hand level, plus the best hand level of any of player 1's opponents.
func SimulateOneHoldemHand(p *poker.Pack, players int, randGen *rand.Rand) *poker.HandOutcome {
	onTable, playerCards := Deal(p, players)
	outcomes := DealOutcomes(onTable, playerCards)
	return calcHandOutcome(outcomes, randGen)
}

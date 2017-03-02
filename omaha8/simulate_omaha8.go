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
	"math/rand"
)

type Omaha8Simulator struct {
	HighSimulator poker.Simulator
	LowSimulator  Omaha8LowSimulator
}

func (s *Omaha8Simulator) reset(players, handsToPlay int) {
	s.HighSimulator.Reset(players, handsToPlay)
	s.LowSimulator.reset(handsToPlay)
}

func (s *Omaha8Simulator) processHand(playerOutcomes []PlayerOutcome, randGen *rand.Rand) {
	randomOpponentIdx := 1 + randGen.Intn(len(playerOutcomes)-1)
	highOutcome := calcHighOutcome(playerOutcomes, randomOpponentIdx)
	lowOutcome := calcLowOutcome(playerOutcomes, randomOpponentIdx)
	s.HighSimulator.ProcessHand(highOutcome)
	s.LowSimulator.processHand(lowOutcome)
}

func (s *Omaha8Simulator) PotsWon() float64 {
	return s.HighSimulator.PotsWon + s.LowSimulator.PotsWon
}

func (s *Omaha8Simulator) PotOddsBreakEven() float64 {
	return poker.PotOddsBreakEven(s.PotsWon(), s.HighSimulator.HandCount)
}

// Over time we can expand this to track the low-hand class counts etc
type Omaha8LowSimulator struct {
	HandCount int
	WinCount  int
	PotsWon   float64
}

func (s *Omaha8LowSimulator) reset(handsToPlay int) {
	s.HandCount = handsToPlay
	s.WinCount = 0
	s.PotsWon = 0
}

func (s *Omaha8LowSimulator) processHand(outcome *poker.HandOutcome) {
	if outcome.Won {
		s.WinCount++
	}
	s.PotsWon += outcome.PotFractionWon
}

func SimulateOmaha8(tableCards, yourCards []poker.Card, players, handsToPlay int, randGen *rand.Rand) *Omaha8Simulator {
	sim := Omaha8Simulator{}
	sim.reset(players, handsToPlay)

	p := poker.NewPack()
	for i := 0; i < handsToPlay; i++ {
		shuffleFixing(&p, tableCards, yourCards, randGen)
		tableCards, playerCards := Deal(&p, players)
		playerOutcomes := PlayerOutcomes(tableCards, playerCards)
		sim.processHand(playerOutcomes, randGen)
	}

	return &sim
}

func shuffleFixing(pack *poker.Pack, tableCards, yourCards []poker.Card, randGen *rand.Rand) {
	// Do a regular shuffle and swap the target cards into place
	pack.Shuffle(randGen)
	for i, c := range tableCards {
		swapIdx := pack.IndexOf(c)
		pack.Cards[i], pack.Cards[swapIdx] = pack.Cards[swapIdx], pack.Cards[i]
	}
	for i, c := range yourCards {
		swapIdx := pack.IndexOf(c)
		targetIdx := 5 + i
		pack.Cards[targetIdx], pack.Cards[swapIdx] = pack.Cards[swapIdx], pack.Cards[targetIdx]
	}
}

func calcHighOutcome(playerOutcomes []PlayerOutcome, randomOpponentIdx int) *poker.HandOutcome {
	ourOutcome := playerOutcomes[0]
	var bestOpponentOutcome PlayerOutcome
	for i := 1; i < len(playerOutcomes); i++ {
		if i == 1 || poker.Beats(playerOutcomes[i].Level.HighLevel, bestOpponentOutcome.Level.HighLevel) {
			bestOpponentOutcome = playerOutcomes[i]
		}
	}
	randomOpponentOutcome := playerOutcomes[randomOpponentIdx]

	return &poker.HandOutcome{ourOutcome.IsHighWinner, bestOpponentOutcome.IsHighWinner, randomOpponentOutcome.IsHighWinner,
		ourOutcome.HighPotFractionWon, bestOpponentOutcome.HighPotFractionWon, randomOpponentOutcome.HighPotFractionWon,
		ourOutcome.Level.HighLevel, bestOpponentOutcome.Level.HighLevel, randomOpponentOutcome.Level.HighLevel}
}

func calcLowOutcome(playerOutcomes []PlayerOutcome, randomOpponentIdx int) *poker.HandOutcome {
	ourOutcome := playerOutcomes[0]
	var bestOpponentOutcome PlayerOutcome
	for i := 1; i < len(playerOutcomes); i++ {
		if i == 1 || poker.BeatsAceToFiveLow(playerOutcomes[i].Level.LowLevel, bestOpponentOutcome.Level.LowLevel) {
			bestOpponentOutcome = playerOutcomes[i]
		}
	}
	randomOpponentOutcome := playerOutcomes[randomOpponentIdx]

	return &poker.HandOutcome{ourOutcome.IsLowWinner, bestOpponentOutcome.IsLowWinner, randomOpponentOutcome.IsLowWinner,
		ourOutcome.LowPotFractionWon, bestOpponentOutcome.LowPotFractionWon, randomOpponentOutcome.LowPotFractionWon,
		ourOutcome.Level.LowLevel, bestOpponentOutcome.Level.LowLevel, randomOpponentOutcome.Level.LowLevel}
}

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
package poker

import (
	"math"
)

type HandOutcome struct {
	Won, OpponentWon, RandomOpponentWon                      bool
	PotFractionWon                                           float64
	BestOpponentPotFractionWon, RandomOpponentPotFractionWon float64
	OurLevel, BestOpponentLevel, RandomOpponentLevel         HandLevel
}

type Simulator struct {
	Players                int
	HandCount              int
	WinCount               int
	JointWinCount          int
	BestOpponentWinCount   int
	RandomOpponentWinCount int
	PotsWon                float64
	BestOpponentPotsWon    float64
	RandomOpponentPotsWon  float64

	OurClassCounts            []int
	BestOpponentClassCounts   []int
	RandomOpponentClassCounts []int

	ClassWinCounts        []int
	ClassJointWinCounts   []int
	ClassBestOppWinCounts []int
	ClassRandOppWinCounts []int

	BestHand          HandLevel
	BestOppHand       HandLevel
	ClassBestHands    []HandLevel
	ClassBestOppHands []HandLevel
}

func (s *Simulator) Reset(players, handsToPlay int) {
	s.Players = players
	s.HandCount = handsToPlay
	s.WinCount = 0
	s.JointWinCount = 0
	s.BestOpponentWinCount = 0
	s.RandomOpponentWinCount = 0
	s.PotsWon = 0
	s.BestOpponentPotsWon = 0
	s.RandomOpponentPotsWon = 0

	s.OurClassCounts = make([]int, MAX_HANDCLASS)
	s.BestOpponentClassCounts = make([]int, MAX_HANDCLASS)
	s.RandomOpponentClassCounts = make([]int, MAX_HANDCLASS)

	s.ClassWinCounts = make([]int, MAX_HANDCLASS)
	s.ClassJointWinCounts = make([]int, MAX_HANDCLASS)
	s.ClassBestOppWinCounts = make([]int, MAX_HANDCLASS)
	s.ClassRandOppWinCounts = make([]int, MAX_HANDCLASS)

	s.BestHand = MinLevel()
	s.BestOppHand = MinLevel()
	s.ClassBestHands = make([]HandLevel, MAX_HANDCLASS)
	for i := range s.ClassBestHands {
		s.ClassBestHands[i] = MinLevel()
	}
	s.ClassBestOppHands = make([]HandLevel, MAX_HANDCLASS)
	for i := range s.ClassBestOppHands {
		s.ClassBestOppHands[i] = MinLevel()
	}
}

func (s *Simulator) ProcessHand(outcome *HandOutcome) {
	if outcome.Won {
		s.WinCount++
		s.ClassWinCounts[outcome.OurLevel.Class]++
	}
	if outcome.OpponentWon {
		s.BestOpponentWinCount++
		s.ClassBestOppWinCounts[outcome.BestOpponentLevel.Class]++
	}
	if outcome.Won && outcome.OpponentWon {
		s.JointWinCount++
		s.ClassJointWinCounts[outcome.OurLevel.Class]++
	}
	if outcome.RandomOpponentWon {
		s.RandomOpponentWinCount++
		s.ClassRandOppWinCounts[outcome.RandomOpponentLevel.Class]++
	}
	s.PotsWon += outcome.PotFractionWon
	s.BestOpponentPotsWon += outcome.BestOpponentPotFractionWon
	s.RandomOpponentPotsWon += outcome.RandomOpponentPotFractionWon
	s.OurClassCounts[outcome.OurLevel.Class]++
	s.BestOpponentClassCounts[outcome.BestOpponentLevel.Class]++
	s.RandomOpponentClassCounts[outcome.RandomOpponentLevel.Class]++

	if Beats(outcome.OurLevel, s.BestHand) {
		s.BestHand = outcome.OurLevel
	}
	if Beats(outcome.BestOpponentLevel, s.BestOppHand) {
		s.BestOppHand = outcome.BestOpponentLevel
	}
	if Beats(outcome.OurLevel, s.ClassBestHands[outcome.OurLevel.Class]) {
		s.ClassBestHands[outcome.OurLevel.Class] = outcome.OurLevel
	}
	if Beats(outcome.BestOpponentLevel, s.ClassBestOppHands[outcome.BestOpponentLevel.Class]) {
		s.ClassBestOppHands[outcome.BestOpponentLevel.Class] = outcome.BestOpponentLevel
	}
}

func (s *Simulator) PotOddsBreakEven() float64 {
	return PotOddsBreakEven(s.PotsWon, s.HandCount)
}

// Calculate the largest bet which would have a positive expected value, relative to the size of the pot.
func PotOddsBreakEven(totalPotsWon float64, handCount int) float64 {
	// If W is the mean number of pots won, then:
	// Expected value of bet = W * (size of pot + bet size) - bet size
	// This is positive iff bet size < size of pot * W / (1 - W)
	meanPotsWon := totalPotsWon / float64(handCount)
	if meanPotsWon == 1.0 {
		return math.Inf(1)
	}
	return meanPotsWon / (1 - meanPotsWon)
}

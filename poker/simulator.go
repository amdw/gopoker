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
	"errors"
	"fmt"
)

type Simulator struct {
	Players                int
	HandCount              int
	WinCount               int
	JointWinCount          int
	BestOpponentWinCount   int
	RandomOpponentWinCount int

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

func (s *Simulator) SimulateHoldem(tableCards, yourCards []Card, players, handsToPlay int) {
	s.Players = players
	s.HandCount = handsToPlay
	s.WinCount = 0
	s.JointWinCount = 0
	s.BestOpponentWinCount = 0
	s.RandomOpponentWinCount = 0

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

	p := NewPack()
	for i := 0; i < handsToPlay; i++ {
		p.shuffleFixing(tableCards, yourCards)
		won, opponentWon, ourLevel, bestOpponentLevel, randomOpponentLevel := p.SimulateOneHoldemHand(players)
		if won {
			s.WinCount++
			s.ClassWinCounts[ourLevel.Class]++
		}
		if opponentWon {
			s.BestOpponentWinCount++
			s.ClassBestOppWinCounts[bestOpponentLevel.Class]++
		}
		if won && opponentWon {
			s.JointWinCount++
			s.ClassJointWinCounts[ourLevel.Class]++
		}
		if !Beats(ourLevel, randomOpponentLevel) && !Beats(bestOpponentLevel, randomOpponentLevel) {
			// The random opponent did at least as well as the winner
			s.RandomOpponentWinCount++
			s.ClassRandOppWinCounts[randomOpponentLevel.Class]++
		}
		s.OurClassCounts[ourLevel.Class]++
		s.BestOpponentClassCounts[bestOpponentLevel.Class]++
		s.RandomOpponentClassCounts[randomOpponentLevel.Class]++

		if Beats(ourLevel, s.BestHand) {
			s.BestHand = ourLevel
		}
		if Beats(bestOpponentLevel, s.BestOppHand) {
			s.BestOppHand = bestOpponentLevel
		}
		if Beats(ourLevel, s.ClassBestHands[ourLevel.Class]) {
			s.ClassBestHands[ourLevel.Class] = ourLevel
		}
		if Beats(bestOpponentLevel, s.ClassBestOppHands[bestOpponentLevel.Class]) {
			s.ClassBestOppHands[bestOpponentLevel.Class] = bestOpponentLevel
		}
	}
}

type StartingPair struct {
	Rank1, Rank2 Rank
	SameSuit     bool
}

func (pair StartingPair) Validate() error {
	if pair.Rank1 == pair.Rank2 && pair.SameSuit {
		return errors.New(fmt.Sprintf("Pair of %vs cannot be the same suit!", pair.Rank1))
	}
	return nil
}

func (pair StartingPair) SampleCards() (Card, Card) {
	err := pair.Validate()
	if err != nil {
		panic(err)
	}
	// Just pick arbitrary suits, either the same or different
	card1 := Card{pair.Rank1, Club}
	card2 := Card{pair.Rank2, Heart}
	if pair.SameSuit {
		card2.Suit = Club
	}
	return card1, card2
}

func (pair StartingPair) RunSimulation(players, handsToPlay int) *Simulator {
	card1, card2 := pair.SampleCards()
	result := &Simulator{}
	result.SimulateHoldem([]Card{}, []Card{card1, card2}, players, handsToPlay)
	return result
}

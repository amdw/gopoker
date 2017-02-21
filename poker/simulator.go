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
	"math"
	"math/rand"
	"time"
)

type Simulator struct {
	Players                int
	HandCount              int
	WinCount               int
	JointWinCount          int
	BestOpponentWinCount   int
	RandomOpponentWinCount int
	PotsWon                float64

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

func (s *Simulator) reset(players, handsToPlay int) {
	s.Players = players
	s.HandCount = handsToPlay
	s.WinCount = 0
	s.JointWinCount = 0
	s.BestOpponentWinCount = 0
	s.RandomOpponentWinCount = 0
	s.PotsWon = 0

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

func (s *Simulator) processHand(outcome HandOutcome) {
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

func (s *Simulator) SimulateHoldem(tableCards, yourCards []Card, players, handsToPlay int) {
	// Very crude attempt to detect situation where exhaustive enumeration is cheaper than simulation
	if len(tableCards) == 5 && len(yourCards) == 2 && players == 2 && handsToPlay > 990 {
		s.enumerateHoldem(tableCards, yourCards, players)
		return
	}
	s.reset(players, handsToPlay)
	p := NewPack()
	for i := 0; i < handsToPlay; i++ {
		p.shuffleFixing(tableCards, yourCards)
		handOutcome := p.SimulateOneHoldemHand(players)
		s.processHand(handOutcome)
	}
}

func (s *Simulator) enumerateHoldem(tableCards, yourCards []Card, players int) {
	s.reset(players, 0)
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	// For now we only enumerate the case where we have only one opponent and a full set of table cards.
	remainingPack := make([]Card, 45)
	i := 0
	isUsed := func(card Card) bool {
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
	for _, card := range NewPack().Cards {
		if !isUsed(card) {
			remainingPack[i] = card
			i++
		}
	}
	opponentHands := allCardCombinations(remainingPack, 2)
	for _, opponentHand := range opponentHands {
		playerCards := [][]Card{yourCards, opponentHand}
		outcomes := DealOutcomes(tableCards, playerCards)
		handOutcome := calcHandOutcome(outcomes, randGen)
		s.processHand(handOutcome)
		s.HandCount++
	}
}

// Calculate the largest bet which would have a positive expected value, relative to the size of the pot.
func (s *Simulator) PotOddsBreakEven() float64 {
	// If W is the mean number of pots won, then:
	// Expected value of bet = W * (size of pot + bet size) - bet size
	// This is positive iff bet size < size of pot * W / (1 - W)
	meanPotsWon := s.PotsWon / float64(s.HandCount)
	if meanPotsWon == 1.0 {
		return math.Inf(1)
	}
	return meanPotsWon / (1 - meanPotsWon)
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

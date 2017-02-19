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

func sortHands(outcomes []PlayerOutcome) {
	sort.Slice(outcomes, func(i, j int) bool {
		iBeatsJ := Beats(outcomes[i].Level, outcomes[j].Level)
		jBeatsI := Beats(outcomes[j].Level, outcomes[i].Level)
		if iBeatsJ && !jBeatsI {
			return true
		}
		if jBeatsI && !iBeatsJ {
			return false
		}
		return outcomes[i].Player < outcomes[j].Player
	})
}

type Pack struct {
	Cards   [52]Card
	randGen *rand.Rand
}

func (p *Pack) initialise() {
	// Not cryptographically secure, but fine for simulation, where performance is more important
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
// It is assumed that there are no duplicate cards in (tableCards+yourCards).
func (p *Pack) shuffleFixing(tableCards, yourCards []Card) {
	if len(tableCards) > 5 || len(yourCards) > 2 {
		panic(fmt.Sprintf("Maximum of 5 table cards and 2 hole cards supported, found %v and %v", len(tableCards), len(yourCards)))
	}

	indexOf := func(cards [52]Card, card Card) int {
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
	p.Shuffle()
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

func (p *Pack) Deal(players int) (onTable []Card, playerCards [][]Card) {
	if players < 1 {
		panic(fmt.Sprintf("At least one player required, found %v", players))
	}

	onTable = p.Cards[0:5]

	playerCards = make([][]Card, players)
	for player := 0; player < players; player++ {
		playerCards[player] = p.Cards[5+2*player : 7+2*player]
	}

	return onTable, playerCards
}

// Assess the hand each player holds and return a sorted list of outcomes
// by hand strength (descending) then player number (ascending).
// Player numbers are in ascending order of playerCards entries, starting with 1.
func DealOutcomes(onTable []Card, playerCards [][]Card) []PlayerOutcome {
	outcomes := make([]PlayerOutcome, len(playerCards))
	for playerIdx, hand := range playerCards {
		// In Hold'em it is not mandatory to use the cards in your hand,
		// so they are all optional
		combinedCards := make([]Card, 7)
		copy(combinedCards[0:5], onTable)
		copy(combinedCards[5:7], hand)
		level, cards := Classify([]Card{}, combinedCards)
		outcomes[playerIdx] = PlayerOutcome{playerIdx + 1, level, cards}
	}
	sortHands(outcomes)
	return outcomes
}

type SimulationResult struct {
	Won, OpponentWon, RandomOpponentWon              bool
	PotFractionWon                                   float64
	OurLevel, BestOpponentLevel, RandomOpponentLevel HandLevel
}

func calcSimResult(outcomes []PlayerOutcome, randGen *rand.Rand) SimulationResult {
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
		if !Beats(outcomes[0].Level, o.Level) {
			// The best hand doesn't beat this hand so it must be a winner
			potSplit++
		}
	}

	opponentWon := !Beats(ourOutcome.Level, opponentOutcomes[0].Level)

	randomOpponentLevel := opponentOutcomes[randGen.Intn(len(opponentOutcomes))].Level
	randomOpponentWon := !Beats(outcomes[0].Level, randomOpponentLevel)

	var potFractionWon float64
	if won {
		potFractionWon = 1.0 / float64(potSplit)
	} else {
		potFractionWon = 0
	}

	return SimulationResult{won, opponentWon, randomOpponentWon, potFractionWon,
		ourOutcome.Level, opponentOutcomes[0].Level, randomOpponentLevel}
}

// Play out one hand of Texas Hold'em and return whether or not player 1 won,
// plus player 1's hand level, plus the best hand level of any of player 1's opponents.
func (p *Pack) SimulateOneHoldemHand(players int) SimulationResult {
	onTable, playerCards := p.Deal(players)
	outcomes := DealOutcomes(onTable, playerCards)
	return calcSimResult(outcomes, p.randGen)
}

func NewPack() Pack {
	var result Pack
	result.initialise()
	return result
}

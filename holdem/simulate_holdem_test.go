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
	"math"
	"math/rand"
	"testing"
)

func TestSimSanity(t *testing.T) {
	players := 5
	simulations := 10000
	sim := SimulateHoldem([]poker.Card{}, []poker.Card{}, players, simulations)
	poker.TestAssertSimSanity(sim, players, simulations, t)
}

func TestTwoPlayers(t *testing.T) {
	simulations := 10000
	sim := SimulateHoldem([]poker.Card{}, []poker.Card{}, 2, simulations)
	poker.TestAssertSimSanity(sim, 2, simulations, t)
	// We can make some extra assertions here, as it's impossible for a pot to be split among opponents
	totalPotsWon := sim.PotsWon + sim.BestOpponentPotsWon
	if math.Abs(totalPotsWon-float64(simulations)) > 1e-6 {
		t.Errorf("Total pots won does not add up: %v us + %v best opponent = %v vs %v", sim.PotsWon, sim.BestOpponentPotsWon, totalPotsWon, simulations)
	}
	totalPotsWon = sim.PotsWon + sim.RandomOpponentPotsWon
	if math.Abs(totalPotsWon-float64(simulations)) > 1e-6 {
		t.Errorf("Total pots won does not add up: %v us + %v random opponent = %v vs %v", sim.PotsWon, sim.RandomOpponentPotsWon, totalPotsWon, simulations)
	}
}

func TestEnumeration(t *testing.T) {
	yourCards := h("9D", "7C")
	tableCards := h("KS", "7D", "AH", "8C", "8D")
	sim := SimulateHoldem(tableCards, yourCards, 2, 10000)

	// Should only do 45C2 = 990 simulations, one for each possible hand our opponent holds
	poker.TestAssertSimSanity(sim, 2, 990, t)

	if sim.OurClassCounts[poker.TwoPair] != 990 {
		t.Errorf("We got two pair but got %v not 990!", sim.OurClassCounts[poker.TwoPair])
	}
	if sim.BestOpponentClassCounts[poker.FourOfAKind] != 1 {
		t.Errorf("Opponent has one way to make quads but found %v", sim.BestOpponentClassCounts[poker.FourOfAKind])
	}
	// There are two 8s, two 7s, three As and three Ks which could be in our opponent's hand.
	// So they have the following ways to make a full house: 3 AAs, 3 KKs, one 77, six A8s, six K8s, four 78s
	// for a total of 23.
	if sim.BestOpponentClassCounts[poker.FullHouse] != 23 {
		t.Errorf("Opponent has 23 ways to make a full house but found %v", sim.BestOpponentClassCounts[poker.FullHouse])
	}
	// There are two 8s left in the deck and 43 other cards, so 86 ways to get exactly three 8s.
	// Of these, 16 also pair another table card giving a full house, so there are 70 ways to be "on trips".
	if sim.BestOpponentClassCounts[poker.ThreeOfAKind] != 70 {
		t.Errorf("Opponent has 70 ways to make trips but found %v", sim.BestOpponentClassCounts[poker.ThreeOfAKind])
	}
	// 3 kings * 40 non-K/8 cards = 120 ways to get kings and another pair
	// 2 sevens * 38 non-7/8/K cards = 76 ways to get sevens and another (non-K) pair
	// 3 aces * 35 non-A/8/K/7 cards = 105 ways to get aces and another (non-K/7) pair
	// 4C2 = 6 pairs for each of 8 values not represented, plus 3C2 = 3 nine pairs, for a total of 51 pocket pairs
	// Total 352 ways to make two pair
	if sim.BestOpponentClassCounts[poker.TwoPair] != 352 {
		t.Errorf("Opponent has 352 ways to make two pair but found %v", sim.BestOpponentClassCounts[poker.TwoPair])
	}
	// And 990 minus all the other scenarios gives 544
	if sim.BestOpponentClassCounts[poker.OnePair] != 544 {
		t.Errorf("Opponent has 544 ways to make one pair but found %v", sim.BestOpponentClassCounts[poker.OnePair])
	}
}

func TestHandOutcomeSanity(t *testing.T) {
	p := poker.NewPack()
	tests := 1000
	randGen := rand.New(rand.NewSource(1234)) // Deterministic for repeatable tests
	for i := 0; i < tests; i++ {
		p.Shuffle(randGen)
		res := SimulateOneHoldemHand(&p, 5, randGen)
		if !res.Won && !res.OpponentWon {
			t.Errorf("Simulator says nobody won!")
		}
		if poker.Beats(res.RandomOpponentLevel, res.BestOpponentLevel) {
			t.Errorf("Random opponent level %v beats best level %v", res.RandomOpponentLevel, res.BestOpponentLevel)
		}
		if res.PotFractionWon < 0 || res.PotFractionWon > 1 {
			t.Errorf("Pot fraction won must be between 0 and 1: %v", res.PotFractionWon)
		}
		if res.BestOpponentPotFractionWon < 0 || res.BestOpponentPotFractionWon > 1 {
			t.Errorf("Best opponent pot fraction won must be between 0 and 1: %v", res.BestOpponentPotFractionWon)
		}
		if res.RandomOpponentPotFractionWon < 0 || res.RandomOpponentPotFractionWon > 1 {
			t.Errorf("Random opponent pot fraction won must be between 0 and 1: %v", res.RandomOpponentPotFractionWon)
		}
		if res.Won {
			if poker.Beats(res.BestOpponentLevel, res.OurLevel) {
				t.Errorf("Simulator says we won but best opponent %v beats our %v", res.BestOpponentLevel, res.OurLevel)
			}
			if poker.Beats(res.RandomOpponentLevel, res.OurLevel) {
				t.Errorf("Simulator says we won but random opponent %v beats our %v", res.RandomOpponentLevel, res.OurLevel)
			}
			if math.Abs(res.PotFractionWon) < 1e-6 {
				t.Errorf("Simulator says we won but we didn't win any of the pot: %v", res.PotFractionWon)
			}
			if res.OpponentWon && math.Abs(res.PotFractionWon-res.BestOpponentPotFractionWon) > 1e-6 {
				t.Errorf("Simulator says both we and best opponent won but we won different amounts: %v vs %v", res.PotFractionWon, res.BestOpponentPotFractionWon)
			}
			if res.RandomOpponentWon && math.Abs(res.PotFractionWon-res.RandomOpponentPotFractionWon) > 1e-6 {
				t.Errorf("Simulator says both we and random opponent won but we won different amounts: %v vs %v", res.PotFractionWon, res.RandomOpponentPotFractionWon)
			}
		} else {
			if !poker.Beats(res.BestOpponentLevel, res.OurLevel) {
				t.Errorf("Simulator says we lost but their %v doesn't beat our %v", res.BestOpponentLevel, res.OurLevel)
			}
			if res.PotFractionWon != 0 {
				t.Errorf("Simulator says we didn't win, but we won chips: %v", res.PotFractionWon)
			}
		}
		if res.OpponentWon {
			if poker.Beats(res.OurLevel, res.BestOpponentLevel) {
				t.Errorf("Simulator says opponent won but our %v beats their %v", res.OurLevel, res.BestOpponentLevel)
			}
			if poker.Beats(res.RandomOpponentLevel, res.BestOpponentLevel) {
				t.Errorf("Simulator says opponent won but random opponent %v beats their %v", res.RandomOpponentLevel, res.BestOpponentLevel)
			}
			if math.Abs(res.PotFractionWon-1.0) < 1e-6 {
				t.Errorf("Simulator says opponent won but we won the whole pot: %v", res.PotFractionWon)
			}
			if math.Abs(res.BestOpponentPotFractionWon) < 1e-6 {
				t.Errorf("Simulator says opponent won, but they didn't win any of the pot: %v", res.BestOpponentPotFractionWon)
			}
			if res.RandomOpponentWon && math.Abs(res.BestOpponentPotFractionWon-res.RandomOpponentPotFractionWon) > 1e-6 {
				t.Errorf("Simulator says both best and random opponent won, but they won different amounts: %v vs %v", res.BestOpponentPotFractionWon, res.RandomOpponentPotFractionWon)
			}
		} else {
			if res.PotFractionWon != 1.0 {
				t.Errorf("Simulator says we were the sole winner but we didn't win the whole pot: %v", res.PotFractionWon)
			}
			if res.BestOpponentPotFractionWon != 0 {
				t.Errorf("Simulator says we were the sole winner but best opponent won chips: %v", res.BestOpponentPotFractionWon)
			}
			if res.RandomOpponentPotFractionWon != 0 {
				t.Errorf("Simulator says we were the sole winner but random opponent won chips: %v", res.RandomOpponentPotFractionWon)
			}
		}
		if res.RandomOpponentWon {
			if !res.OpponentWon {
				t.Errorf("Simulator says random opponent won (%v) but not best opponent (%v)!", res.RandomOpponentLevel, res.BestOpponentLevel)
			}
			if poker.Beats(res.OurLevel, res.RandomOpponentLevel) {
				t.Errorf("Simulator says random opponent won but our %v beats their %v", res.OurLevel, res.RandomOpponentLevel)
			}
			if poker.Beats(res.BestOpponentLevel, res.RandomOpponentLevel) {
				t.Errorf("Simulator says random opponent won but best opponent %v beats their %v", res.BestOpponentLevel, res.RandomOpponentLevel)
			}
			if math.Abs(res.PotFractionWon-1.0) < 1e-6 {
				t.Errorf("Simulator says random opponent won but we won the whole pot: %v", res.PotFractionWon)
			}
			if math.Abs(res.RandomOpponentPotFractionWon) < 1e-6 {
				t.Errorf("Simulator says random opponent won but they didn't win any chips: %v", res.RandomOpponentPotFractionWon)
			}
		}
	}
}

func TestFixedShuffle(t *testing.T) {
	pack := poker.NewPack()
	randGen := rand.New(rand.NewSource(1234)) // Deterministic for repeatable tests

	myCards := h("KS", "AC")
	tableCards := h("10D", "2C", "AS", "4D", "6H")
	for testNum := 0; testNum < 1000; testNum++ {
		shuffleFixing(&pack, tableCards, myCards, randGen)
		tCards, pCards := Deal(&pack, 5)
		if !poker.CardsEqual(tableCards, tCards) {
			t.Errorf("Expected table cards %q, found %q", tableCards, tCards)
		}
		if !poker.CardsEqual(myCards, pCards[0]) {
			t.Errorf("Expected player 1 cards %q, found %q", myCards, pCards[0])
		}
		containsAny := func(cards, testCards []poker.Card) bool {
			for _, c := range cards {
				for _, tc := range testCards {
					if c == tc {
						return true
					}
				}
			}
			return false
		}
		for i := 1; i < len(pCards); i++ {
			if containsAny(pCards[i], myCards) {
				t.Errorf("Player %v's cards %q should not contain any of player 1's cards %q.", i+1, pCards[i], myCards)
			}
			if containsAny(pCards[i], tableCards) {
				t.Errorf("Player %v's cards %q should not contain any table cards %q", i+1, pCards[i], tableCards)
			}
		}
		poker.TestPackPermutation(&pack, t)
	}
}

func sp(r1s, r2s string, suited bool) StartingPair {
	r1, err := poker.MakeRank(r1s)
	if err != nil {
		panic(fmt.Sprintf("Cannot make rank from %v", r1s))
	}
	r2, err := poker.MakeRank(r2s)
	if err != nil {
		panic(fmt.Sprintf("Cannot make rank from %v", r2s))
	}
	return StartingPair{r1, r2, suited}
}

func TestPairs(t *testing.T) {
	pairs := []StartingPair{sp("K", "Q", false), sp("K", "Q", true), sp("K", "K", false)}
	players := 6
	simCount := 1000
	for _, pair := range pairs {
		sim := pair.RunSimulation(players, simCount)
		poker.TestAssertSimSanity(sim, players, simCount, t)
	}
}

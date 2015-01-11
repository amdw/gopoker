package poker

import (
	"math/rand"
	"testing"
)

func TestFixedShuffle(t *testing.T) {
	pack := NewPack()
	myCards := h("KS", "AC")
	tableCards := h("10D", "2C", "AS")
	pack.randGen = rand.New(rand.NewSource(1234)) // Deterministic for repeatable tests

	for i := 0; i < 100; i++ {
		pack.Shuffle()
		pack.shuffleFixing(tableCards, myCards)
		for j, c := range h("10D", "2C", "AS") {
			if pack.Cards[j] != c {
				t.Fatalf("Expected %v at Cards[%v], found %v", c, j, pack.Cards[j])
			}
		}
		for j, c := range h("KS", "AC") {
			if pack.Cards[j+5] != c {
				t.Fatalf("Expected %v at Cards[%v], found %v", c, j+5, pack.Cards[j+5])
			}
		}
		permCheck := make([][]int, 4)
		for i := 0; i < 4; i++ {
			permCheck[i] = make([]int, 13)
		}
		for _, c := range pack.Cards {
			permCheck[c.Suit][c.Rank]++
		}
		for s := range permCheck {
			for r, count := range permCheck[s] {
				if count != 1 {
					t.Fatalf("Expected exactly one %v%v in pack after shuffle, found %v", Rank(r).String(), Suit(s).String(), count)
				}
			}
		}
	}
}

func TestSimInternalSanity(t *testing.T) {
	p := NewPack()
	tests := 1000
	for i := 0; i < tests; i++ {
		p.Shuffle()
		won, ourLevel, bestOpponentLevel, randomOpponentLevel := p.SimulateOneHoldemHand(5)
		if Beats(randomOpponentLevel, bestOpponentLevel) {
			t.Errorf("Random opponent level %v beats best level %v", randomOpponentLevel, bestOpponentLevel)
		}
		if won {
			if Beats(bestOpponentLevel, ourLevel) {
				t.Errorf("Simulator says we won but best opponent %v beats our %v", bestOpponentLevel, ourLevel)
			}
			if Beats(randomOpponentLevel, ourLevel) {
				t.Errorf("Simulator says we won but random opponent %v beats our %v", randomOpponentLevel, ourLevel)
			}
		} else {
			if Beats(ourLevel, bestOpponentLevel) {
				t.Errorf("Simulator says we didn't win but %v beats %v", ourLevel, bestOpponentLevel)
			}
		}
	}
}

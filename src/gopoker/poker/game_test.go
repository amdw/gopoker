package poker

import (
	"testing"
)

func TestSimInternalSanity(t *testing.T) {
	p := NewPack()
	tests := 1000
	for i := 0; i < tests; i++ {
		p.Shuffle()
		won, ourLevel, bestOpponentLevel := p.SimulateOneHoldemHand(5)
		if won {
			if Beats(bestOpponentLevel, ourLevel) {
				t.Errorf("Simulator says we won but %v beats %v", bestOpponentLevel, ourLevel)
			}
		} else {
			if Beats(ourLevel, bestOpponentLevel) {
				t.Errorf("Simulator says we didn't win but %v beats %v", ourLevel, bestOpponentLevel)
			}
		}
	}
}

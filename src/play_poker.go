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
along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"flag"
	"fmt"
	"poker"
)

// Evaluate a few deals of Texas Hold'em
func main() {
	// Initialise pack
	pack := poker.NewPack()

	var handsToPlay int
	var verbose bool
	flag.IntVar(&handsToPlay, "hands", 10000, "How many hands to play?")
	flag.BoolVar(&verbose, "verbose", false, "Print the result of every hand?")
	flag.Parse()

	frequencies := make(map[poker.HandClass]int)

	for handCount := 1; handCount <= handsToPlay; handCount++ {
		pack.Shuffle()
		onTable, playerCards, sortedOutcomes := pack.PlayHoldem(1)
		if verbose {
			fmt.Printf("On table: %v; in hand: %v; outcome: %v %v\n", onTable, playerCards[0], sortedOutcomes.Outcomes[0].Level, sortedOutcomes.Outcomes[0].Cards)
		}
		frequencies[sortedOutcomes.Outcomes[0].Level.Class]++
	}
	fmt.Printf("Hand frequencies (total %v):\n", handsToPlay)
	for c := poker.MAX_HANDCLASS - 1; c >= 0; c-- {
		class := poker.HandClass(c)
		freq := frequencies[class]
		fmt.Printf("%v\t%v (%.2f%%)\n", class, freq, float32(freq)*100.0/float32(handsToPlay))
	}
}

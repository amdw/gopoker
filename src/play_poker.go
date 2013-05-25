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

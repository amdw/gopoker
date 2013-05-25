package main

import (
	"flag"
	"fmt"
	"math/rand"
	"poker"
	"time"
)

// Evaluate a few deals of Texas Hold'em
func main() {
	// Initialise pack
	pack := make([]poker.Card, 52)
	i := 0
	for s := 0; s < 4; s++ {
		for r := 0; r < 13; r++ {
			pack[i] = poker.Card{poker.Rank(r), poker.Suit(s)}
			i++
		}
	}

	var handsToPlay int
	var verbose bool
	flag.IntVar(&handsToPlay, "hands", 10000, "How many hands to play?")
	flag.BoolVar(&verbose, "verbose", false, "Print the result of every hand?")
	flag.Parse()

	frequencies := make(map[poker.HandClass]int)
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for handCount := 1; handCount <= handsToPlay; handCount++ {
		// Shuffle (we only care about the first 7 cards)
		for i := 0; i < 7; i++ {
			j := randGen.Intn(52-i) + i
			pack[i], pack[j] = pack[j], pack[i]
		}
		result := poker.Classify(pack[0:2], pack[2:7])
		if verbose {
			fmt.Printf("Hand %v: %q\n", handCount, result)
		}
		frequencies[result.Class]++
	}
	fmt.Printf("Hand frequencies (total %v):\n", handsToPlay)
	for class, freq := range frequencies {
		fmt.Printf("%v\t%v\n", class, freq)
	}
}

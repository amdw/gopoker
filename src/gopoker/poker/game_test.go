package poker

import (
	"math/rand"
	"testing"
)

type GameTest struct {
	HandCards     []Card
	TableCards    []Card
	ExpectedLevel HandLevel
	ExpectedCards []Card
}

var gameTests = []GameTest{
	GameTest{h("10S", "JS"), h("2H", "QS", "6D", "KS", "AS"), HandLevel{StraightFlush, []Rank{Ace}}, h("10S", "JS", "QS", "KS", "AS")},
	GameTest{h("AH", "3H"), h("2H", "4H", "5H", "2C", "3C"), HandLevel{StraightFlush, []Rank{Five}}, h("AH", "2H", "3H", "4H", "5H")},
	GameTest{h("10S", "JS"), h("10C", "QD", "10D", "10H", "2S"), HandLevel{FourOfAKind, []Rank{Ten, Queen}}, h("10S", "10C", "10D", "10H", "QD")},
	GameTest{h("10S", "10C"), h("10D", "10H", "JD", "QD", "KD"), HandLevel{FourOfAKind, []Rank{Ten, King}}, h("10S", "10C", "10D", "10H", "KD")},
	GameTest{h("2S", "2H"), h("3D", "3H", "3C", "QD", "KS"), HandLevel{FullHouse, []Rank{Three, Two}}, h("2S", "2H", "3H", "3C", "3D")},
	GameTest{h("2S", "3S"), h("4C", "4S", "2D", "2H", "3C"), HandLevel{FullHouse, []Rank{Two, Four}}, h("2S", "2D", "2H", "4C", "4S")},
	GameTest{h("6H", "8H"), h("9H", "10H", "2H", "3S", "7C"), HandLevel{Flush, []Rank{Ten, Nine, Eight, Six, Two}}, h("10H", "9H", "8H", "6H", "2H")},
	GameTest{h("6S", "8H"), h("9H", "10H", "JH", "QH", "7H"), HandLevel{StraightFlush, []Rank{Queen}}, h("QH", "JH", "10H", "9H", "8H")},
	GameTest{h("AS", "JH"), h("QC", "KD", "10S", "2C", "3C"), HandLevel{Straight, []Rank{Ace, King, Queen, Jack, Ten}}, h("AS", "KD", "QC", "JH", "10S")},
	GameTest{h("AS", "3H"), h("2C", "4C", "5D", "KS", "JC"), HandLevel{Straight, []Rank{Five, Four, Three, Two, Ace}}, h("AS", "2C", "3H", "4C", "5D")},
	GameTest{h("6S", "6D"), h("6C", "KH", "JC", "7H", "2S"), HandLevel{ThreeOfAKind, []Rank{Six, King, Jack}}, h("6S", "6D", "6C", "KH", "JC")},
	GameTest{h("6S", "2S"), h("6C", "KH", "JC", "7H", "6D"), HandLevel{ThreeOfAKind, []Rank{Six, King, Jack}}, h("6S", "6D", "6C", "KH", "JC")},
	GameTest{h("6S", "4D"), h("6D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}}, h("6S", "6D", "4D", "4S", "AH")},
	GameTest{h("6S", "6D"), h("4D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}}, h("6S", "6D", "4D", "4S", "AH")},
	GameTest{h("AS", "AH"), h("2S", "4C", "6D", "8S", "10D"), HandLevel{OnePair, []Rank{Ace, Ten, Eight, Six}}, h("AS", "AH", "10D", "8S", "6D")},
	GameTest{h("AS", "2S"), h("AH", "4C", "6D", "8S", "10D"), HandLevel{OnePair, []Rank{Ace, Ten, Eight, Six}}, h("AS", "AH", "10D", "8S", "6D")},
	GameTest{h("2S", "4S"), h("5D", "7S", "8S", "QH", "KH"), HandLevel{HighCard, []Rank{King, Queen, Eight, Seven, Five}}, h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("2S", "KH"), h("5D", "7S", "8S", "QH", "4S"), HandLevel{HighCard, []Rank{King, Queen, Eight, Seven, Five}}, h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("8S", "KH"), h("5D", "7S", "2S", "QH", "4S"), HandLevel{HighCard, []Rank{King, Queen, Eight, Seven, Five}}, h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("KS", "8C"), h("QH", "10C", "9H", "7H", "6S"), HandLevel{Straight, []Rank{Ten, Nine, Eight, Seven, Six}}, h("10C", "9H", "8C", "7H", "6S")},
	GameTest{h("KS", "10H"), h("8C", "QC", "9H", "7H", "6S"), HandLevel{Straight, []Rank{Ten, Nine, Eight, Seven, Six}}, h("10H", "9H", "8C", "7H", "6S")},
	GameTest{h("KS", "10H"), h("8C", "5C", "9H", "7H", "6S"), HandLevel{Straight, []Rank{Ten, Nine, Eight, Seven, Six}}, h("10H", "9H", "8C", "7H", "6S")},
	// Uses none of the hand cards
	GameTest{h("2S", "3S"), h("AH", "10H", "JH", "QH", "KH"), HandLevel{StraightFlush, []Rank{Ace}}, h("AH", "KH", "QH", "JH", "10H")},
	// Uses one of the hand cards
	GameTest{h("AH", "3S"), h("2S", "10H", "JH", "QH", "KH"), HandLevel{StraightFlush, []Rank{Ace}}, h("AH", "KH", "QH", "JH", "10H")},
}

func TestHoldem(t *testing.T) {
	playerCount := 2
	for i, testCase := range gameTests {
		if len(testCase.TableCards) != 5 {
			t.Fatalf("Wrong number of table cards %v", i)
		}
		pack := NewPack()
		pack.Shuffle()
		pack.shuffleFixing(testCase.TableCards, testCase.HandCards)
		tableCards, playerCards, outcomes := pack.PlayHoldem(playerCount)
		if len(tableCards) != 5 {
			t.Fatalf("Expected 5 table cards, found %v", len(tableCards))
		}
		for j, c := range tableCards {
			if c != testCase.TableCards[j] {
				t.Errorf("Expected %v at tableCards[%v], found %v", testCase.TableCards[j], j, c)
			}
		}
		if len(playerCards) != playerCount {
			t.Fatalf("Expected %v sets of player cards, found %v", playerCount, len(playerCards))
		}
		for j, c := range playerCards[0] {
			if c != testCase.HandCards[j] {
				t.Errorf("Expected %v at playerCards[0][%v], found %v", testCase.HandCards[j], j, c)
			}
		}
		if len(outcomes) != playerCount {
			t.Fatalf("Expected %v outcomes, found %v", playerCount, len(outcomes))
		}
		foundPlayer := false
		for _, outcome := range outcomes {
			if outcome.Player != 1 {
				continue
			}
			foundPlayer = true
			if !levelsEqual(outcome.Level, testCase.ExpectedLevel) {
				t.Errorf("Expected %v for case %v, found %v", testCase.ExpectedLevel, i, outcome.Level)
			}
			if len(outcome.Cards) != 5 {
				t.Fatalf("Expected 5 cards for case %v, found %v", i, len(outcome.Cards))
			}
			if !cardsEqual(testCase.ExpectedCards, outcome.Cards) {
				t.Errorf("Expected cards %q for case %v, found %q", testCase.ExpectedCards, i, outcome.Cards)
			}
		}
		if !foundPlayer {
			t.Errorf("Could not find player 1 in outcomes")
		}
	}
}

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

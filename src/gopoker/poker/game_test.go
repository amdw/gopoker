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
	GameTest{h("AS", "JH"), h("QC", "KD", "10S", "2C", "3C"), HandLevel{Straight, []Rank{Ace}}, h("AS", "KD", "QC", "JH", "10S")},
	GameTest{h("AS", "3H"), h("2C", "4C", "5D", "KS", "JC"), HandLevel{Straight, []Rank{Five}}, h("AS", "2C", "3H", "4C", "5D")},
	GameTest{h("6S", "6D"), h("6C", "KH", "JC", "7H", "2S"), HandLevel{ThreeOfAKind, []Rank{Six, King, Jack}}, h("6S", "6D", "6C", "KH", "JC")},
	GameTest{h("6S", "2S"), h("6C", "KH", "JC", "7H", "6D"), HandLevel{ThreeOfAKind, []Rank{Six, King, Jack}}, h("6S", "6D", "6C", "KH", "JC")},
	GameTest{h("6S", "4D"), h("6D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}}, h("6S", "6D", "4D", "4S", "AH")},
	GameTest{h("6S", "6D"), h("4D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}}, h("6S", "6D", "4D", "4S", "AH")},
	GameTest{h("AS", "AH"), h("2S", "4C", "6D", "8S", "10D"), HandLevel{OnePair, []Rank{Ace, Ten, Eight, Six}}, h("AS", "AH", "10D", "8S", "6D")},
	GameTest{h("AS", "2S"), h("AH", "4C", "6D", "8S", "10D"), HandLevel{OnePair, []Rank{Ace, Ten, Eight, Six}}, h("AS", "AH", "10D", "8S", "6D")},
	GameTest{h("2S", "4S"), h("5D", "7S", "8S", "QH", "KH"), HandLevel{HighCard, []Rank{King, Queen, Eight, Seven, Five}}, h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("2S", "KH"), h("5D", "7S", "8S", "QH", "4S"), HandLevel{HighCard, []Rank{King, Queen, Eight, Seven, Five}}, h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("8S", "KH"), h("5D", "7S", "2S", "QH", "4S"), HandLevel{HighCard, []Rank{King, Queen, Eight, Seven, Five}}, h("5D", "7S", "8S", "QH", "KH")},
	GameTest{h("KS", "8C"), h("QH", "10C", "9H", "7H", "6S"), HandLevel{Straight, []Rank{Ten}}, h("10C", "9H", "8C", "7H", "6S")},
	GameTest{h("KS", "10H"), h("8C", "QC", "9H", "7H", "6S"), HandLevel{Straight, []Rank{Ten}}, h("10H", "9H", "8C", "7H", "6S")},
	GameTest{h("KS", "10H"), h("8C", "5C", "9H", "7H", "6S"), HandLevel{Straight, []Rank{Ten}}, h("10H", "9H", "8C", "7H", "6S")},
	// Uses none of the hand cards
	GameTest{h("2S", "3S"), h("AH", "10H", "JH", "QH", "KH"), HandLevel{StraightFlush, []Rank{Ace}}, h("AH", "KH", "QH", "JH", "10H")},
	// Uses one of the hand cards
	GameTest{h("AH", "3S"), h("2S", "10H", "JH", "QH", "KH"), HandLevel{StraightFlush, []Rank{Ace}}, h("AH", "KH", "QH", "JH", "10H")},
}

func TestHoldem(t *testing.T) {
	playerCount := 4
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
		for j := 1; j < len(outcomes); j++ {
			for k := 0; k < j; k++ {
				prevLevel := outcomes[k].Level
				thisLevel := outcomes[j].Level
				if Beats(thisLevel, prevLevel) {
					t.Errorf("%v beats %v but came after it", thisLevel, prevLevel)
				}
				if !Beats(thisLevel, prevLevel) && !Beats(prevLevel, thisLevel) && outcomes[k].Player > outcomes[j].Player {
					t.Errorf("%v and %v have equal hands (%v, %v) but %v came first", outcomes[k].Player, outcomes[j].Player, outcomes[k].Cards, outcomes[j].Cards, outcomes[k].Player)
				}
			}
		}
	}
}

func TestFixedShuffle(t *testing.T) {
	pack := NewPack()
	pack.randGen = rand.New(rand.NewSource(1234)) // Deterministic for repeatable tests

	myCards := h("KS", "AC")
	tableCards := h("10D", "2C", "AS", "4D", "6H")
	pack.Shuffle()
	pack.shuffleFixing(tableCards, myCards)
	tCards, pCards, _ := pack.PlayHoldem(2)
	if !cardsEqual(tableCards, tCards) {
		t.Errorf("Expected table cards %q, found %q", tableCards, tCards)
	}
	if !cardsEqual(myCards, pCards[0]) {
		t.Errorf("Expected player 1 cards %q, found %q", myCards, pCards[0])
	}
	containsAny := func(cards, testCards []Card) bool {
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

func TestSimInternalSanity(t *testing.T) {
	p := NewPack()
	tests := 1000
	for i := 0; i < tests; i++ {
		p.Shuffle()
		won, opponentWon, ourLevel, bestOpponentLevel, randomOpponentLevel := p.SimulateOneHoldemHand(5)
		if !won && !opponentWon {
			t.Errorf("Simulator says nobody won!")
		}
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
			if !Beats(bestOpponentLevel, ourLevel) {
				t.Errorf("Simulator says we lost but their %v doesn't beat our %v", bestOpponentLevel, ourLevel)
			}
		}
		if opponentWon {
			if Beats(ourLevel, bestOpponentLevel) {
				t.Errorf("Simulator says opponent won but our %v beats their %v", ourLevel, bestOpponentLevel)
			}
			if Beats(randomOpponentLevel, bestOpponentLevel) {
				t.Errorf("Simulator says opponent won but random opponent %v beats their %v", randomOpponentLevel, bestOpponentLevel)
			}
		}
	}
}

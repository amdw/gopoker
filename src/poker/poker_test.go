package poker

import (
	"fmt"
	"sort"
	"testing"
)

func TestCardBasics(t *testing.T) {
	for s := Heart; s <= Club; s++ {
		for r := Two; r <= Ace; r++ {
			cs := fmt.Sprintf("%v%v", r, s)
			c1 := C(cs)
			c2 := Card{r, s}
			if c1 != c2 {
				t.Errorf("Expected %q, found %q", c2, c1)
			}
			if c1.String() != cs {
				t.Errorf("Expected %q, found %q", cs, c1.String())
			}
		}
	}
}

type LevelTest struct {
	l1        HandLevel
	l2        HandLevel
	isGreater bool
	isLess    bool
}

func hl(class HandClass, tiebreaks []Rank) HandLevel {
	return HandLevel{class, tiebreaks, []Card{}}
}

var levelTests = []LevelTest{
	{hl(StraightFlush, []Rank{Ace}), hl(StraightFlush, []Rank{Ace}), false, false},
	{hl(StraightFlush, []Rank{Ace}), hl(StraightFlush, []Rank{King}), true, false},
	{hl(FourOfAKind, []Rank{Nine, Ten}), hl(StraightFlush, []Rank{Two}), false, true},
}

func TestLevels(t *testing.T) {
	for _, ltst := range levelTests {
		gt := Beats(ltst.l1, ltst.l2)
		lt := Beats(ltst.l2, ltst.l1)
		if gt != ltst.isGreater {
			t.Errorf("Expected %q beats %q == %v, found %v", ltst.l1, ltst.l2, ltst.isGreater, gt)
		}
		if lt != ltst.isLess {
			t.Errorf("Expected %q beats %q == %v, found %v", ltst.l2, ltst.l1, ltst.isLess, lt)
		}
	}
}

type ClassTest struct {
	mandatory []Card
	optional  []Card
	expected  HandLevel
}

func h(cards ...string) []Card {
	result := make([]Card, len(cards))
	for i, c := range cards {
		result[i] = C(c)
	}
	return result
}

var classTests = []ClassTest{
	ClassTest{h("AS", "KS", "QS", "JS", "10S"), h(), HandLevel{StraightFlush, []Rank{Ace}, h("AS", "KS", "QS", "JS", "10S")}},
	ClassTest{h("9D", "10S", "9S", "9H", "9C"), h(), HandLevel{FourOfAKind, []Rank{Nine, Ten}, h("9D", "10S", "9S", "9H", "9C")}},
	ClassTest{h("10S", "JS"), h("2H", "QS", "6D", "KS", "AS"), HandLevel{StraightFlush, []Rank{Ace}, h("10S", "JS", "QS", "KS", "AS")}},
	ClassTest{h("10S", "JS"), h("10C", "QD", "10D", "10H", "2S"), HandLevel{FourOfAKind, []Rank{Ten, Jack}, h("10S", "10C", "10D", "10H", "JS")}},
	ClassTest{h("10S", "10C"), h("10D", "10H", "JD", "QD", "KD", "AD"), HandLevel{FourOfAKind, []Rank{Ten, Ace}, h("10S", "10C", "10D", "10H", "AD")}},
	ClassTest{h("2S", "2H"), h("3D", "3H", "3C", "QD", "KS"), HandLevel{FullHouse, []Rank{Three, Two}, h("2S", "2H", "3H", "3C", "3D")}},
	ClassTest{h("2S", "3S"), h("4H", "4D", "4C", "4S", "2D", "2H", "3C"), HandLevel{FullHouse, []Rank{Two, Three}, h("2S", "2D", "2H", "3S", "3C")}},
	ClassTest{h("6H", "8H"), h("9H", "10H", "2H", "3S", "7C"), HandLevel{Flush, []Rank{Ten, Nine, Eight, Six, Two}, h("10H", "9H", "8H", "6H", "2H")}},
	ClassTest{h("6S", "8H"), h("9H", "10H", "JH", "QH", "7H"), HandLevel{Straight, []Rank{Ten, Nine, Eight, Seven, Six}, h("10H", "9H", "8H", "7H", "6S")}},
	ClassTest{h("AS", "3H"), h("2C", "4C", "5D", "KS", "JC"), HandLevel{Straight, []Rank{Five, Four, Three, Two, Ace}, h("AS", "2C", "3H", "4C", "5D")}},
	ClassTest{h("6S", "6D"), h("6C", "KH", "JC", "7H", "2S"), HandLevel{ThreeOfAKind, []Rank{Six, King, Jack}, h("6S", "6D", "6C", "KH", "JC")}},
	ClassTest{h("6S", "2S"), h("6C", "KH", "JC", "7H", "6D"), HandLevel{ThreeOfAKind, []Rank{Six, King, Two}, h("6S", "6D", "6C", "KH", "2S")}},
	ClassTest{h("6S", "4D"), h("6D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}, h("6S", "6D", "4D", "4S", "AH")}},
	ClassTest{h("6S", "6D"), h("4D", "QS", "4S", "AH", "3C"), HandLevel{TwoPair, []Rank{Six, Four, Ace}, h("6S", "6D", "4D", "4S", "AH")}},
	ClassTest{h("KS", "8C"), h("10H", "10C", "9H", "7H", "6S"), HandLevel{OnePair, []Rank{Ten, King, Nine, Eight}, h("KS", "8C", "10H", "10C", "9H")}},
	ClassTest{h("KS", "10H"), h("8C", "10C", "9H", "7H", "6S"), HandLevel{OnePair, []Rank{Ten, King, Nine, Eight}, h("KS", "8C", "10H", "10C", "9H")}},
	ClassTest{h("KS", "10H"), h("8C", "5C", "9H", "7H", "6S"), HandLevel{HighCard, []Rank{King, Ten, Nine, Eight, Seven}, h("KS", "10H", "9H", "8C", "7H")}},
}

func levelsEqual(l1 HandLevel, l2 HandLevel) bool {
	if l1.Class != l2.Class {
		return false
	}
	if len(l1.Tiebreaks) != len(l2.Tiebreaks) {
		return false
	}
	for i := 0; i < len(l1.Tiebreaks); i++ {
		if l1.Tiebreaks[i] != l2.Tiebreaks[i] {
			return false
		}
	}
	if len(l1.Cards) != len(l2.Cards) {
		return false
	}
	sort.Sort(CardSorter{l1.Cards, false})
	sort.Sort(CardSorter{l2.Cards, false})
	for i := 0; i < len(l1.Cards); i++ {
		if l1.Cards[i] != l2.Cards[i] {
			return false
		}
	}
	return true
}

func TestClassification(t *testing.T) {
	for _, ct := range classTests {
		c := Classify(ct.mandatory, ct.optional)
		if !levelsEqual(ct.expected, c) {
			t.Errorf("Expected %q, found %q for %q / %q", ct.expected, c, ct.mandatory, ct.optional)
		}
	}
}

func TestBuildStraights(t *testing.T) {
	suits := [][]Suit{{Spade}, {Club, Heart}, {Spade}, {Spade}, {Spade}}
	expectedOutput := [][]Suit{{Spade, Club, Spade, Spade, Spade}, {Spade, Heart, Spade, Spade, Spade}}
	actualOutput := buildStraights(suits)
	for i := range expectedOutput {
		for j := range expectedOutput[i] {
			if expectedOutput[i][j] != actualOutput[i][j] {
				t.Errorf("Expected %q at %v but found %q", expectedOutput[i], i, actualOutput[i])
				break
			}
		}
	}
}

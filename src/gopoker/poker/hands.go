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

This file contains code dealing with basic rules of poker: hand rankings, classifications and so forth.
*/
package poker

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Suit int8

const (
	Heart Suit = iota
	Diamond
	Spade
	Club
)

func (s Suit) String() string {
	switch s {
	case Heart:
		return "H"
	case Diamond:
		return "D"
	case Spade:
		return "S"
	case Club:
		return "C"
	default:
		return fmt.Sprintf("Unknown[%v]", int(s))
	}
}

// Render as suitable HTML Unicode
func (s Suit) HTML() string {
	switch s {
	case Heart:
		return `<span style="color:red">&#9829;</span>`
	case Diamond:
		return `<span style="color:red">&#9830;</span>`
	case Spade:
		return "&#9824;"
	case Club:
		return "&#9827;"
	default:
		return fmt.Sprintf("Unknown[%v]", int(s))
	}
}

type Rank int8

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	switch r {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return fmt.Sprintf("Unknown[%v]", int(r))
	}
}

type Card struct {
	Rank
	Suit
}

func (c Card) String() string {
	return fmt.Sprint(c.Rank.String(), c.Suit.String())
}

func (c Card) HTML() string {
	return fmt.Sprintf("%v%v", c.Rank.String(), c.Suit.HTML())
}

// Construct a card from text, e.g. "QD" for queen of diamonds
func MakeCard(c string) (Card, error) {
	re := regexp.MustCompile("^([0123456789AJQK]+)([CDHS])$")
	match := re.FindStringSubmatch(strings.ToUpper(c))
	if match == nil {
		return Card{}, errors.New(fmt.Sprintf("Illegally formatted card %q", c))
	}

	var rank Rank
	switch match[1] {
	case "2":
		rank = Two
	case "3":
		rank = Three
	case "4":
		rank = Four
	case "5":
		rank = Five
	case "6":
		rank = Six
	case "7":
		rank = Seven
	case "8":
		rank = Eight
	case "9":
		rank = Nine
	case "10":
		rank = Ten
	case "J":
		rank = Jack
	case "Q":
		rank = Queen
	case "K":
		rank = King
	case "A":
		rank = Ace
	}

	var suit Suit
	switch match[2] {
	case "C":
		suit = Club
	case "D":
		suit = Diamond
	case "H":
		suit = Heart
	case "S":
		suit = Spade
	}

	return Card{rank, suit}, nil
}

// Abbreviated constructor function for a card, to save typing.
func C(c string) Card {
	card, err := MakeCard(c)
	if err != nil {
		panic(err.Error())
	}
	return card
}

// Classification of a poker hand (e.g. "straight flush"). For any two
// hands with different classifications, the higher-classified hand
// will always beat the lower.
type HandClass int8

const (
	HighCard HandClass = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	MAX_HANDCLASS // Just a convenience value for iteration
)

func (hc HandClass) String() string {
	switch hc {
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	default:
		return fmt.Sprintf("Unknown (%v)", int(hc))
	}
}

// The full information needed to determine whether one hand beats another.
// Tiebreak ranks are used to determine which hand wins if both are in the same class.
// For example, two full houses are compared first by the three-card set, then by the two-card set,
// so e.g. for threes-over-Jacks, the tiebreaks should be {Three,Jack}.
type HandLevel struct {
	Class     HandClass
	Tiebreaks []Rank
}

// Function which returns a hand level beaten by any legitimate level
func MinLevel() HandLevel {
	return HandLevel{HighCard, []Rank{Two, Two, Two, Two, Two}}
}

func (hl HandLevel) PrettyTiebreaks() string {
	rankStrings := make([]string, len(hl.Tiebreaks))
	for i, r := range hl.Tiebreaks {
		rankStrings[i] = r.String()
	}
	return strings.Join(rankStrings, ", ")
}

func (hl HandLevel) String() string {
	return fmt.Sprintf("%v [%v]", hl.Class, hl.PrettyTiebreaks())
}

func (hl HandLevel) PrettyPrint() string {
	switch hl.Class {
	case StraightFlush:
		return fmt.Sprintf("Straight Flush: %v high", hl.Tiebreaks[0])
	case FourOfAKind:
		return fmt.Sprintf("Four %vs (plus %v)", hl.Tiebreaks[0], hl.Tiebreaks[1])
	case FullHouse:
		return fmt.Sprintf("Full house: %vs over %vs", hl.Tiebreaks[0], hl.Tiebreaks[1])
	case Flush:
		return fmt.Sprintf("Flush: %v", hl.PrettyTiebreaks())
	case Straight:
		return fmt.Sprintf("Straight: %v high", hl.Tiebreaks[0])
	case ThreeOfAKind:
		return fmt.Sprintf("Three %vs (plus %v, %v)", hl.Tiebreaks[0], hl.Tiebreaks[1], hl.Tiebreaks[2])
	case TwoPair:
		return fmt.Sprintf("Two pair: %vs and %vs (plus %v)", hl.Tiebreaks[0], hl.Tiebreaks[1], hl.Tiebreaks[2])
	case OnePair:
		return fmt.Sprintf("Pair %vs (plus %v, %v, %v)", hl.Tiebreaks[0], hl.Tiebreaks[1], hl.Tiebreaks[2], hl.Tiebreaks[3])
	case HighCard:
		return fmt.Sprintf("High card: %v", hl.PrettyTiebreaks())
	default:
		panic(fmt.Sprintf("Unknown class %v", hl.Class))
	}
}

// All the possible sets of ranks which make up straights, starting with the highest-value
var straights = make([][]Rank, 0)
var ranksDesc = []Rank{Ace, King, Queen, Jack, Ten, Nine, Eight, Seven, Six, Five, Four, Three, Two}
var allSuits = []Suit{Club, Diamond, Heart, Spade}

func init() {
	// Initialise the set of straights
	for high := Ace; high >= Six; high-- {
		straight := make([]Rank, 5)
		for i := 0; i < 5; i++ {
			straight[i] = Rank(int(high) - i)
		}
		straights = append(straights, straight)
	}
	// Ace-low is a special case
	straights = append(straights, []Rank{Five, Four, Three, Two, Ace})
}

// Compares two hands to see if one beats the other.
// All this needs to do is compare the levels; if they match, we compare the tiebreak ranks lexicographically.
func Beats(h1 HandLevel, h2 HandLevel) bool {
	if h1.Class != h2.Class {
		return h1.Class > h2.Class
	}
	for i := 0; i < len(h1.Tiebreaks) && i < len(h2.Tiebreaks); i++ {
		if h1.Tiebreaks[i] != h2.Tiebreaks[i] {
			return h1.Tiebreaks[i] > h2.Tiebreaks[i]
		}
	}
	return false
}

// Simple linear search
func containsAllCards(cards []Card, testSubset []Card) bool {
	for _, tc := range testSubset {
		found := false
		for _, c := range cards {
			if c == tc {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Implements the sort.Interface interface to allow cards to be sorted in descending rank order.
type CardSorter struct {
	Cards  []Card
	AceLow bool
}

func (cs CardSorter) Len() int {
	return len(cs.Cards)
}

func (cs CardSorter) Swap(i, j int) {
	cs.Cards[i], cs.Cards[j] = cs.Cards[j], cs.Cards[i]
}

// Cards are to be sorted in descending order of rank, so Less(c1,c2) is true iff c1 has a HIGHER rank than c2
func (cs CardSorter) Less(i, j int) bool {
	if cs.Cards[i].Rank != cs.Cards[j].Rank {
		if cs.AceLow {
			if cs.Cards[i].Rank == Ace {
				return false
			}
			if cs.Cards[j].Rank == Ace {
				return true
			}
		}
		return cs.Cards[i].Rank > cs.Cards[j].Rank
	}
	// Sort suits in an arbitrary way just to give a consistent ordering
	return cs.Cards[i].Suit < cs.Cards[j].Suit
}

var noLevel = HandLevel{OnePair, []Rank{}}

// See if this hand forms a straight flush. If so, return the level indicating that. HandLevel should be ignored if valid is false.
func classifyStraightFlush(cards []Card) (hl HandLevel, ok bool) {
	firstSuit := cards[0].Suit
	ranksForFirstSuit := make([]int, 13)
	for _, c := range cards {
		if c.Suit != firstSuit {
			return noLevel, false // Not all the same suit -> no straight flush
		}
		ranksForFirstSuit[c.Rank]++
	}

	for _, straight := range straights {
		foundStraightFlush := true
		for _, r := range straight {
			if ranksForFirstSuit[r] == 0 {
				foundStraightFlush = false
				break
			}
		}
		if foundStraightFlush {
			return HandLevel{StraightFlush, straight[0:1]}, true
		}
	}
	return noLevel, false
}

// See if "four of a kind" can be formed from this set of cards. If so, return the level of the best such hand; if not, indicate invalid.
func classifyFourOfAKind(cards []Card, countsByRank []int) (hl HandLevel, ok bool) {
	for _, r := range ranksDesc {
		if countsByRank[r] < 4 {
			continue
		}
		// Since the cards are sorted by rank, the other one will either be the first or the last
		otherRank := cards[0].Rank
		if otherRank == r {
			otherRank = cards[4].Rank
		}
		return HandLevel{FourOfAKind, []Rank{r, otherRank}}, true
	}
	return noLevel, false
}

// See if a full house can be composed from this set of cards. If so, return the level of the best such hand; otherwise indicate invalid.
func classifyFullHouse(cards []Card, countsByRank []int) (hl HandLevel, ok bool) {
	overUnderRanks := make([]Rank, 2)
	for _, r := range ranksDesc {
		switch countsByRank[r] {
		case 3:
			overUnderRanks[0] = r
		case 2:
			overUnderRanks[1] = r
		case 0:
			continue
		default:
			return noLevel, false
		}
	}
	return HandLevel{FullHouse, overUnderRanks}, true
}

// Can we make a flush out of these cards? If so, return the level; otherwise, indicate invalid.
func classifyFlush(cards []Card) (hl HandLevel, ok bool) {
	countsBySuit := make([]int, 4)
	for _, c := range cards {
		countsBySuit[c.Suit]++
	}

	for _, count := range countsBySuit {
		switch count {
		case 5:
			ranks := make([]Rank, len(cards))
			for i, c := range cards {
				ranks[i] = c.Rank
			}
			return HandLevel{Flush, ranks}, true
		case 0:
			continue
		default:
			return noLevel, false
		}
	}
	return noLevel, false // Should never hit this
}

// Can we make a straight out of these cards? If so, return the level; otherwise, indicate invalid.
func classifyStraight(cards []Card, countsByRank []int) (hl HandLevel, ok bool) {
	for _, straight := range straights {
		straightRealised := true
		for _, r := range straight {
			if countsByRank[r] != 1 {
				straightRealised = false
				break
			}
		}
		if straightRealised {
			return HandLevel{Straight, straight}, true
		}
	}
	return noLevel, false
}

// Can we build three-of-a-kind from this set of cards? If so, return the level, otherwise indicate failure.
func classifyThreeOfAKind(cards []Card, countsByRank []int) (hl HandLevel, ok bool) {
	for _, r := range ranksDesc {
		if countsByRank[r] < 3 {
			continue
		}
		handRanks := []Rank{r}
		for _, c := range cards {
			if c.Rank != r {
				handRanks = append(handRanks, c.Rank)
			}
		}
		return HandLevel{ThreeOfAKind, handRanks}, true
	}
	return noLevel, false
}

// Can we get two pairs out of this set of cards? If so, return the level, otherwise indicate failure.
func classifyTwoPair(cards []Card, countsByRank []int) (hl HandLevel, ok bool) {
	handRanks := make([]Rank, 3)
	found := make([]bool, 3)
	for _, r := range ranksDesc {
		switch countsByRank[r] {
		case 1:
			if found[2] {
				return noLevel, false // We have exactly one of more than one rank -> can't be two pair
			}
			handRanks[2] = r
			found[2] = true
		case 2:
			if found[0] {
				handRanks[1] = r
				found[1] = true
			} else {
				handRanks[0] = r
				found[0] = true
			}
		case 0:
			continue
		default:
			return noLevel, false
		}
	}
	return HandLevel{TwoPair, handRanks}, true
}

// Can we get a one-pair out of this set of cards? If so, return the level, otherwise indicate failure.
func classifyOnePair(cards []Card, countsByRank []int) (hl HandLevel, ok bool) {
	for _, r := range ranksDesc {
		if countsByRank[r] == 2 {
			handRanks := []Rank{r}
			for _, c := range cards {
				if c.Rank != r {
					handRanks = append(handRanks, c.Rank)
				}
			}
			return HandLevel{OnePair, handRanks}, true
		}
	}
	return noLevel, false
}

// If we have to call this function, we can't have anything better than a high-card, so rank it.
func classifyHighCard(cards []Card) HandLevel {
	handRanks := make([]Rank, len(cards))
	for i, c := range cards {
		handRanks[i] = c.Rank
	}
	return HandLevel{HighCard, handRanks}
}

// Classify a hand of five cards
func classifyHand(cards []Card) HandLevel {
	if len(cards) != 5 {
		panic(fmt.Sprintf("Expected exactly five cards, found %v", len(cards)))
	}
	// First sort the cards, as this makes some of the functions easier to write
	sort.Sort(CardSorter{cards, false})

	countsByRank := make([]int, 13)
	for _, c := range cards {
		countsByRank[c.Rank]++
	}

	if result, ok := classifyStraightFlush(cards); ok {
		return result
	}
	if result, ok := classifyFourOfAKind(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyFullHouse(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyFlush(cards); ok {
		return result
	}
	if result, ok := classifyStraight(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyThreeOfAKind(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyTwoPair(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyOnePair(cards, countsByRank); ok {
		return result
	}
	return classifyHighCard(cards)
}

// A sub-function with an extra argument startSkippingAt to avoid duplication of results
func allChoicesSkipping(cards []Card, num, startSkippingAt int) [][]Card {
	if num >= len(cards) {
		return [][]Card{cards}
	}

	result := [][]Card{}

	// Call the function recursively with every possible one-smaller combination, starting skipping at the appropriate location.
	for i := startSkippingAt; i < len(cards); i++ {
		nextSmaller := make([]Card, len(cards)-1)
		j := 0
		for k, c := range cards {
			if k == i {
				continue
			}
			nextSmaller[j] = c
			j++
		}
		subChoices := allChoicesSkipping(nextSmaller, num, i)
		if len(subChoices) > 0 {
			newResult := make([][]Card, len(result)+len(subChoices))
			copy(newResult[0:len(result)], result)
			copy(newResult[len(result):], subChoices)
			result = newResult
		}
	}

	return result
}

// All ways to choose n cards from a set
func allChoices(cards []Card, num int) [][]Card {
	return allChoicesSkipping(cards, num, 0)
}

// Classifies a poker hand composed of some mandatory cards (which MUST be in the constructed hand)
// and some optional cards (which MAY be used to construct the hand).
// For example, for Texas Hold'em, there will be two mandatory cards and five optional ones.
func Classify(mandatory, optional []Card) (HandLevel, []Card) {
	// Construct all possible hands and find the best one
	allPossibleOptionals := allChoices(optional, 5-len(mandatory))

	allPossibleHands := make([][]Card, len(allPossibleOptionals))
	for i, os := range allPossibleOptionals {
		hand := make([]Card, 5)
		copy(hand[0:len(mandatory)], mandatory)
		copy(hand[len(mandatory):], os)
		allPossibleHands[i] = hand
	}

	bestHand := allPossibleHands[0]
	bestRank := classifyHand(bestHand)

	for i := 1; i < len(allPossibleHands); i++ {
		hand := allPossibleHands[i]
		rank := classifyHand(hand)
		if Beats(rank, bestRank) {
			bestHand = hand
			bestRank = rank
		}
	}

	return bestRank, bestHand
}

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
package poker

import (
	"fmt"
	"strings"
)

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
var ranksDescAceLow = []Rank{King, Queen, Jack, Ten, Nine, Eight, Seven, Six, Five, Four, Three, Two, Ace}
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
			return HandLevel{Straight, straight[0:1]}, true
		}
	}
	return noLevel, false
}

// Can we build three-of-a-kind from this set of cards? If so, return the level, otherwise indicate failure.
func classifyThreeOfAKind(cards []Card, countsByRank []int, rankOrder []Rank) (hl HandLevel, ok bool) {
	for _, r := range rankOrder {
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
func classifyTwoPair(cards []Card, countsByRank []int, rankOrder []Rank) (hl HandLevel, ok bool) {
	handRanks := make([]Rank, 3)
	found := make([]bool, 3)
	for _, r := range rankOrder {
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
func classifyOnePair(cards []Card, countsByRank []int, rankOrder []Rank) (hl HandLevel, ok bool) {
	for _, r := range rankOrder {
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

func rankCounts(cards []Card) []int {
	countsByRank := make([]int, 13)
	for _, c := range cards {
		countsByRank[c.Rank]++
	}
	return countsByRank
}

// Classify a hand of five cards
func ClassifyHand(cards []Card) HandLevel {
	if len(cards) != 5 {
		panic(fmt.Sprintf("Expected exactly five cards, found %v", len(cards)))
	}
	// First sort the cards, as this makes some of the functions easier to write
	SortCards(cards, false)

	countsByRank := rankCounts(cards)
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
	if result, ok := classifyThreeOfAKind(cards, countsByRank, ranksDesc); ok {
		return result
	}
	if result, ok := classifyTwoPair(cards, countsByRank, ranksDesc); ok {
		return result
	}
	if result, ok := classifyOnePair(cards, countsByRank, ranksDesc); ok {
		return result
	}
	return classifyHighCard(cards)
}

// Ace-to-five low classification is very similar to standard classification - we just ignore straights and flushes.
func ClassifyAceToFiveLow(cards []Card) HandLevel {
	if len(cards) != 5 {
		panic(fmt.Sprintf("Expected exactly five cards, found %v", len(cards)))
	}
	// Sort ace-low
	SortCards(cards, true)

	countsByRank := rankCounts(cards)
	if result, ok := classifyFourOfAKind(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyFullHouse(cards, countsByRank); ok {
		return result
	}
	if result, ok := classifyThreeOfAKind(cards, countsByRank, ranksDescAceLow); ok {
		return result
	}
	if result, ok := classifyTwoPair(cards, countsByRank, ranksDescAceLow); ok {
		return result
	}
	if result, ok := classifyOnePair(cards, countsByRank, ranksDescAceLow); ok {
		return result
	}
	return classifyHighCard(cards)
}

// Determine whether one hand beats another using the ace-to-five low system.
func BeatsAceToFiveLow(l1, l2 HandLevel) bool {
	if l1.Class != l2.Class {
		return l1.Class < l2.Class
	}
	for i := 0; i < len(l1.Tiebreaks) && i < len(l2.Tiebreaks); i++ {
		if l1.Tiebreaks[i] != l2.Tiebreaks[i] {
			return IsRankLess(l1.Tiebreaks[i], l2.Tiebreaks[i], true)
		}
	}
	return false
}

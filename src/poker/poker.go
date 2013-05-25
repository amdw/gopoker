/*
 Code dealing with basic rules of poker: hand rankings, classifications and so forth.
*/
package poker

import (
	"fmt"
	"regexp"
	"sort"
)

type Suit int

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

	}
	panic(fmt.Sprintf("Illegal Suit value %s", s))
}

type Rank int

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
	}
	panic(fmt.Sprintf("Illegal Rank value %v", r))
}

type Card struct {
	Rank
	Suit
}

func (c Card) String() string {
	return fmt.Sprint(c.Rank.String(), c.Suit.String())
}

// Abbreviated constructor function for a card, to save typing.
func C(c string) Card {
	re := regexp.MustCompile("^([0123456789AJQK]+)([CDHS])$")
	match := re.FindStringSubmatch(c)
	if match == nil {
		panic(fmt.Sprintf("Illegally formatted card %q", c))
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

	return Card{rank, suit}
}

// Classification of a poker hand (e.g. "straight flush"). For any two
// hands with different classifications, the higher-classified hand
// will always beat the lower.
type HandClass int

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
	}
	panic(fmt.Sprintf("Illegal HandClass %v", hc))
}

// The full information needed to determine whether one hand beats another.
// Tiebreak ranks are used to determine which hand wins if both are in the same class.
// For example, two full houses are compared first by the three-card set, then by the two-card set,
// so e.g. for threes-over-Jacks, the tiebreaks should be {Three,Jack}.
// The actual cards in the hand are included for informational purposes only, and may or may not be populated.
type HandLevel struct {
	Class     HandClass
	Tiebreaks []Rank
	Cards     []Card
}

func (hl HandLevel) String() string {
	if len(hl.Cards) == 0 {
		return fmt.Sprintf("%v (%v)", hl.Class, hl.Tiebreaks)
	}
	return fmt.Sprintf("%v (%q)", hl.Class, hl.Cards)
}

// All the possible sets of ranks which make up straights, starting with the highest-value
var straights = make([][]Rank, 10)

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

func (cs CardSorter) Less(i, j int) bool {
	if cs.Cards[i].Rank != cs.Cards[j].Rank {
		if cs.AceLow {
			if cs.Cards[i].Rank == Ace {
				return true
			}
			if cs.Cards[j].Rank == Ace {
				return false
			}
		}
		return cs.Cards[i].Rank > cs.Cards[j].Rank
	}
	// Sort suits in an arbitrary way just to give a consistent ordering
	return cs.Cards[i].Suit < cs.Cards[j].Suit
}

var noLevel = HandLevel{OnePair, []Rank{}, []Card{}}

// See if this hand forms a straight flush. If so, return the level indicating that. HandLevel should be ignored if valid is false.
func classifyStraightFlush(mandatory, optional []Card) (hl HandLevel, valid bool) {
	// Mash all the cards together and see if there is a straight flush which contains the mandatory cards
	valuesBySuit := make([][]int, 4)
	for s := 0; s < 4; s++ {
		valuesBySuit[s] = make([]int, 13)
	}
	for _, c := range mandatory {
		valuesBySuit[c.Suit][c.Rank]++
	}
	for _, c := range optional {
		valuesBySuit[c.Suit][c.Rank]++
	}
	straightFlushes := make([][]Card, 0)
	for s := 0; s < 4; s++ {
		for _, straight := range straights {
			foundStraightFlush := true
			for _, r := range straight {
				if valuesBySuit[s][r] == 0 {
					foundStraightFlush = false
					break
				}
			}
			if foundStraightFlush {
				sf := make([]Card, 5)
				for i, r := range straight {
					sf[i] = Card{r, Suit(s)}
				}
				if containsAllCards(sf, mandatory) {
					straightFlushes = append(straightFlushes, sf)
					break // This is the best one for this suit
				}
			}
		}
	}
	if len(straightFlushes) == 0 {
		return noLevel, false
	}
	bestStraightFlush := straightFlushes[0]
	for i := 1; i < len(straightFlushes); i++ {
		if straightFlushes[i][0].Rank > bestStraightFlush[0].Rank {
			bestStraightFlush = straightFlushes[i]
		}
	}
	return HandLevel{StraightFlush, []Rank{bestStraightFlush[0].Rank}, bestStraightFlush}, true
}

// See if "four of a kind" can be formed from this set of cards. If so, return the level of the best such hand; if not, indicate invalid.
func classifyFourOfAKind(mandatory, optional []Card) (hl HandLevel, valid bool) {
	// Again, just mash everything together and find the hands matching the pattern which contain the mandatory cards
	ranks := make([]int, 13)
	for _, c := range mandatory {
		ranks[c.Rank]++
	}
	for _, c := range optional {
		ranks[c.Rank]++
	}
	for r := 12; r >= 0; r-- {
		if ranks[r] < 4 {
			continue
		}
		hand := make([]Card, 5)
		for i := 0; i < 4; i++ {
			hand[i] = Card{Rank(r), Suit(i)}
		}
		// If the mandatory cards hold something other than this rank, use the best of these;
		// otherwise, use the best of the other ranks from the optional set.
		findMissingCard := func(candidates []Card) (Card, bool) {
			for _, candidate := range candidates {
				if candidate.Rank == Rank(r) {
					continue
				}
				return candidate, true
			}
			return Card{}, false
		}
		for _, candidateList := range [][]Card{mandatory, optional} {
			if mc, ok := findMissingCard(candidateList); ok {
				hand[4] = mc
				if containsAllCards(hand, mandatory) {
					return HandLevel{FourOfAKind, []Rank{Rank(r), mc.Rank}, hand}, true
				}
			}
		}
	}
	return noLevel, false
}

// See if a full house can be composed from this set of cards. If so, return the level of the best such hand; otherwise indicate invalid.
func classifyFullHouse(mandatory, optional []Card) (hl HandLevel, valid bool) {
	// Once again, we find all full houses from the full set of cards, rejecting those which don't contain all the mandatory cards
	ranks := make([][]Suit, 13)
	for _, c := range mandatory {
		ranks[c.Rank] = append(ranks[c.Rank], c.Suit)
	}
	for _, c := range optional {
		ranks[c.Rank] = append(ranks[c.Rank], c.Suit)
	}
	for overRank := 12; overRank >= 0; overRank-- {
		if len(ranks[overRank]) < 3 {
			continue
		}
		for underRank := 12; underRank >= 0; underRank-- {
			if underRank == overRank || len(ranks[underRank]) < 2 {
				continue
			}
			hand := make([]Card, 5)
			for i := 0; i < 3; i++ {
				hand[i] = Card{Rank(overRank), ranks[overRank][i]}
			}
			for i := 0; i < 2; i++ {
				hand[i+3] = Card{Rank(underRank), ranks[underRank][i]}
			}
			if containsAllCards(hand, mandatory) {
				return HandLevel{FullHouse, []Rank{Rank(overRank), Rank(underRank)}, hand}, true
			}
		}
	}
	return noLevel, false
}

// Can we make a flush out of these cards? If so, return the level; otherwise, indicate invalid.
func classifyFlush(mandatory, optional []Card) (hl HandLevel, valid bool) {
	ranks := make([][]Rank, 4)
	for _, c := range mandatory {
		ranks[c.Suit] = append(ranks[c.Suit], c.Rank)
	}
	for _, c := range optional {
		ranks[c.Suit] = append(ranks[c.Suit], c.Rank)
	}

	flushes := make([][]Card, 0)
	for s := 0; s < 4; s++ {
		if len(ranks[s]) < 5 {
			continue // No flushes in this suit
		}
		allCardsForThisSuit := make([]Card, len(ranks[s]))
		for i, r := range ranks[s] {
			allCardsForThisSuit[i] = Card{r, Suit(s)}
		}
		if !containsAllCards(allCardsForThisSuit, mandatory) {
			continue // There's a mandatory card not in this suit
		}
		flush := make([]Card, 5)
		copy(flush, mandatory)
		i := len(mandatory)
		for _, c := range optional {
			if i >= 5 {
				break
			}
			if c.Suit == Suit(s) {
				flush[i] = c
				i++
			}
		}
		sort.Sort(CardSorter{flush, false})
		flushes = append(flushes, flush)
	}
	if len(flushes) == 0 {
		return noLevel, false
	}
	bestFlush := flushes[0]
	for i := 1; i < len(flushes); i++ {
		for j := 0; j < 5; j++ {
			if flushes[i][j].Rank == bestFlush[j].Rank {
				continue
			}
			if flushes[i][j].Rank > bestFlush[j].Rank {
				bestFlush = flushes[i]
			}
		}
	}
	handRanks := make([]Rank, 5)
	for i, c := range bestFlush {
		handRanks[i] = c.Rank
	}
	return HandLevel{Flush, handRanks, bestFlush}, true
}

// Recursive function to build all possible straights using a Cartesian product.
// Ranks are totally ignored in this function: they are assumed to be known by the caller.
// Instead, the input is the set of suits available at each rank.
// Example: {{S},{C,H},{S},{S},{S}} => {{S,C,S,S,S},{S,H,S,S,S}}.
// TODO: Optimise through earlier filtering of hands not containing mandatory cards?
func buildStraights(suits [][]Suit) [][]Suit {
	result := [][]Suit{{}}
	for _, suitsForRank := range suits {
		for _, suit := range suitsForRank {
			for i, h := range result {
				result[i] = append(h, suit)
			}
		}
	}
	return result
}

// Can we make a straight out of these cards? If so, return the level; otherwise, indicate invalid.
func classifyStraight(mandatory, optional []Card) (hl HandLevel, valid bool) {
	suits := make([][]Suit, 13)
	for _, c := range mandatory {
		suits[c.Rank] = append(suits[c.Rank], c.Suit)
	}
	for _, c := range optional {
		suits[c.Rank] = append(suits[c.Rank], c.Suit)
	}

	ourStraights := make([][]Card, 0)
	for _, possStraight := range straights {
		realised := true
		for _, rank := range possStraight {
			if len(suits[rank]) == 0 {
				realised = false
				break
			}
		}
		if !realised {
			continue
		}
		suitsForThisStraight := make([][]Suit, 5)
		for i, rank := range possStraight {
			suitsForThisStraight[i] = suits[rank]
		}
		suitSets := buildStraights(suitsForThisStraight)
		for _, suitSet := range suitSets {
			hand := make([]Card, 5)
			for i, suit := range suitSet {
				hand[i] = Card{possStraight[i], suit}
			}
			if containsAllCards(hand, mandatory) {
				ourStraights = append(ourStraights, hand)
			}
		}
	}

	if len(ourStraights) == 0 {
		return noLevel, false
	}
	bestStraight := ourStraights[0]
	for i := 1; i < len(ourStraights); i++ {
		for j := 0; j < 5; j++ {
			if ourStraights[i][j].Rank == bestStraight[j].Rank {
				continue
			}
			if ourStraights[i][j].Rank > bestStraight[j].Rank {
				bestStraight = ourStraights[i]
				break
			}
		}
	}
	ourRanks := make([]Rank, 5)
	for i, c := range bestStraight {
		ourRanks[i] = c.Rank
	}
	return HandLevel{Straight, ourRanks, bestStraight}, true
}

// Classifies a poker hand composed of some mandatory cards (which MUST be in the constructed hand)
// and some optional cards (which MAY be used to construct the hand).
// For example, for Texas Hold'em, there will be two mandatory cards and five optional ones.
func Classify(mandatory, optional []Card) HandLevel {
	// First sort the cards, as this makes some of the functions easier to write
	sort.Sort(CardSorter{mandatory, false})
	sort.Sort(CardSorter{optional, false})

	if result, ok := classifyStraightFlush(mandatory, optional); ok {
		return result
	}
	if result, ok := classifyFourOfAKind(mandatory, optional); ok {
		return result
	}
	if result, ok := classifyFullHouse(mandatory, optional); ok {
		return result
	}
	if result, ok := classifyFlush(mandatory, optional); ok {
		return result
	}
	if result, ok := classifyStraight(mandatory, optional); ok {
		return result
	}
	return noLevel
	//return HandLevel{StraightFlush, []Rank{Ace}}
}

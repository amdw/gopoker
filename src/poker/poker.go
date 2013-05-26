/*
 Code dealing with basic rules of poker: hand rankings, classifications and so forth.
*/
package poker

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
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
	}
	panic(fmt.Sprintf("Illegal suit value %v", s))
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

func (c Card) HTML() string {
	return fmt.Sprintf("%v%v", c.Rank.String(), c.Suit.HTML())
}

// Construct a card from text, e.g. "QD" for queen of diamonds
func MakeCard(c string) (Card, error) {
	re := regexp.MustCompile("^([0123456789AJQK]+)([CDHS])$")
	match := re.FindStringSubmatch(c)
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
	}
	panic(fmt.Sprintf("Illegal HandClass %v", hc))
}

// The full information needed to determine whether one hand beats another.
// Tiebreak ranks are used to determine which hand wins if both are in the same class.
// For example, two full houses are compared first by the three-card set, then by the two-card set,
// so e.g. for threes-over-Jacks, the tiebreaks should be {Three,Jack}.
type HandLevel struct {
	Class     HandClass
	Tiebreaks []Rank
}

func (hl HandLevel) String() string {
	rankStrings := make([]string, len(hl.Tiebreaks))
	for i, r := range hl.Tiebreaks {
		rankStrings[i] = r.String()
	}

	return fmt.Sprintf("%v [%v]", hl.Class, strings.Join(rankStrings, ","))
}

// All the possible sets of ranks which make up straights, starting with the highest-value
var straights = make([][]Rank, 10)
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

var noLevel = HandLevel{OnePair, []Rank{}}

// TODO: These classification functions should be optimised and refactored to make their methods more consistent and avoid repeated
// recalculation of the same data (e.g. ranks by suit, suits by rank, best remaining card).

// See if this hand forms a straight flush. If so, return the level indicating that. HandLevel should be ignored if valid is false.
func classifyStraightFlush(mandatory, optional []Card) (hl HandLevel, cards []Card, valid bool) {
	// Mash all the cards together and see if there is a straight flush which contains the mandatory cards
	valuesBySuit := make([][]int, 4)
	for _, s := range allSuits {
		valuesBySuit[s] = make([]int, 13)
	}
	for _, c := range mandatory {
		valuesBySuit[c.Suit][c.Rank]++
	}
	for _, c := range optional {
		valuesBySuit[c.Suit][c.Rank]++
	}
	straightFlushes := make([][]Card, 0)
	for _, s := range allSuits {
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
					sf[i] = Card{r, s}
				}
				if containsAllCards(sf, mandatory) {
					straightFlushes = append(straightFlushes, sf)
					break // This is the best one for this suit
				}
			}
		}
	}
	if len(straightFlushes) == 0 {
		return noLevel, []Card{}, false
	}
	bestStraightFlush := straightFlushes[0]
	for i := 1; i < len(straightFlushes); i++ {
		if straightFlushes[i][0].Rank > bestStraightFlush[0].Rank {
			bestStraightFlush = straightFlushes[i]
		}
	}
	return HandLevel{StraightFlush, []Rank{bestStraightFlush[0].Rank}}, bestStraightFlush, true
}

// See if "four of a kind" can be formed from this set of cards. If so, return the level of the best such hand; if not, indicate invalid.
func classifyFourOfAKind(mandatory, optional []Card) (hl HandLevel, cards []Card, valid bool) {
	// Again, just mash everything together and find the hands matching the pattern which contain the mandatory cards
	ranks := make([]int, 13)
	for _, c := range mandatory {
		ranks[c.Rank]++
	}
	for _, c := range optional {
		ranks[c.Rank]++
	}
	for _, r := range ranksDesc {
		if ranks[r] < 4 {
			continue
		}
		hand := make([]Card, 5)
		for _, i := range allSuits {
			hand[i] = Card{r, i}
		}
		// If the mandatory cards hold something other than this rank, use the best of these;
		// otherwise, use the best of the other ranks from the optional set.
		findMissingCard := func(candidates []Card) (Card, bool) {
			for _, candidate := range candidates {
				if candidate.Rank == r {
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
					return HandLevel{FourOfAKind, []Rank{r, mc.Rank}}, hand, true
				}
			}
		}
	}
	return noLevel, []Card{}, false
}

// See if a full house can be composed from this set of cards. If so, return the level of the best such hand; otherwise indicate invalid.
func classifyFullHouse(mandatory, optional []Card) (hl HandLevel, cards []Card, valid bool) {
	// Once again, we find all full houses from the full set of cards, rejecting those which don't contain all the mandatory cards
	ranks := make([][]Suit, 13)
	for _, c := range mandatory {
		ranks[c.Rank] = append(ranks[c.Rank], c.Suit)
	}
	for _, c := range optional {
		ranks[c.Rank] = append(ranks[c.Rank], c.Suit)
	}
	for _, overRank := range ranksDesc {
		if len(ranks[overRank]) < 3 {
			continue
		}
		for _, underRank := range ranksDesc {
			if underRank == overRank || len(ranks[underRank]) < 2 {
				continue
			}
			hand := make([]Card, 5)
			for i := 0; i < 3; i++ {
				hand[i] = Card{overRank, ranks[overRank][i]}
			}
			for i := 0; i < 2; i++ {
				hand[i+3] = Card{underRank, ranks[underRank][i]}
			}
			if containsAllCards(hand, mandatory) {
				return HandLevel{FullHouse, []Rank{overRank, underRank}}, hand, true
			}
		}
	}
	return noLevel, []Card{}, false
}

// Can we make a flush out of these cards? If so, return the level; otherwise, indicate invalid.
func classifyFlush(mandatory, optional []Card) (hl HandLevel, cards []Card, valid bool) {
	ranks := make([][]Rank, 4)
	for _, c := range mandatory {
		ranks[c.Suit] = append(ranks[c.Suit], c.Rank)
	}
	for _, c := range optional {
		ranks[c.Suit] = append(ranks[c.Suit], c.Rank)
	}

	flushes := make([][]Card, 0)
	for _, s := range allSuits {
		if len(ranks[s]) < 5 {
			continue // No flushes in this suit
		}
		allCardsForThisSuit := make([]Card, len(ranks[s]))
		for i, r := range ranks[s] {
			allCardsForThisSuit[i] = Card{r, s}
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
			if c.Suit == s {
				flush[i] = c
				i++
			}
		}
		sort.Sort(CardSorter{flush, false})
		flushes = append(flushes, flush)
	}
	if len(flushes) == 0 {
		return noLevel, []Card{}, false
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
	return HandLevel{Flush, handRanks}, bestFlush, true
}

// Build all possible straights using a Cartesian product.
// Ranks are totally ignored in this function: they are assumed to be known by the caller.
// Instead, the input is the set of suits available at each rank.
// Example: {{S},{C,H},{S},{S},{S}} => {{S,C,S,S,S},{S,H,S,S,S}}.
// TODO: Optimise through earlier filtering of hands not containing mandatory cards?
func buildStraights(suits [][]Suit) [][]Suit {
	result := [][]Suit{{}}
	for _, suitsForRank := range suits {
		newResult := [][]Suit{}
		for _, suit := range suitsForRank {
			for _, h := range result {
				newResult = append(newResult, append(h, suit))
			}
		}
		result = newResult
	}
	return result
}

// Can we make a straight out of these cards? If so, return the level; otherwise, indicate invalid.
func classifyStraight(mandatory, optional []Card) (hl HandLevel, cards []Card, valid bool) {
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
		return noLevel, []Card{}, false
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
	return HandLevel{Straight, ourRanks}, bestStraight, true
}

// Can we build three-of-a-kind from this set of cards? If so, return the level, otherwise indicate failure.
func classifyThreeOfAKind(mandatory, optional []Card) (hl HandLevel, cards []Card, ok bool) {
	ranks := make([]int, 13)
	for _, c := range mandatory {
		ranks[c.Rank]++
	}
	for _, c := range optional {
		ranks[c.Rank]++
	}
	for _, r := range ranksDesc {
		if ranks[r] < 3 {
			continue
		}
		hand := make([]Card, 5)
		copy(hand, mandatory)
		missingTripleCards := 3
		for _, c := range mandatory {
			if c.Rank == r {
				missingTripleCards--
			}
		}
		i := len(mandatory)
		for _, c := range optional {
			if i >= 5 || missingTripleCards == 0 {
				break
			}
			if c.Rank == r {
				hand[i] = c
				i++
				missingTripleCards--
			}
		}
		if missingTripleCards > 0 {
			continue // We can't get all the cards from the triple into our hand
		}
		// Fill in with the best available remaining cards
		for _, c := range optional {
			if i >= 5 {
				break
			}
			if c.Rank == r {
				continue
			}
			hand[i] = c
			i++
		}
		if i < 5 {
			panic(fmt.Sprintf("Not enough cards: %q %q", mandatory, optional))
		}
		sort.Sort(CardSorter{hand, false})
		handRanks := []Rank{r}
		for _, c := range hand {
			if c.Rank != r {
				handRanks = append(handRanks, c.Rank)
			}
		}
		return HandLevel{ThreeOfAKind, handRanks}, hand, true
	}
	return noLevel, []Card{}, false
}

// Can we get two pairs out of this set of cards? If so, return the level, otherwise indicate failure.
func classifyTwoPair(mandatory, optional []Card) (hl HandLevel, cards []Card, ok bool) {
	suitsByRank := make([][]Suit, 13)
	mandatorySuitsByRank := make([][]Suit, 13)
	for _, c := range mandatory {
		suitsByRank[c.Rank] = append(suitsByRank[c.Rank], c.Suit)
		mandatorySuitsByRank[c.Rank] = append(mandatorySuitsByRank[c.Rank], c.Suit)
	}
	for _, c := range optional {
		suitsByRank[c.Rank] = append(suitsByRank[c.Rank], c.Suit)
	}
	ranksWithPairs := make([]Rank, 0)
	for _, r := range ranksDesc {
		if len(suitsByRank[r]) >= 2 {
			ranksWithPairs = append(ranksWithPairs, r)
		}
	}
	getBestOtherCard := func(cards []Card, pr1, pr2 Rank) (Card, bool) {
		for _, c := range cards {
			if c.Rank != pr1 && c.Rank != pr2 {
				return c, true
			}
		}
		return Card{}, false
	}
	getPairCards := func(rank Rank) []Card {
		result := []Card{}
		for _, c := range mandatory {
			if c.Rank == rank {
				result = append(result, c)
			}
		}
		for _, c := range optional {
			if len(result) == 2 {
				break
			}
			if c.Rank == rank {
				result = append(result, c)
			}
		}
		return result
	}
	for i, pairRank1 := range ranksWithPairs {
		for j := i + 1; j < len(ranksWithPairs); j++ {
			pairRank2 := ranksWithPairs[j]
			hand := make([]Card, 5)
			var bestOtherCard Card
			if boc, ok := getBestOtherCard(mandatory, pairRank1, pairRank2); ok {
				bestOtherCard = boc
			} else if boc, ok := getBestOtherCard(optional, pairRank1, pairRank2); ok {
				bestOtherCard = boc
			} else {
				panic(fmt.Sprintf("Not enough cards: %q %q", mandatory, optional))
			}
			copy(hand[0:2], getPairCards(pairRank1))
			copy(hand[2:4], getPairCards(pairRank2))
			hand[4] = bestOtherCard
			if containsAllCards(hand, mandatory) {
				return HandLevel{TwoPair, []Rank{pairRank1, pairRank2, bestOtherCard.Rank}}, hand, true
			}
		}
	}
	return noLevel, []Card{}, false
}

func classifyOnePair(mandatory, optional []Card) (hl HandLevel, cards []Card, ok bool) {
	ranks := make([]int, 13)
	for _, c := range mandatory {
		ranks[c.Rank]++
	}
	for _, c := range optional {
		ranks[c.Rank]++
	}
	for _, r := range ranksDesc {
		if ranks[r] < 2 {
			continue
		}
		hand := []Card{}
		pairCardsFound := 0
		for _, c := range mandatory {
			hand = append(hand, c)
			if c.Rank == r {
				pairCardsFound--
			}
		}
		for _, c := range optional {
			if pairCardsFound == 2 {
				break
			}
			if c.Rank == r {
				hand = append(hand, c)
			}
		}
		for _, c := range optional {
			if len(hand) >= 5 {
				break
			}
			if c.Rank == r {
				continue
			}
			hand = append(hand, c)
		}
		sort.Sort(CardSorter{hand, false})
		handRanks := []Rank{r}
		for _, c := range hand {
			if c.Rank != r {
				handRanks = append(handRanks, c.Rank)
			}
		}
		return HandLevel{OnePair, handRanks}, hand, true
	}
	return noLevel, []Card{}, false
}

func classifyHighCard(mandatory, optional []Card) (HandLevel, []Card) {
	hand := make([]Card, 5)
	copy(hand, mandatory)
	for i := 0; i < 5-len(mandatory); i++ {
		hand[i+len(mandatory)] = optional[i]
	}
	sort.Sort(CardSorter{hand, false})
	handRanks := make([]Rank, 5)
	for i, c := range hand {
		handRanks[i] = c.Rank
	}
	return HandLevel{HighCard, handRanks}, hand
}

// Classifies a poker hand composed of some mandatory cards (which MUST be in the constructed hand)
// and some optional cards (which MAY be used to construct the hand).
// For example, for Texas Hold'em, there will be two mandatory cards and five optional ones.
func Classify(mandatory, optional []Card) (HandLevel, []Card) {
	// First sort the cards, as this makes some of the functions easier to write
	sort.Sort(CardSorter{mandatory, false})
	sort.Sort(CardSorter{optional, false})

	if result, cards, ok := classifyStraightFlush(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyFourOfAKind(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyFullHouse(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyFlush(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyStraight(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyThreeOfAKind(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyTwoPair(mandatory, optional); ok {
		return result, cards
	}
	if result, cards, ok := classifyOnePair(mandatory, optional); ok {
		return result, cards
	}
	return classifyHighCard(mandatory, optional)
}

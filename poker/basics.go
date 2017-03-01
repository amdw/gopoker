/*
Copyright 2017 Andrew Medworth

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

// Construct a rank from text
func MakeRank(r string) (Rank, error) {
	var rank Rank
	switch r {
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
	default:
		return 0, errors.New(fmt.Sprintf("Illegally formatted rank '%q'", r))
	}
	return rank, nil
}

// Construct a card from text, e.g. "QD" for queen of diamonds
func MakeCard(c string) (Card, error) {
	re := regexp.MustCompile("^([0123456789AJQK]+)([CDHS])$")
	match := re.FindStringSubmatch(strings.ToUpper(c))
	if match == nil {
		return Card{}, errors.New(fmt.Sprintf("Illegally formatted card %q", c))
	}

	rank, err := MakeRank(match[1])
	if err != nil {
		return Card{}, err
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

// Simple linear search
func containsAllCards(cards, testSubset []Card) bool {
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

func SortCards(cards []Card, aceLow bool) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Rank != cards[j].Rank {
			if aceLow {
				if cards[i].Rank == Ace {
					return false
				}
				if cards[j].Rank == Ace {
					return true
				}
			}
			return cards[i].Rank > cards[j].Rank
		}
		// Sort suits in an arbitrary way just to give a consistent ordering
		return cards[i].Suit < cards[j].Suit
	})
}

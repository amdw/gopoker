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
package poker_http

import (
	"errors"
	"fmt"
	"github.com/amdw/gopoker/poker"
	"net/http"
	"strconv"
	"strings"
)

func duplicateCheck(tableCards, yourCards []poker.Card) (ok bool, dupeCard poker.Card) {
	allCards := make([]poker.Card, len(yourCards)+len(tableCards))
	copy(allCards, yourCards)
	copy(allCards[len(yourCards):], tableCards)
	cardDupeCheck := make(map[string]int)
	for _, c := range allCards {
		cardDupeCheck[c.String()]++
		if cardDupeCheck[c.String()] > 1 {
			return false, c
		}
	}
	return true, poker.Card{}
}

type simulationParams struct {
	players               int
	tableCards, yourCards []poker.Card
	handsToPlay           int
	forceComputation      bool
}

func getSimulationParams(req *http.Request) (params simulationParams, err error) {
	players, err := getPlayers(req)
	if err != nil {
		return simulationParams{}, errors.New(fmt.Sprintf("Could not get player count: %v", err))
	}

	params = simulationParams{players, []poker.Card{}, []poker.Card{}, 10000, false}

	if forceStrs, ok := req.Form[forceComputeKey]; ok && len(forceStrs) == 1 && strings.EqualFold(forceStrs[0], "true") {
		params.forceComputation = true
	}

	extractCards := func(key string) ([]poker.Card, error) {
		cards := []poker.Card{}
		if cardsStrs, ok := req.Form[key]; ok && len(cardsStrs) > 0 && len(cardsStrs[0]) > 0 {
			cardsSplit := strings.Split(strings.Replace(cardsStrs[0], " ", "", -1), ",")
			cards = make([]poker.Card, len(cardsSplit))
			for i, cstr := range cardsSplit {
				card, err := poker.MakeCard(cstr)
				if err != nil {
					return cards, errors.New(fmt.Sprintf("Illegally formatted card %q", cstr))
				}
				cards[i] = card
			}
		}
		return cards, nil
	}
	params.yourCards, err = extractCards(yourCardsKey)
	if err != nil {
		return params, err
	}
	if len(params.yourCards) > 2 {
		return params, errors.New(fmt.Sprintf("Maximum of 2 player cards allowed, found %v", len(params.yourCards)))
	}
	params.tableCards, err = extractCards(tableCardsKey)
	if err != nil {
		return params, err
	}
	if len(params.tableCards) > 5 {
		return params, errors.New(fmt.Sprintf("Maximum of 5 table cards allowed, found %v", len(params.tableCards)))
	}
	// Check for duplicate cards
	if ok, dupeCard := duplicateCheck(params.tableCards, params.yourCards); !ok {
		return params, errors.New(fmt.Sprintf("Found duplicate card %v in specification", dupeCard))
	}

	if handsToPlayStrs, ok := req.Form[simCountKey]; ok && len(handsToPlayStrs) > 0 {
		handsToPlayParsed, err := strconv.ParseInt(handsToPlayStrs[0], 10, 32)
		if err != nil {
			return params, errors.New(fmt.Sprintf("Could not parse simcount: %v", err.Error()))
		}
		params.handsToPlay = int(handsToPlayParsed)
	}

	return params, nil
}

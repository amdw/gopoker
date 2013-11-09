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
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"poker"
	"strconv"
	"strings"
)

func menu(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<html><head><title>Poker</title></head><body><h1>Poker</h1><ul>")
	fmt.Fprintf(w, `<li><a href="/play">Play</a></li>`)
	fmt.Fprintf(w, `<li><a href="/simulate">Simulate</a></li>`)
	fmt.Fprintf(w, "</ul></body></html>")
}

func formatCards(cards []poker.Card) string {
	cardStrings := make([]string, len(cards))
	for i, c := range cards {
		cardStrings[i] = c.HTML()
	}
	return strings.Join(cardStrings, ", ")
}

const playersKey = "players"

func getPlayers(req *http.Request) (int, error) {
	players := 5
	if plstrs, ok := req.Form[playersKey]; ok && len(plstrs) > 0 {
		pl, err := strconv.ParseInt(plstrs[0], 10, 32)
		if err != nil {
			return 0, err
		}
		players = int(pl)
	}
	return players, nil
}

func playHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintf(w, "<html><head><title>A game of Texas Hold'em</title></head><body><h1>A game of Texas Hold'em</h1>")
	players, err := getPlayers(req)
	if err != nil {
		// Use template to sanitise user input for security
		t := template.Must(template.New("error").Parse("<p>Could not parse players as integer: {{.}}</p></body></html>"))
		t.Execute(w, err.Error())
		return
	}

	fmt.Fprintf(w, `<form method="get">Players: <input type="text" name="%v" value="%v"/><input type="submit" value="Rerun"/></form>`, playersKey, players)

	pack := poker.NewPack()
	pack.Shuffle()
	onTable, playerCards, sortedOutcomes := pack.PlayHoldem(players)
	fmt.Fprintf(w, "<h2>Table cards</h2><p>%v</p>", formatCards(onTable))
	fmt.Fprintf(w, "<h2>Player cards</h2><ul>")
	for player := 0; player < players; player++ {
		fmt.Fprintf(w, "<li>Player %v: %v</li>", player+1, formatCards(playerCards[player]))
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "<h2>Results</h2><table><tr><th>Position</th><th>Player</th><th>Hand</th><th>Cards</th></tr>")
	for i, outcome := range sortedOutcomes.Outcomes {
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>", i+1, outcome.Player, outcome.Level, formatCards(outcome.Cards))
	}
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</body></html>")
}

const yourCardsKey = "yours"
const tableCardsKey = "table"
const simCountKey = "simcount"

func duplicateCheck(yourCards, tableCards []poker.Card) (ok bool, dupeCard poker.Card) {
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

func simulationParams(req *http.Request) (yourCards, tableCards []poker.Card, handsToPlay int, err error) {
	yourCards = []poker.Card{}
	tableCards = []poker.Card{}
	handsToPlay = 10000

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
	yourCards, err = extractCards(yourCardsKey)
	if err != nil {
		return yourCards, tableCards, handsToPlay, err
	}
	if len(yourCards) > 2 {
		return yourCards, tableCards, handsToPlay, errors.New(fmt.Sprintf("Maximum of 2 player cards allowed, found %v", len(yourCards)))
	}
	tableCards, err = extractCards(tableCardsKey)
	if err != nil {
		return yourCards, tableCards, handsToPlay, err
	}
	if len(tableCards) > 5 {
		return yourCards, tableCards, handsToPlay, errors.New(fmt.Sprintf("Maximum of 5 table cards allowed, found %v", len(tableCards)))
	}
	// Check for duplicate cards
	if ok, dupeCard := duplicateCheck(yourCards, tableCards); !ok {
		return yourCards, tableCards, handsToPlay, errors.New(fmt.Sprintf("Found duplicate card %v in specification", dupeCard))
	}

	if handsToPlayStrs, ok := req.Form[simCountKey]; ok && len(handsToPlayStrs) > 0 {
		handsToPlayParsed, err := strconv.ParseInt(handsToPlayStrs[0], 10, 32)
		if err != nil {
			return yourCards, tableCards, handsToPlay, errors.New(fmt.Sprintf("Could not parse simcount: %v", err.Error()))
		}
		handsToPlay = int(handsToPlayParsed)
	}

	return yourCards, tableCards, handsToPlay, nil
}

func simulateHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintf(w, "<html><head><title>Texas Hold'em simulator</title></head><body><h1>Texas Hold'em Simulator</h1>")

	players, err := getPlayers(req)
	if err != nil {
		// Use a template for security as error messages will often contain raw user input
		t := template.Must(template.New("error").Parse("<p>Could not get player count: {{.}}</p></body></html>"))
		t.Execute(w, err.Error())
		return
	}

	yourCards, tableCards, handsToPlay, err := simulationParams(req)
	if err != nil {
		t := template.Must(template.New("error").Parse("<p>Could not get simulation parameters: {{.}}</p></body></html>"))
		t.Execute(w, err.Error())
		return
	}

	simulator := poker.Simulator{}
	simulator.SimulateHoldem(yourCards, tableCards, players, handsToPlay)
	fmt.Fprintf(w, "<h2>Simulation outcome</h2>")
	fmt.Fprintf(w, `<form method="get">`)
	fmt.Fprintf(w, `<p><input type="submit" value="Rerun"/> <a href="/simulate">Reset</a></p>`)
	fmt.Fprintf(w, "<table>")
	cardText := func(cards []poker.Card) string {
		text := make([]string, len(cards))
		for i, c := range cards {
			text[i] = c.String()
		}
		return strings.Join(text, ",")
	}
	fmt.Fprintf(w, `<tr><td><b>Players</b></td><td><input type="text" name="%v" value="%v"/></td></tr>`, playersKey, players)
	fmt.Fprintf(w, `<tr><td><b>Your cards</b><br/><i>(comma-separated, e.g. 'KD,10H')</i></td><td>%v <input type="text" name="%v" value="%v"/></td></tr>`, formatCards(yourCards), yourCardsKey, cardText(yourCards))
	fmt.Fprintf(w, `<tr><td><b>Table cards</b></td><td>%v <input type="text" name="%v" value="%v"/></td></tr>`, formatCards(tableCards), tableCardsKey, cardText(tableCards))
	fmt.Fprintf(w, `<tr><td><b>Simulations</b></td><td><input type="text" name="%v" value="%v"/></td></tr>`, simCountKey, simulator.HandCount)
	fmt.Fprintf(w, "<tr><td><b>Wins</b></td><td>%v (%.1f%%)", simulator.WinCount, (float32(simulator.WinCount)*100.0)/float32(simulator.HandCount))
	fmt.Fprintf(w, "<tr><td><b>Your hand</b></td><td><ul>")
	printClassCounts := func(counts []int) {
		fmt.Fprintf(w, "<table>")
		for class, freq := range counts {
			fmt.Fprintf(w, `<tr><td>%v</td><td style="text-align:right">%v</td><td style="text-align:right">%.1f%%</td></tr>`, poker.HandClass(class).String(), freq, (float32(freq) * 100.0 / float32(simulator.HandCount)))
		}
		fmt.Fprintf(w, "</table>")
	}
	printClassCounts(simulator.OurClassCounts)
	fmt.Fprintf(w, "</ul></td></tr>")
	fmt.Fprintf(w, "<tr><td><b>Opponent's hand</b></td><td><ul>")
	printClassCounts(simulator.OpponentClassCounts)
	fmt.Fprintf(w, "</ul></td></tr>")
	fmt.Fprintf(w, "</table></form></body></html>")
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Listen port for HTTP server")
	flag.Parse()

	log.Printf("Listening on port %v...\n", port)

	http.HandleFunc("/", menu)
	http.HandleFunc("/play", playHoldem)
	http.HandleFunc("/simulate", simulateHoldem)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

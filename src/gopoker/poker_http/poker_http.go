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
	"gopoker/poker"
	"html/template"
	"log"
	"net/http"
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

func formatPct(num, denom int) string {
	result := ""
	if denom != 0 {
		result = fmt.Sprintf("%.1f%%", float32(num)*100.0/float32(denom))
	}
	return result
}

func printResultTable(w http.ResponseWriter, simulator poker.Simulator) {
	cssClass := func(isNum, isZero bool) string {
		result := "countTable"
		if isNum {
			result += " numcell"
		}
		if isZero {
			result += " zero"
		}
		return result
	}
	printStringCell := func(content string) {
		fmt.Fprintf(w, `<td class="%v">%v</td>`, cssClass(false, false), content)
	}
	printNumCell := func(content int) {
		fmt.Fprintf(w, `<td class="%v">%v</td>`, cssClass(true, content == 0), content)
	}
	printPctCell := func(num, denom int) {
		fmt.Fprintf(w, `<td class="%v">%v</td>`, cssClass(true, num == 0), formatPct(num, denom))
	}
	printRow := func(handClass string, classFreq, winCount, oppCount, oppWinCount int, isSummary bool) {
		cssClass := ""
		if isSummary {
			cssClass = ` class="summary"`
		}
		fmt.Fprintf(w, `<tr%v>`, cssClass)
		printStringCell(handClass)
		printNumCell(classFreq)
		printPctCell(classFreq, simulator.HandCount)
		printNumCell(winCount)
		printPctCell(winCount, simulator.HandCount)
		printNumCell(oppCount)
		printPctCell(oppCount, simulator.HandCount)
		printNumCell(oppWinCount)
		printPctCell(oppWinCount, simulator.HandCount)
		fmt.Fprintf(w, "</tr>")
	}
	fmt.Fprintf(w, `<table class="countTable"><tr><th class="countTable" rowspan="2">Hand</th><th class="countTable" colspan="4">For you</th><th class="countTable" colspan="4">For opponent</th></tr>`)
	fmt.Fprintf(w, `<tr><th class="countTable" colspan="2">Freq</th><th class="countTable" colspan="2">Wins</th><th class="countTable" colspan="2">Freq</th><th class="countTable" colspan="2">Wins</th></tr>`)
	printRow("All", simulator.HandCount, simulator.WinCount, simulator.HandCount, simulator.HandCount-simulator.WinCount, true)
	for class := range simulator.OurClassCounts {
		printRow(poker.HandClass(class).String(), simulator.OurClassCounts[class], simulator.ClassWinCounts[class], simulator.OpponentClassCounts[class], simulator.ClassOppWinCounts[class], false)
	}
	fmt.Fprintf(w, "</table>")
}

func simulateHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintf(w, "<html><head><title>Texas Hold'em simulator</title>")
	fmt.Fprintf(w, "<style>")
	fmt.Fprintf(w, "td.formcell { vertical-align: top }\n")
	fmt.Fprintf(w, "table.countTable { border-collapse: collapse; }\n")
	fmt.Fprintf(w, "th.countTable, td.countTable { border: 1px solid black; padding: 3px; }\n")
	fmt.Fprintf(w, "td.numcell { text-align: right }\n")
	fmt.Fprintf(w, "td.zero { color: lightgrey }\n")
	fmt.Fprintf(w, ".summary { font-weight: bold }\n")
	fmt.Fprintf(w, "</style>")
	fmt.Fprintf(w, "</head><body><h1>Texas Hold'em Simulator</h1>")

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
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Players</b></td><td><input type="text" name="%v" value="%v"/></td></tr>`, playersKey, players)
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Your cards</b><br/><i>(comma-separated, e.g. 'KD,10H')</i></td><td>%v <input type="text" name="%v" value="%v"/></td></tr>`, formatCards(yourCards), yourCardsKey, cardText(yourCards))
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Table cards</b></td><td>%v <input type="text" name="%v" value="%v"/></td></tr>`, formatCards(tableCards), tableCardsKey, cardText(tableCards))
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Simulations</b></td><td><input type="text" name="%v" value="%v"/></td></tr>`, simCountKey, simulator.HandCount)
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Results</b></td><td>`)
	printResultTable(w, simulator)
	fmt.Fprintf(w, "</td></tr>")
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

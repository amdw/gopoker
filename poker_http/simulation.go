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
package poker_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amdw/gopoker/poker"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const yourCardsKey = "yours"
const tableCardsKey = "table"
const simCountKey = "simcount"

func summariseCards(cards []poker.Card) string {
	if len(cards) == 0 {
		return "empty"
	}
	return formatCards(cards)
}

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

func simulationParams(req *http.Request) (tableCards, yourCards []poker.Card, handsToPlay int, err error) {
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
		return tableCards, yourCards, handsToPlay, err
	}
	if len(yourCards) > 2 {
		return tableCards, yourCards, handsToPlay, errors.New(fmt.Sprintf("Maximum of 2 player cards allowed, found %v", len(yourCards)))
	}
	tableCards, err = extractCards(tableCardsKey)
	if err != nil {
		return tableCards, yourCards, handsToPlay, err
	}
	if len(tableCards) > 5 {
		return tableCards, yourCards, handsToPlay, errors.New(fmt.Sprintf("Maximum of 5 table cards allowed, found %v", len(tableCards)))
	}
	// Check for duplicate cards
	if ok, dupeCard := duplicateCheck(tableCards, yourCards); !ok {
		return tableCards, yourCards, handsToPlay, errors.New(fmt.Sprintf("Found duplicate card %v in specification", dupeCard))
	}

	if handsToPlayStrs, ok := req.Form[simCountKey]; ok && len(handsToPlayStrs) > 0 {
		handsToPlayParsed, err := strconv.ParseInt(handsToPlayStrs[0], 10, 32)
		if err != nil {
			return tableCards, yourCards, handsToPlay, errors.New(fmt.Sprintf("Could not parse simcount: %v", err.Error()))
		}
		handsToPlay = int(handsToPlayParsed)
	}

	return tableCards, yourCards, handsToPlay, nil
}

func printResultGraph(w http.ResponseWriter, title string, handNames []string, series []map[string]interface{}, id string) {
	graphDef := map[string]interface{}{
		"chart":       map[string]string{"type": "column"},
		"title":       map[string]interface{}{"text": title, "useHTML": true},
		"xAxis":       map[string]interface{}{"categories": handNames},
		"yAxis":       map[string]interface{}{"title": map[string]string{"text": "Probability (%)"}, "min": 0, "max": 100},
		"series":      series,
		"plotOptions": map[string]interface{}{"series": map[string]string{"stacking": "normal"}},
		"tooltip":     map[string]string{"pointFormat": "{series.name}: <b>{point.y:.1f}%</b>"},
	}
	fmt.Fprintf(w, `<div id="%s" style="height: 400px"></div>`, id)
	fmt.Fprintln(w, `<script>`)
	fmt.Fprintf(w, `$(function () { $('#%s').highcharts(`, id)
	json.NewEncoder(w).Encode(graphDef)
	fmt.Fprintln(w, `)});`)
	fmt.Fprintln(w, `</script>`)
}

func makeYourSeries(simulator poker.Simulator) ([]string, []map[string]interface{}) {
	handNames := make([]string, len(simulator.ClassWinCounts))
	soleWinData := make([]float64, len(handNames))
	jointWinData := make([]float64, len(handNames))
	lossData := make([]float64, len(handNames))
	for class := range simulator.ClassWinCounts {
		handNames[class] = poker.HandClass(class).String()
		soleWinData[class] = 100.0 * float64(simulator.ClassWinCounts[class]-simulator.ClassJointWinCounts[class]) / float64(simulator.HandCount)
		jointWinData[class] = 100.0 * float64(simulator.ClassJointWinCounts[class]) / float64(simulator.HandCount)
		lossData[class] = 100.0 * float64(simulator.OurClassCounts[class]-simulator.ClassWinCounts[class]) / float64(simulator.HandCount)
	}
	overallData := []interface{}{
		map[string]interface{}{"name": "Sole winner", "y": 100.0 * float64(simulator.WinCount-simulator.JointWinCount) / float64(simulator.HandCount)},
		map[string]interface{}{"name": "Joint winner", "y": 100.0 * float64(simulator.JointWinCount) / float64(simulator.HandCount)},
		map[string]interface{}{"name": "Loser", "y": 100.0 * float64(simulator.HandCount-simulator.WinCount) / float64(simulator.HandCount)},
	}
	//"type": "column",
	series := []map[string]interface{}{
		map[string]interface{}{"name": "Sole winner", "data": soleWinData},
		map[string]interface{}{"name": "Joint winner", "data": jointWinData},
		map[string]interface{}{"name": "Loser", "data": lossData},
		map[string]interface{}{"name": "Overall", "type": "pie", "data": overallData, "size": 100, "center": []string{"60%", "25%"}, "showInLegend": false, "dataLabels": map[string]interface{}{"enabled": true, "format": "{point.name} {y:.1f}%"}},
	}
	return handNames, series
}

func makeBestOppSeries(simulator poker.Simulator) []map[string]interface{} {
	winData := make([]float64, len(simulator.ClassBestOppWinCounts))
	lossData := make([]float64, len(winData))
	for class := range simulator.ClassBestOppWinCounts {
		winData[class] = 100.0 * float64(simulator.ClassBestOppWinCounts[class]) / float64(simulator.HandCount)
		lossData[class] = 100.0 * float64(simulator.BestOpponentClassCounts[class]-simulator.ClassBestOppWinCounts[class]) / float64(simulator.HandCount)
	}
	overallData := []interface{}{
		map[string]interface{}{"name": "Winner", "y": 100.0 * float64(simulator.BestOpponentWinCount) / float64(simulator.HandCount)},
		map[string]interface{}{"name": "Loser", "y": 100.0 * float64(simulator.HandCount-simulator.BestOpponentWinCount) / float64(simulator.HandCount)},
	}
	series := []map[string]interface{}{
		map[string]interface{}{"name": "Winner (sole or joint)", "data": winData},
		map[string]interface{}{"name": "Loser", "data": lossData},
		map[string]interface{}{"name": "Overall", "type": "pie", "data": overallData, "size": 100, "center": []string{"60%", "25%"}, "showInLegend": false, "dataLabels": map[string]interface{}{"enabled": true, "format": "{point.name} {y:.1f}%"}},
	}
	return series
}

func printResultGraphs(w http.ResponseWriter, simulator poker.Simulator, tableCards, yourCards []poker.Card) {
	fmt.Fprintln(w, `<div class="row">`)

	fmt.Fprintln(w, `<div class="col-md-6">`)
	handNames, yourSeries := makeYourSeries(simulator)
	printResultGraph(w, "Your outcomes", handNames, yourSeries, "wingraph")
	fmt.Fprintln(w, `</div>`)

	fmt.Fprintln(w, `<div class="col-md-6">`)
	bestOppSeries := makeBestOppSeries(simulator)
	printResultGraph(w, "Best opponent outcomes", handNames, bestOppSeries, "bestoppwingraph")
	fmt.Fprintln(w, `</div>`)

	fmt.Fprintln(w, `</div>`)
}

func formatPct(num, denom int) string {
	result := ""
	if denom != 0 {
		result = fmt.Sprintf("%.1f%%", float32(num)*100.0/float32(denom))
	}
	return result
}

type playerStats struct {
	ClassFreq     int
	WinCount      int
	JointWinCount int
	BestHand      string
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
	printHeadCell := func(content string, colspan int) {
		colspanStr := ""
		if colspan > 1 {
			colspanStr = fmt.Sprintf(` colspan="%v"`, colspan)
		}
		fmt.Fprintf(w, `<th class="countTable"%v>%v</th>`, colspanStr, content)
	}
	printRow := func(handClass string, yourStats, bestOppStats, randOppStats playerStats, isSummary bool) {
		cssClass := ""
		if isSummary {
			cssClass = ` class="summary"`
		}
		fmt.Fprintf(w, `<tr%v>`, cssClass)
		printStringCell(handClass)
		printNumCell(yourStats.ClassFreq)
		printPctCell(yourStats.ClassFreq, simulator.HandCount)
		printNumCell(yourStats.WinCount)
		printPctCell(yourStats.WinCount, simulator.HandCount)
		printNumCell(yourStats.JointWinCount)
		printPctCell(yourStats.JointWinCount, simulator.HandCount)
		printStringCell(yourStats.BestHand)
		printNumCell(bestOppStats.ClassFreq)
		printPctCell(bestOppStats.ClassFreq, simulator.HandCount)
		printNumCell(bestOppStats.WinCount)
		printPctCell(bestOppStats.WinCount, simulator.HandCount)
		printStringCell(bestOppStats.BestHand)
		printNumCell(randOppStats.ClassFreq)
		printPctCell(randOppStats.ClassFreq, simulator.HandCount)
		printNumCell(randOppStats.WinCount)
		printPctCell(randOppStats.WinCount, simulator.HandCount)
		fmt.Fprintf(w, "</tr>")
	}
	fmt.Fprintf(w, `<table class="countTable"><tr><th class="countTable" rowspan="2">Hand</th>`)
	printHeadCell("For you", 7)
	printHeadCell("For best opponent", 5)
	printHeadCell("For random opponent", 5)
	fmt.Fprintf(w, `</tr><tr>`)
	printHeadCell("Freq", 2)
	printHeadCell("Wins", 2)
	printHeadCell("Joint wins", 2)
	printHeadCell("Best hand", 1)
	printHeadCell("Freq", 2)
	printHeadCell("Wins", 2)
	printHeadCell("Best hand", 1)
	printHeadCell("Freq", 2)
	printHeadCell("Wins", 2)
	fmt.Fprintf(w, `</tr>`)
	printRow("All", playerStats{simulator.HandCount, simulator.WinCount, simulator.JointWinCount, simulator.BestHand.PrettyPrint()}, playerStats{simulator.HandCount, simulator.BestOpponentWinCount, -1, simulator.BestOppHand.PrettyPrint()}, playerStats{simulator.HandCount, simulator.RandomOpponentWinCount, -1, ""}, true)
	for class := range simulator.OurClassCounts {
		bestHand := ""
		if simulator.OurClassCounts[class] > 0 {
			bestHand = simulator.ClassBestHands[class].PrettyPrint()
		}
		bestOppHand := ""
		if simulator.BestOpponentClassCounts[class] > 0 {
			bestOppHand = simulator.ClassBestOppHands[class].PrettyPrint()
		}
		printRow(poker.HandClass(class).String(), playerStats{simulator.OurClassCounts[class], simulator.ClassWinCounts[class], simulator.ClassJointWinCounts[class], bestHand}, playerStats{simulator.BestOpponentClassCounts[class], simulator.ClassBestOppWinCounts[class], -1, bestOppHand}, playerStats{simulator.RandomOpponentClassCounts[class], simulator.ClassRandOppWinCounts[class], -1, ""}, false)
	}
	fmt.Fprintf(w, "</table>")
}

func sampleCards(inputTableCards, inputYourCards []poker.Card) ([]string, []string) {
	samplePack := poker.SamplePack(inputTableCards, inputYourCards)
	tableCards, yourCards, _ := samplePack.PlayHoldem(1)
	tcStrings := make([]string, len(tableCards))
	for i, c := range tableCards {
		tcStrings[i] = c.String()
	}
	ycStrings := make([]string, len(yourCards[0]))
	for i, c := range yourCards[0] {
		ycStrings[i] = c.String()
	}
	return tcStrings, ycStrings
}

func SimulateHoldem(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, `<html lang="en">`)
	fmt.Fprintln(w, "<head><title>Texas Hold'em simulator</title>")
	fmt.Fprintln(w, `<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">`)
	fmt.Fprintln(w, "<style>")
	fmt.Fprintln(w, "td.formcell { vertical-align: top }")
	fmt.Fprintln(w, "table.countTable { border-collapse: collapse; }")
	fmt.Fprintln(w, "th.countTable { text-align: center }")
	fmt.Fprintln(w, "th.countTable, td.countTable { border: 1px solid black; padding: 3px; }")
	fmt.Fprintln(w, "td.numcell { text-align: right }")
	fmt.Fprintln(w, "td.zero { color: lightgrey }")
	fmt.Fprintln(w, ".summary { font-weight: bold }")
	fmt.Fprintln(w, "</style>")
	fmt.Fprintf(w, `<script src="//ajax.googleapis.com/ajax/libs/jquery/1.8.2/jquery.min.js"></script>`)
	fmt.Fprintf(w, `<script src="//code.highcharts.com/highcharts.js"></script>`)
	fmt.Fprintf(w, "</head><body>")
	fmt.Fprintln(w, `<div class="container-fluid">`)
	fmt.Fprintln(w, "<h1>Texas Hold'em Simulator</h1>")

	players, err := getPlayers(req)
	if err != nil {
		// Use a template for security as error messages will often contain raw user input
		t := template.Must(template.New("error").Parse("<p>Could not get player count: {{.}}</p></div></body></html>"))
		t.Execute(w, err.Error())
		return
	}

	tableCards, yourCards, handsToPlay, err := simulationParams(req)
	if err != nil {
		t := template.Must(template.New("error").Parse("<p>Could not get simulation parameters: {{.}}</p></div></body></html>"))
		t.Execute(w, err.Error())
		return
	}

	simulator := poker.Simulator{}
	simulator.SimulateHoldem(tableCards, yourCards, players, handsToPlay)
	simTableCards, simYourCards := sampleCards(tableCards, yourCards)
	fmt.Fprintln(w, `<div class="row"><div class="col-xs-12">`)
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
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Players</b></td><td><input id="playercount" type="text" name="%v" value="%v"/></td>`, playersKey, players)
	fmt.Fprintf(w, `<td><button type="button" onclick="$('#playercount').val(Math.max(2, $('#playercount').val()-1))">Fewer</button></td>`)
	fmt.Fprintln(w, `</tr>`)
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Your cards</b><br/><i>(comma-separated, e.g. 'KD,10H')</i></td><td>%v <input type="text" id="yourcards" name="%v" value="%v"/></td>`, formatCards(yourCards), yourCardsKey, cardText(yourCards))
	fmt.Fprintf(w, `<td><button type="button" onclick="$('#yourcards').val('%s')">Use sample</button></td></tr>`, strings.Join(simYourCards, ","))
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Table cards</b></td><td>%v <input id="tablecards" type="text" name="%v" value="%v"/></td>`, formatCards(tableCards), tableCardsKey, cardText(tableCards))
	fmt.Fprintln(w, `<td>`)
	fmt.Fprintf(w, `<button type="button" onclick="$('#tablecards').val('%s')">Use sample flop</button>`, strings.Join(simTableCards[:3], ","))
	fmt.Fprintf(w, `<button type="button" onclick="$('#tablecards').val('%s')">Use sample turn</button>`, strings.Join(simTableCards[:4], ","))
	fmt.Fprintf(w, `<button type="button" onclick="$('#tablecards').val('%s')">Use sample river</button>`, strings.Join(simTableCards[:5], ","))
	fmt.Fprintln(w, `</td></tr>`)
	fmt.Fprintf(w, `<tr><td class="formcell"><b>Simulations</b></td><td><input type="text" name="%v" value="%v"/></td></tr>`, simCountKey, simulator.HandCount)
	fmt.Fprintf(w, "</td></tr></table>")
	fmt.Fprintf(w, "</form></div></div>")

	fmt.Fprintf(w, "<h2>Results</h2>")
	printResultGraphs(w, simulator, tableCards, yourCards)
	fmt.Fprintln(w, `<div class="row"><div class="col-xs-12">`)
	printResultTable(w, simulator)
	fmt.Fprintln(w, `</div></div>`)

	fmt.Fprintf(w, "</div></body></html>")
}

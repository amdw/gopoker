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
package poker_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amdw/gopoker/holdem"
	"github.com/amdw/gopoker/poker"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strings"
)

const yourCardsKey = "yours"
const tableCardsKey = "table"
const simCountKey = "simcount"
const forceComputeKey = "compute"

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

func makeYourSeries(simulator *poker.Simulator) ([]string, []map[string]interface{}) {
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
	series := []map[string]interface{}{
		map[string]interface{}{"name": "Sole winner", "data": soleWinData},
		map[string]interface{}{"name": "Joint winner", "data": jointWinData},
		map[string]interface{}{"name": "Loser", "data": lossData},
		map[string]interface{}{"name": "Overall", "type": "pie", "data": overallData, "size": 100, "center": []string{"60%", "25%"}, "showInLegend": false, "dataLabels": map[string]interface{}{"enabled": true, "format": "{point.name} {y:.1f}%"}},
	}
	return handNames, series
}

func makeBestOppSeries(simulator *poker.Simulator) []map[string]interface{} {
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

func printResultGraphs(w http.ResponseWriter, simulator *poker.Simulator) {
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
		result = fmt.Sprintf("%.1f%%", float64(num)*100.0/float64(denom))
	}
	return result
}

type playerStats struct {
	ClassFreq     int
	WinCount      int
	JointWinCount int
	BestHand      string
}

func printResultTable(w http.ResponseWriter, simulator *poker.Simulator) {
	cssClass := func(isNum, isZero bool) string {
		classes := []string{}
		if isNum {
			classes = append(classes, "numcell")
		}
		if isZero {
			classes = append(classes, "zero")
		}
		return strings.Join(classes, " ")
	}
	printStringCell := func(content string) {
		fmt.Fprintf(w, `<td class="%v">%v</td>`, cssClass(false, false), content)
		fmt.Fprintln(w)
	}
	printNumCell := func(content int) {
		fmt.Fprintf(w, `<td class="%v">%v</td>`, cssClass(true, content == 0), content)
		fmt.Fprintln(w)
	}
	printPctCell := func(num, denom int) {
		fmt.Fprintf(w, `<td class="%v">%v</td>`, cssClass(true, num == 0), formatPct(num, denom))
		fmt.Fprintln(w)
	}
	printHeadCell := func(content string, colspan int) {
		colspanStr := ""
		if colspan > 1 {
			colspanStr = fmt.Sprintf(` colspan="%v"`, colspan)
		}
		fmt.Fprintf(w, `<th%v>%v</th>`, colspanStr, content)
		fmt.Fprintln(w)
	}
	printRow := func(handClass string, yourStats, bestOppStats, randOppStats playerStats, isSummary bool) {
		cssClass := ""
		if isSummary {
			cssClass = ` class="summary"`
		}
		fmt.Fprintf(w, `<tr%v>`, cssClass)
		fmt.Fprintln(w)
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
		fmt.Fprintln(w, "</tr>")
	}
	fmt.Fprintln(w, `<div class="table-responsive"><table class="table table-bordered table-condensed"><tr><th rowspan="2">Hand</th>`)
	printHeadCell("For you", 7)
	printHeadCell("For best opponent", 5)
	printHeadCell("For random opponent", 4)
	fmt.Fprintln(w, `</tr><tr>`)
	printHeadCell("Freq", 2)
	printHeadCell("Wins", 2)
	printHeadCell("Joint wins", 2)
	printHeadCell("Best hand", 1)
	printHeadCell("Freq", 2)
	printHeadCell("Wins", 2)
	printHeadCell("Best hand", 1)
	printHeadCell("Freq", 2)
	printHeadCell("Wins", 2)
	fmt.Fprintln(w, `</tr>`)
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
	fmt.Fprintf(w, "</table></div>")
}

func loadStaticFiles(staticBaseDir string) (*os.File, *os.File, *os.File, error) {
	filenames := []string{"simulation_head.html", "simulation_foot.html", "simulation.js"}
	files := make([]*os.File, len(filenames))
	for i, filename := range filenames {
		path := path.Join(staticBaseDir, filename)
		var err error
		files[i], err = os.Open(path)
		if err != nil {
			return nil, nil, nil, errors.New(fmt.Sprintf("Could not load %v: %v", path, err))
		}
	}
	return files[0], files[1], files[2], nil
}

func writeStaticFile(file *os.File, w http.ResponseWriter) bool {
	_, err := io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not write static file: %v", err), http.StatusInternalServerError)
		return false
	}
	return true
}

func cardsJson(cards []poker.Card) string {
	cardStrings := make([]string, len(cards))
	for i, card := range cards {
		cardStrings[i] = card.String()
	}
	jsonBytes, err := json.Marshal(cardStrings)
	if err != nil {
		panic(fmt.Sprintf("Unable to marshal cards %v: %v", cards, err))
	}
	return string(jsonBytes)
}

func SimulateHoldem(staticBaseDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		headFile, footFile, jsFile, err := loadStaticFiles(staticBaseDir)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error loading static files: %v", err), http.StatusInternalServerError)
			return
		}
		defer headFile.Close()
		defer footFile.Close()
		defer jsFile.Close()

		params, err := getSimulationParams(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get simulation parameters: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if !writeStaticFile(headFile, w) {
			return
		}

		breakEvenStr := "undefined"

		if len(params.tableCards) > 0 || len(params.yourCards) > 0 || params.forceComputation {
			simulator := holdem.SimulateHoldem(params.tableCards, params.yourCards, params.players, params.handsToPlay)

			fmt.Fprintf(w, "<h2>Results</h2>")

			breakEven := simulator.PotOddsBreakEven()
			if math.IsInf(breakEven, 1) {
				breakEvenStr = "Infinity"
			} else {
				breakEvenStr = fmt.Sprintf("%v", breakEven)
			}
			fmt.Fprintln(w, `<div class="row"><div class="col-xs-12"><div class="form-group"><form>`)
			fmt.Fprintln(w, `<label for="potsize">Pot size</label>`)
			fmt.Fprintln(w, `<input id="potsize" type="text" name="potsize" ng-model="potSize" class="form-control"/>`)
			fmt.Fprintln(w, `<span ng-bind-html="potOddsMessage()"></span>`)
			fmt.Fprintln(w, `</form></div></div></div>`)

			printResultGraphs(w, simulator)

			fmt.Fprintln(w, `<div class="row"><div class="col-xs-12">`)
			printResultTable(w, simulator)
			fmt.Fprintln(w, `</div></div>`)
		}

		fmt.Fprintln(w, "<script>")
		fmt.Fprintf(w, "var initPlayerCount = %v;\n", params.players)
		fmt.Fprintf(w, "var initYourCards = %v;\n", cardsJson(params.yourCards))
		fmt.Fprintf(w, "var initTableCards = %v;\n", cardsJson(params.tableCards))
		fmt.Fprintf(w, "var initSimCount = %v;\n", params.handsToPlay)
		fmt.Fprintf(w, "var potOddsBreakEven = %v;\n", breakEvenStr)

		if !writeStaticFile(jsFile, w) {
			return
		}
		fmt.Fprintln(w, "</script>")

		if !writeStaticFile(footFile, w) {
			return
		}
	}
}

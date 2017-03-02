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
	"fmt"
	"github.com/amdw/gopoker/poker"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

const baseUrl = "http://example.com"

func assertStatus(expectedStatus int, expectedContentType string, rec *httptest.ResponseRecorder, t *testing.T) {
	if rec.Code != expectedStatus {
		t.Errorf("Expected HTTP status %v, got %v: %v", expectedStatus, rec.Code, rec.Body.String())
	}
	contentType := rec.Result().Header["Content-Type"]
	if !reflect.DeepEqual([]string{expectedContentType}, contentType) {
		t.Errorf("Expected content type %v, found %v: %v", expectedContentType, contentType, rec.Body.String())
	}

}

func assertOkHtml(rec *httptest.ResponseRecorder, t *testing.T) {
	assertStatus(200, "text/html; charset=utf-8", rec, t)
}

func assertOkJson(rec *httptest.ResponseRecorder, t *testing.T) {
	assertStatus(200, "application/json", rec, t)
}

func assertBadRequest(rec *httptest.ResponseRecorder, t *testing.T) {
	assertStatus(http.StatusBadRequest, "text/plain; charset=utf-8", rec, t)
}

func TestPlayHoldem(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/holdem/play", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	PlayHoldem(rec, req)
	assertOkHtml(rec, t)
}

func TestPlayHoldemBadPlayers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/holdem/play?players=wibble", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	PlayHoldem(rec, req)
	assertBadRequest(rec, t)
}

func setupSimStaticAssets(t *testing.T) string {
	dir, err := ioutil.TempDir("", "gopokersimulatorstatic")
	if err != nil {
		t.Fatalf("Could not create temp dir for static assets: %v", err)
	}
	t.Log("Created temp dir", dir)

	filenames := []string{"simulation_head.html", "simulation_foot.html", "simulation.js"}
	for _, filename := range filenames {
		path := path.Join(dir, filename)
		err := ioutil.WriteFile(path, []byte(filename), 0644)
		if err != nil {
			t.Fatalf("Could not write temp file %v: %v", path, err)
		}
	}

	return dir
}

func TestHoldemSim(t *testing.T) {
	dir := setupSimStaticAssets(t)
	defer os.RemoveAll(dir)

	urls := []string{
		fmt.Sprintf("%v/holdem/simulate?compute=false", baseUrl),
		fmt.Sprintf("%v/holdem/simulate?compute=true", baseUrl),
	}
	for _, url := range urls {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Could not generate HTTP request: %v", err)
		}
		SimulateHoldem(dir)(rec, req)
		assertOkHtml(rec, t)
	}
}

func TestHoldemSimInputValidation(t *testing.T) {
	dir := setupSimStaticAssets(t)
	defer os.RemoveAll(dir)

	tooManyTableCards := "yours=" + url.QueryEscape("AS,QD") + "&table=" + url.QueryEscape("2S,3S,4S,5S,6S,7S")
	duplicateCard := "yours=" + url.QueryEscape("AS,QD") + "&table=" + url.QueryEscape("QD,2S,3S")
	tests := map[string]string{
		"players=wibble":                          "Could not get player count",
		"yours=" + url.QueryEscape("AS,QZ"):       "Illegally formatted card \"QZ\"",
		"table=" + url.QueryEscape("2D,3S,QZ,AD"): "Illegally formatted card \"QZ\"",
		"yours=" + url.QueryEscape("AS,QD,3S"):    "Maximum of 2 player cards allowed, found 3",
		tooManyTableCards:                         "Maximum of 5 table cards allowed, found 6",
		duplicateCard:                             "Found duplicate card QD",
		"simcount=wibble":                         "Could not parse simcount",
	}

	for query, expectedError := range tests {
		rec := httptest.NewRecorder()
		url := fmt.Sprintf("%v/holdem/simulate?%v", baseUrl, query)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Could not generate HTTP request: %v", err)
		}
		SimulateHoldem(dir)(rec, req)
		assertBadRequest(rec, t)
		if !strings.Contains(rec.Body.String(), expectedError) {
			t.Errorf("Could not find expected error '%v' in response for %v: %v", expectedError, url, rec.Body.String())
		}
	}
}

func TestHoldemStartingCardsHome(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/holdem/startingcards", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	// Create a temp dir to hold the HTML (doesn't matter that it's not real HTML)
	// This is easier than locating the actual HTML file reliably...
	dir, err := ioutil.TempDir("", "gopokerstartcardsstatic")
	if err != nil {
		t.Fatalf("Could not create temp dir for HTML: %v", err)
	}
	defer os.RemoveAll(dir)
	t.Log("Created temp dir", dir)
	filename := path.Join(dir, "starting_cards.html")
	err = ioutil.WriteFile(filename, []byte("temp html"), 0644)
	if err != nil {
		t.Fatalf("Could not write temp HTML file %v: %v", filename, err)
	}
	StartingCards(dir)(rec, req)
	assertOkHtml(rec, t)
}

func TestHoldemStartingCardsExecute(t *testing.T) {
	rec := httptest.NewRecorder()
	handCount := 12345
	players := 8
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/holdem/startingcards/sim?rank1=10&rank2=Q&samesuit=false&handstoplay=%v&players=%v", baseUrl, handCount, players), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateStartingCards(rec, req)
	assertOkJson(rec, t)
	sim := poker.Simulator{}
	json.Unmarshal(rec.Body.Bytes(), &sim)
	if sim.HandCount != handCount {
		t.Errorf("Expected hand count %v, found %v", handCount, sim.HandCount)
	}
	if sim.Players != players {
		t.Errorf("Expected players %v, found %v", players, sim.Players)
	}
}

func TestHoldemBadStartingCards(t *testing.T) {
	// Can't be both same rank and same suit
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/holdem/startingcards/sim?rank1=A&rank2=A&samesuit=true", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateStartingCards(rec, req)
	assertBadRequest(rec, t)

	// Bad rank
	rec = httptest.NewRecorder()
	req, err = http.NewRequest("GET", fmt.Sprintf("%v/holdem/startingcards/sim?rank1=A&rank2=Z&samesuit=false", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateStartingCards(rec, req)
	assertBadRequest(rec, t)
}

func TestOmaha8Simulation(t *testing.T) {
	dir := setupSimStaticAssets(t)
	defer os.RemoveAll(dir)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/omaha8/simulate", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateOmaha8(rec, req)
	assertOkHtml(rec, t)
}

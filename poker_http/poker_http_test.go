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
	"fmt"
	"github.com/amdw/gopoker/poker"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

const baseUrl = "http://example.com"

func TestGame(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/play", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	PlayHoldem(rec, req)
	if rec.Code != 200 {
		t.Errorf("Got HTTP error %s", rec.Code)
	}
}

func TestSim(t *testing.T) {
	// Create temp dir to hold the static assets
	dir, err := ioutil.TempDir("", "gopokersimulatorstatic")
	if err != nil {
		t.Fatalf("Could not create temp dir for static assets: %v", err)
	}
	defer os.RemoveAll(dir)
	t.Log("Created temp dir", dir)

	filenames := []string{"simulation_head.html", "simulation_foot.html", "simulation.js"}
	for _, filename := range filenames {
		path := path.Join(dir, filename)
		err = ioutil.WriteFile(path, []byte(filename), 0644)
		if err != nil {
			t.Fatalf("Could not write temp file %v: %v", path, err)
		}
	}

	urls := []string{
		fmt.Sprintf("%v/simulate?runsim=false", baseUrl),
		fmt.Sprintf("%v/simulate?runsim=true", baseUrl),
	}
	for _, url := range urls {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Could not generate HTTP request: %v", err)
		}
		SimulateHoldem(dir)(rec, req)
		if rec.Code != 200 {
			t.Errorf("Got HTTP error %v: %v", rec.Code, rec.Body.String())
		}
	}
}

func TestStartingCardsHome(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/startingcards", baseUrl), nil)
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
	if rec.Code != 200 {
		t.Errorf("Got HTTP error %v: %v", rec.Code, rec.Body.String())
	}
}

func TestStartingCardsExecute(t *testing.T) {
	rec := httptest.NewRecorder()
	handCount := 12345
	players := 8
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/startingcards/sim?rank1=10&rank2=Q&samesuit=false&handstoplay=%v&players=%v", baseUrl, handCount, players), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateStartingCards(rec, req)
	if rec.Code != 200 {
		t.Fatalf("Got HTTP error %v: %v", rec.Code, rec.Body.String())
	}
	sim := poker.Simulator{}
	json.Unmarshal(rec.Body.Bytes(), &sim)
	if sim.HandCount != handCount {
		t.Errorf("Expected hand count %v, found %v", handCount, sim.HandCount)
	}
	if sim.Players != players {
		t.Errorf("Expected players %v, found %v", players, sim.Players)
	}
}

func TestBadStartingCards(t *testing.T) {
	// Can't be both same rank and same suit
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/startingcards/sim?rank1=A&rank2=A&samesuit=true", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateStartingCards(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v: %v", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	// Bad rank
	rec = httptest.NewRecorder()
	req, err = http.NewRequest("GET", fmt.Sprintf("%v/startingcards/sim?rank1=A&rank2=Z&samesuit=false", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateStartingCards(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v: %v", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

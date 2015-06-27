package poker_http

import (
	"encoding/json"
	"fmt"
	"gopoker/poker"
	"net/http"
	"net/http/httptest"
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
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/simulate", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	SimulateHoldem(rec, req)
	if rec.Code != 200 {
		t.Errorf("Got HTTP error %v", rec.Code)
	}
}

func TestStartingCardsHome(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/startingcards", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	StartingCards(rec, req)
	if rec.Code != 200 {
		t.Errorf("Got HTTP error %v", rec.Code)
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

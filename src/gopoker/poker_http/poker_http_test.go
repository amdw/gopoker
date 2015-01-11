package poker_http

import (
	"fmt"
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

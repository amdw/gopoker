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
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestPlayOmaha8(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/omaha8/play", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not generate HTTP request: %v", err)
	}
	PlayOmaha8(rec, req)
	if rec.Code != 200 {
		t.Errorf("Got HTTP error %v: %v", rec.Code, rec.Body.String())
	}
	contentType := rec.Result().Header["Content-Type"]
	if !reflect.DeepEqual([]string{"text/html; charset=utf-8"}, contentType) {
		t.Errorf("Expected HTML response, found %v", contentType)
	}
}

func TestPlayOmaha8ErrorHandling(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/omaha8/play?players=wibble", baseUrl), nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}
	PlayOmaha8(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %v, found %v: %v", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	contentType := rec.Result().Header["Content-Type"]
	if !reflect.DeepEqual([]string{"text/plain; charset=utf-8"}, contentType) {
		t.Errorf("Expected plain-text response, found %v", contentType)
	}
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

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")

	return router
}

// Not really unit tests of router
func TestGetRequest(t *testing.T) {

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Router().ServeHTTP(w, r)

}

func TestPutRequest(t *testing.T) {

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	Router().ServeHTTP(w, r)

}

func testLogin(t *testing.T) {

}

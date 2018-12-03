// Mock Tests.
package main

import (
	"dataRep/repo"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var env *Env

func TestSetupSuite(t *testing.T) {
	dataSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8", dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		t.Fatal(err)
	}

	env = &Env{&repo.DB{db}}
}

func TestEnv_CabTripsHandler(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=D7D598CD99978BD012A87A76A7C891B7,F81EBC0D7805F6AF3E7C57038C951D4B&date=2013-12-01&skipCache=false",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	assert.Equal(t,
		`{
  	"D7D598CD99978BD012A87A76A7C891B7": "3",
  	"F81EBC0D7805F6AF3E7C57038C951D4B": "Data not found."
  }`,
		rr.Body.String(),
		"handler returned unexpected body")
}

func TestEnv_CabTripsHandlerInValidMedallion(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=abcd&date=2013-12-01&skipCache=false",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	assert.Equal(t,
		`{
  	"abcd": "Data not found."
  }`,
		rr.Body.String(),
		"handler returned unexpected body")
}

func TestEnv_CabTripsHandlerInValidDate(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=D7D598CD99978BD012A87A76A7C891B7&date=21321&skipCache=true",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	assert.Equal(t,
		`{
  	"D7D598CD99978BD012A87A76A7C891B7": "Data not found."
  }`,
		rr.Body.String(),
		"handler returned unexpected body")
}

func TestEnv_CabTripsHandlerInValidskipCache(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=D7D598CD99978BD012A87A76A7C891B7,F81EBC0D7805F6AF3E7C57038C951D4B&date=2013-12-01&skipCache=xyz",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	assert.Equal(t,
		`{
  	"D7D598CD99978BD012A87A76A7C891B7": "3",
  	"F81EBC0D7805F6AF3E7C57038C951D4B": "Data not found."
  }`,
		rr.Body.String(),
		"handler returned unexpected body")
}

func TestEnv_CabTripsHandlerNoskipCache(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=D7D598CD99978BD012A87A76A7C891B7&date=2013-12-01",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	assert.Equal(t,
		`{
  	"D7D598CD99978BD012A87A76A7C891B7": "3"
  }`,
		rr.Body.String(),
		"handler returned unexpected body")
}

func TestEnv_CabTripsHandlerEmptymedallion(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=&date=2013-12-01",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}
}

func TestEnv_CabTripsHandlerEmptyDate(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/trips?medallion=D7D598CD99978BD012A87A76A7C891B7&date=",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.CabTripsHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}
}

func TestEnv_ClearCacheHandler(t *testing.T) {
	req, err := http.NewRequest("GET",
		"http://localhost:8984/cab/clearcache",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(env.ClearCacheHandler)
	h.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

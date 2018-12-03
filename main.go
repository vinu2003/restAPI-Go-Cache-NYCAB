package main

import (
	"bytes"
	"dataRep/redisCache"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"dataRep/repo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
)

type Env struct {
	db repo.Datastore
}

const (
	dbName = "mysqld@localhost"
	dbPass = "safeinJesus123"
	dbHost = "localhost"
	dbPort = "3306"
)

var storage redisCache.Storage

func init() {
	var err error
	if storage, err = redisCache.NewStorage("redis://:United123@localhost:6379/1"); err != nil {
		panic(err)
	}
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "  ", "\t")
	return out.Bytes(), err
}

func writeJson(w http.ResponseWriter, data interface{}) {
	bJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	prettyB, _ := prettyprint(bJson)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyB)
}

// Get the port env variable.
func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8984"
	}
	return ":" + port
}

func main() {
	dbSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8", dbPass, dbHost, dbPort, dbName)
	db, err := repo.NewDB(dbSource)
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/cab/trips", cached("15s", env.CabTripsHandler))
	r.HandleFunc("/cab/clearcache", env.ClearCacheHandler)

	log.Fatal(http.ListenAndServe(port(), handlers.CORS()(r)))
}

// GET METHOD to retrieve trip details.
func (env *Env) CabTripsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	queryVals := r.URL.Query()
	medallionArg := queryVals.Get("medallion")

	medallionArr := strings.Split(medallionArg, ",")

	date := queryVals.Get("date")

	var (
		result map[string]int
		err    error
	)
	result, err = env.db.GetTrips(medallionArr, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, result)
}

// Handler to flush the DB cache.
func (env *Env) ClearCacheHandler(w http.ResponseWriter, r *http.Request) {
	err := storage.FlushDB()
	if err.Err() != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}

	writeJson(w, err.String())
}

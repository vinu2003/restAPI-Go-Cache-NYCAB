package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

func invokeHandler(w http.ResponseWriter, r *http.Request, handler func(w http.ResponseWriter, r *http.Request), duration string) {
	c := httptest.NewRecorder()
	handler(c, r)

	for k, v := range c.Header() {
		w.Header()[k] = v
	}

	w.WriteHeader(c.Code)
	content := c.Body.Bytes()

	// Whenever there is a new read from database  *update* the cache for that particluar key.
	if d, err := time.ParseDuration(duration); err == nil {
		status, err := storage.Set(r.RequestURI, content, d)
		if err != nil {
			http.Error(w, status.String(), http.StatusExpectationFailed)
		} else {
			w.Header().Add("Cache Status", fmt.Sprintf("New page cached: %s for %s\n", r.RequestURI, duration))
		}
	} else {
		log.Println("Page not cached. err: ", err)
	}

	w.Write(content)
}

func cached(duration string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryVals := r.URL.Query()

		var cache bool
		cacheArg := queryVals.Get("skipCache")
		if strings.Contains(cacheArg, "true") || strings.Contains(cacheArg, "false") {
			b, err := strconv.ParseBool(cacheArg)
			if err != nil {
				cache = false
			}
			cache = b
		} else {
			cache = false
		}

		index := strings.Index(r.RequestURI, "&skipCache")
		if index != -1 {
			r.RequestURI = r.RequestURI[:index]
		}

		// skipCache true, read from DB, skipping to look from cache.
		// skipCache false, make an attempt to read from cache...
		if !cache {
			content, _ := storage.Get(r.RequestURI)
			if content != nil {
				w.Header().Add("Status", "Cache Hit!")
				w.Write(content)
			} else {
				invokeHandler(w, r, handler, duration)
			}
		} else {
			invokeHandler(w, r, handler, duration)

		}

	})
}

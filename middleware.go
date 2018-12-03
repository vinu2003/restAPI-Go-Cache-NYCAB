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

	if d, err := time.ParseDuration(duration); err == nil {
		fmt.Printf("New page cached: %s for %s\n", r.RequestURI, duration)
		storage.Set(r.RequestURI, content, d)
	} else {
		fmt.Printf("Page not cached. err: %s\n", err)
	}

	w.Write(content)
}

func cached(duration string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RequestURI:", r.RequestURI)
		queryVals := r.URL.Query()

		var cache bool
		cacheArg := queryVals.Get("skipCache")
		if strings.Contains(cacheArg, "true") || strings.Contains(cacheArg, "false") {
			b, err := strconv.ParseBool(cacheArg)
			if err != nil {
				log.Println("INFO: Unable to Parse bool value - ", err)
				log.Println("INFO: skipCache set to false")
				cache = false
			}
			cache = b
		} else {
			log.Println("INFO: skipCache set to false")
			cache = false
		}

		index := strings.Index(r.RequestURI, "&skipCache")
		fmt.Println("index: ", index)
		if index != -1 {
			r.RequestURI = r.RequestURI[:index]
		}
		fmt.Println("New RequestURI : ", r.RequestURI)

		// skipCache true, read from DB, skipping to look from cache.
		// skipCache false, make an attempt to read from cache...
		if !cache {
			content := storage.Get(r.RequestURI)
			if content != nil {
				fmt.Println("Cache Hit!")
				w.Write(content)
			} else {
				invokeHandler(w, r, handler, duration)
			}
		} else {
			invokeHandler(w, r, handler, duration)

		}

	})
}

// cache implements redis get/set functionalities.
package redisCache

import (
	"fmt"
	r "gopkg.in/redis.v5"
	"log"
	"strings"
	"time"
)

//Storage mecanism for caching strings.
type Storage interface {
	Get(key string) []byte
	Set(key string, content []byte, duration time.Duration)
	FlushDB() r.StatusCmd
}

// Storage mechanism for caching strings.
type Store struct {
	client *r.Client
}

//NewStorage creates a new redis storage.
func NewStorage(url string) (*Store, error) {
	var (
		opts *r.Options
		err  error
	)
	if opts, err = r.ParseURL(url); err != nil {
		return nil, err
	}

	fmt.Println("NewStorage : ", opts)
	return &Store{
		client: r.NewClient(opts),
	}, nil
}

// Get a cached content by key.
func (s *Store) Get(key string) []byte {
	val, err := s.client.Get("_PAGE_CACHE_" + key).Bytes()
	if err == nil && len(val) == 0 {
		log.Println("GET error: key doesn't exists :", err)
	}
	return val
}

// Set a cached content by key.
func (s *Store) Set(key string, content []byte, duration time.Duration) {
	err := s.client.Set("_PAGE_CACHE_"+key, content, duration)
	// nil doesnt mean it is an Error. It can be OK.
	if err != nil && !(strings.Contains(err.String(), "OK")) {
		log.Println("SET Error: ", err)
	} else if strings.Contains(err.String(), "OK") {
		log.Println(err.String())
	}
}

// FlushDB.
func (s *Store) FlushDB() r.StatusCmd {
	err := s.client.FlushDb()
	if err != nil && !(strings.Contains(err.String(), "OK")) {
		log.Println("FlushDB Error: ", err)
	} else if strings.Contains(err.String(), "OK") {
		log.Println(err.String())
	}

	return *err
}

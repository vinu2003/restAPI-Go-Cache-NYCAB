// cache implements redis get/set flush functionalities.
package redisCache

import (
	"errors"
	"fmt"
	r "gopkg.in/redis.v5"
	"strings"
	"time"
)

//Storage mechanism for caching strings.
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte, duration time.Duration) (r.StatusCmd, error)
	FlushDB() (r.StatusCmd, error)
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
		return nil, errors.New(fmt.Sprintf("ERROR: Redis URL parsing %v", err))
	}

	fmt.Println("NewStorage : ", opts)
	return &Store{
		client: r.NewClient(opts),
	}, nil
}

// Get a cached content by key.
func (s *Store) Get(key string) ([]byte, error) {
	val, err := s.client.Get("_PAGE_CACHE_" + key).Bytes()
	if err == nil && len(val) == 0 {
		return nil, errors.New(fmt.Sprintf("ERROR: Redis SET %v", err))
	}
	return val, nil
}

// Set a cached content by key.
func (s *Store) Set(key string, content []byte, duration time.Duration) (r.StatusCmd, error) {
	cmdStatus := s.client.Set("_PAGE_CACHE_"+key, content, duration)

	// nil doesnt mean it is an Error. It can be OK.
	if cmdStatus != nil && !(strings.Contains(cmdStatus.String(), "OK")) {
		return *cmdStatus, errors.New(fmt.Sprintf("ERROR: Redis SET %v", cmdStatus.Err()))
	}

	return *cmdStatus, nil
}

// FlushDB.
func (s *Store) FlushDB() (r.StatusCmd, error) {
	cmdStatus := s.client.FlushDb()
	if cmdStatus != nil && !(strings.Contains(cmdStatus.String(), "OK")) {
		return *cmdStatus, errors.New(fmt.Sprintf("ERROR: Redis FlushDB %v", cmdStatus.Err()))
	}

	return *cmdStatus, nil
}

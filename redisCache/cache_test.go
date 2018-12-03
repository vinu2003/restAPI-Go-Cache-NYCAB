package redisCache

import (
	"bytes"
	"testing"
	"time"
)

var redisURL = "redis://:United123@localhost:6379/1"

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestWrongURL(t *testing.T) {
	storage, err := NewStorage("http://abcd")
	if err == nil || storage != nil {
		t.Fail()
	}
}

func TestGetEmpty(t *testing.T) {
	storage, _ := NewStorage(redisURL)
	content, _ := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func TestGetValue(t *testing.T) {
	storage, _ := NewStorage(redisURL)
	storage.Set("MY_KEY", []byte("123456"), parse("10s"))

	content, _ := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte("123456"))
}

func TestGetExpiredValue(t *testing.T) {
	storage, _ := NewStorage(redisURL)
	storage.Set("MY_KEY", []byte("123456"), parse("1s"))
	time.Sleep(parse("1s"))
	content, _ := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func TestStore_FlushDB(t *testing.T) {
	storage, _ := NewStorage(redisURL)
	storage.Set("MY_KEY", []byte("123456"), parse("50s"))

	// flush
	storage.FlushDB()

	content, _ := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func assertContentEquals(t *testing.T, content, expected []byte) {
	if !bytes.Equal(content, expected) {
		t.Errorf("content should '%s', but was '%s'", expected, content)
	}
}

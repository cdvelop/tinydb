package tinydb

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type failStore struct{}

func (f *failStore) GetFile(filePath string) ([]byte, error)    { return nil, nil }
func (f *failStore) SetFile(filePath string, data []byte) error { return errors.New("disk full") }

func TestGet(t *testing.T) {
	store := newMockStore()
	db, _ := New("test.db", nil, store)
	db.Set("foo", "bar")

	t.Run("gets an existing key", func(t *testing.T) {
		val, err := db.Get("foo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "bar" {
			t.Errorf("expected value 'bar', got '%s'", val)
		}
	})

	t.Run("returns an error for a non-existent key", func(t *testing.T) {
		_, err := db.Get("baz")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestSet(t *testing.T) {
	store := newMockStore()
	db, _ := New("test.db", nil, store)

	t.Run("sets a new key-value pair", func(t *testing.T) {
		err := db.Set("foo", "bar")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, _ := db.Get("foo")
		if val != "bar" {
			t.Errorf("expected value 'bar', got '%s'", val)
		}
	})

	t.Run("updates an existing key", func(t *testing.T) {
		err := db.Set("foo", "baz")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, _ := db.Get("foo")
		if val != "baz" {
			t.Errorf("expected value 'baz', got '%s'", val)
		}
	})
}

func TestLogger(t *testing.T) {
	store := newMockStore()
	var buf bytes.Buffer
	db, _ := New("test.db", &buf, store)

	// Successful operations should not log
	db.Set("foo", "bar")
	if buf.Len() != 0 {
		t.Errorf("expected no logs for successful operations, got '%s'", buf.String())
	}

	// Now simulate a failing store to ensure errors are logged
	fs := &failStore{}
	var buf2 bytes.Buffer
	db2, _ := New("test.db", &buf2, fs)
	_ = db2.Set("a", "b")
	if !strings.Contains(buf2.String(), "error persisting") {
		t.Errorf("expected error log for failing persist, got '%s'", buf2.String())
	}
}

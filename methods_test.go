package tinydb

import (
	"bytes"
	"strings"
	"testing"
)

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

	db.Set("foo", "bar")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "insert key=foo") {
		t.Errorf("expected log to contain 'insert key=foo', got '%s'", logOutput)
	}
}

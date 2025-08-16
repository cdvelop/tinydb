# TinyDB
TinyGo–compatible key–value store with a minimal API.
It provides a simple way to persist string-based key–value pairs using a custom Store backend interface.
Unlike traditional databases, tinydb avoids heavy dependencies (like JSON, SQL, or reflection) and relies only on io.Writer for logging and a user–provided Store implementation for persistence.

---

## Features
- Minimal API: only `Get` and `Set`.
- Works entirely with **strings** (`key` and `value`).
- Pluggable storage with the `Store` interface.
- TinyGo–friendly (no `fmt`, no JSON).
- Data stored as `key=value` per line.

---

## API

### KVStore Interface
```go
type KVStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

Store Interface

type Store interface {
    GetFile(filePath string) ([]byte, error)    // Load file or create empty
    SetFile(filePath string, data []byte) error // Save full DB
}

Constructor

db, err := tinydb.New("mydb.tdb", logger, store)
if err != nil {
    panic(err)
}

name → logical DB name (usually a file path).

logger io.Writer → optional writer for logs.

store Store → backend implementation for persistence.


---

Example Usage

package main

import (
    "os"
    "github.com/cdvelop/tinydb"
)

// Example FileStore (simplified)
type FileStore struct{}

func (fs FileStore) GetFile(path string) ([]byte, error) {
    return os.ReadFile(path)
}
func (fs FileStore) SetFile(path string, data []byte) error {
    return os.WriteFile(path, data, 0644)
}

func main() {
    db, _ := tinydb.New("settings.tdb", os.Stdout, FileStore{})

    db.Set("username", "cesar")
    db.Set("theme", "dark")

    val, _ := db.Get("theme")
    println("Theme:", val)
}


---

Storage Format

tinydb stores data as simple key=value lines:

username=cesar
theme=dark
window=1024x768


---

License

MIT

---



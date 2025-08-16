package tinydb

import (
	"io"

	. "github.com/cdvelop/tinystring"
)

type pair struct {
	Key   string
	Value string
}

type TinyDB struct {
	name   string
	data   []pair
	logger io.Writer
	store  Store

	raw builder
}

// New creates or loads a database
func New(name string, logger io.Writer, store Store) (*TinyDB, error) {
	db := &TinyDB{
		name:   name,
		data:   make([]pair, 0),
		logger: logger,
		store:  store,
		raw:    Convert(),
	}

	// try to load DB from Store
	raw, err := store.GetFile(name)
	if err == nil && len(raw) > 0 {
		lines := Convert(string(raw)).Split("\n")
		for _, line := range lines {
			if Convert(line).TrimSpace().String() == "" {
				continue
			}
			kv := Convert(line).Split("=")
			if len(kv) == 2 {
				db.data = append(db.data, pair{
					Key:   kv[0],
					Value: kv[1],
				})
			}
		}
	}

	return db, nil
}

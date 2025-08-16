package tinydb

import (
	"io"
	. "github.com/cdvelop/tinystring"
)

// Store define la interfaz de persistencia
type Store interface {
	GetFile(filePath string) ([]byte, error)     // leer archivo (crear si no existe)
	SetFile(filePath string, data []byte) error // guardar archivo completo
}

// KVStore API mínima
type KVStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type pair struct {
	Key   string
	Value string
}

type tinydb struct {
	name   string
	data   []pair
	logger io.Writer
	store  Store
}

// New crea o carga una base de datos
func New(name string, logger io.Writer, store Store) (KVStore, error) {
	db := &tinydb{
		name:   name,
		data:   make([]pair, 0),
		logger: logger,
		store:  store,
	}

	// intentar cargar DB desde el Store
	raw, err := store.GetFile(name)
	if err == nil && len(raw) > 0 {
		lines := Split(string(raw), "\n")
		for _, line := range lines {
			if TrimSpace(line) == "" {
				continue
			}
			kv := tinystring.SplitN(line, "=", 2)
			if len(kv) == 2 {
				db.data = append(db.data, pair{
					Key:   kv[0],
					Value: kv[1],
				})
			}
		}
		db.log("db loaded: " + name)
	}

	return db, nil
}

func (t *tinydb) Get(key string) (string, error) {
	for _, p := range t.data {
		if p.Key == key {
			return p.Value, nil
		}
	}
	return "", Err("key not found: ", key)
}

func (t *tinydb) Set(key, value string) error {
	// buscar si existe
	for i, p := range t.data {
		if p.Key == key {
			t.data[i].Value = value
			return t.persist("update key=" + key)
		}
	}

	// insertar nuevo
	t.data = append(t.data, pair{Key: key, Value: value})
	return t.persist("insert key=" + key)
}

func (t *tinydb) persist(msg string) error {
	var lines []string
	for _, p := range t.data {
		lines = append(lines, p.Key+"="+p.Value)
	}

	raw := Join(lines, "\n")

	if err := t.store.SetFile(t.name, []byte(raw)); err != nil {
		return err
	}

	t.log(msg)
	return nil
}

func (t *tinydb) log(msg string) {
	if t.logger != nil {
		t.logger.Write([]byte(msg + "\n"))
	}
}

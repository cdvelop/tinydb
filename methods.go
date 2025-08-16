package tinydb

import (
	. "github.com/cdvelop/tinystring"
)

func (t *tinydb) Get(key string) (string, error) {
	for _, p := range t.data {
		if p.Key == key {
			return p.Value, nil
		}
	}
	return "", Err("key not found: ", key)
}

func (t *tinydb) Set(key, value string) error {
	// search if it exists
	for i, p := range t.data {
		if p.Key == key {
			t.data[i].Value = value
			return t.persist("update key=" + key)
		}
	}

	// insert new
	t.data = append(t.data, pair{Key: key, Value: value})
	return t.persist("insert key=" + key)
}

func (t *tinydb) persist(msg string) error {
	var lines []string
	for _, p := range t.data {
		lines = append(lines, p.Key+"="+p.Value)
	}

	raw := Convert(lines).Join("\n").String()

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

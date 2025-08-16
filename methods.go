package tinydb

import (
	. "github.com/cdvelop/tinystring"
)

func (t *TinyDB) Get(key string) (string, error) {
	for _, p := range t.data {
		if p.Key == key {
			return p.Value, nil
		}
	}
	return "", Err("key not found: ", key)
}

func (t *TinyDB) Set(key, value string) error {
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

func (t *TinyDB) persist(msg string) error {
	raw := Convert()
	for _, p := range t.data {
		raw.Write(p.Key)
		raw.Write("=")
		raw.Write(p.Value)
		raw.Write("\n")
	}

	raw := Convert(lines).Join("\n").String()

	if err := t.store.SetFile(t.name, raw.Byte()); err != nil {
		// log only on error
		t.log("error persisting: " + err.Error())
		return err
	}

	return nil
}

func (t *TinyDB) log(msg string) {
	if t.logger != nil {
		t.logger.Write([]byte(msg + "\n"))
	}
}

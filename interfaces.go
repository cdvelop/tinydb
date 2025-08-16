package tinydb

// Store defines the persistence interface
type Store interface {
	GetFile(filePath string) ([]byte, error)
	SetFile(filePath string, data []byte) error
}

// KVStore defines the minimum API
type KVStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type builder interface {
	Reset() error
	Write(v any) error
	Bytes() []byte
}

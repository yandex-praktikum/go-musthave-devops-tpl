package storage

import "errors"

type Storager interface {
	Write(key, value string) error
	Read(key string) (string, error)
}

// MemoryRepo структура
type MemoryRepo struct {
	db map[string]string
}

func NewMemoryRepo() MemoryRepo {
	return MemoryRepo{
		db: make(map[string]string),
	}
}

func (m MemoryRepo) Write(key, value string) error {
	m.db[key] = value
	return nil
}

func (m MemoryRepo) Read(key string) (string, error) {
	value, err := m.db[key]
	if !err {
		return "", errors.New("Значение по ключу не найдено, ключ: " + key)
	}

	return value, nil
}

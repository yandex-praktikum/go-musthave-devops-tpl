package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Gauge float64
type Counter int64

type Storager interface {
	Len() int
	Write(key, value string) error
	Read(key string) (string, error)
	Delete(key string) (string, bool)
	GetSchemaDump() map[string]string
}

// MemoryRepo структура
type MemoryRepo struct {
	db map[string]string
	*sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		db:      make(map[string]string),
		RWMutex: &sync.RWMutex{},
	}
}

func (m *MemoryRepo) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.db)
}

func (m MemoryRepo) Write(key, value string) error {
	m.Lock()
	defer m.Unlock()
	m.db[key] = value
	return nil
}

func (m *MemoryRepo) Delete(key string) (string, bool) {
	m.Lock()
	defer m.Unlock()
	oldValue, ok := m.db[key]
	if ok {
		delete(m.db, key)
	}
	return oldValue, ok
}

func (m MemoryRepo) Read(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()
	value, err := m.db[key]
	if !err {
		return "", errors.New("Значение по ключу не найдено, ключ: " + key)
	}

	return value, nil
}

func (m MemoryRepo) GetSchemaDump() map[string]string {
	return m.db
}

//MemStatsMemoryRepo - репо для приходящей статистики
type MemStatsMemoryRepo struct {
	storage Storager
}

//Создание MemStatsMemoryRepo с дефолтными значениями
func NewMemStatsMemoryRepo() MemStatsMemoryRepo {
	var memStatsStorage MemStatsMemoryRepo
	memStatsStorage.storage = NewMemoryRepo()

	memStatsStorage.storage.Write("Alloc", "0")
	memStatsStorage.storage.Write("BuckHashSys", "0")
	memStatsStorage.storage.Write("Frees", "0")
	memStatsStorage.storage.Write("GCCPUFraction", "0")
	memStatsStorage.storage.Write("GCSys", "0")

	memStatsStorage.storage.Write("HeapAlloc", "0")
	memStatsStorage.storage.Write("HeapIdle", "0")
	memStatsStorage.storage.Write("HeapInuse", "0")
	memStatsStorage.storage.Write("HeapObjects", "0")
	memStatsStorage.storage.Write("HeapReleased", "0")

	memStatsStorage.storage.Write("HeapSys", "0")
	memStatsStorage.storage.Write("LastGC", "0")
	memStatsStorage.storage.Write("Lookups", "0")
	memStatsStorage.storage.Write("MCacheInuse", "0")
	memStatsStorage.storage.Write("MCacheSys", "0")

	memStatsStorage.storage.Write("MSpanInuse", "0")
	memStatsStorage.storage.Write("MSpanSys", "0")
	memStatsStorage.storage.Write("Mallocs", "0")
	memStatsStorage.storage.Write("NextGC", "0")
	memStatsStorage.storage.Write("NumForcedGC", "0")

	memStatsStorage.storage.Write("NumGC", "0")
	memStatsStorage.storage.Write("OtherSys", "0")
	memStatsStorage.storage.Write("PauseTotalNs", "0")
	memStatsStorage.storage.Write("StackInuse", "0")
	memStatsStorage.storage.Write("StackSys", "0")

	memStatsStorage.storage.Write("Sys", "0")
	memStatsStorage.storage.Write("TotalAlloc", "0")
	memStatsStorage.storage.Write("PollCount", "0")
	memStatsStorage.storage.Write("RandomValue", "0")

	return memStatsStorage
}

func (memStatsStorage MemStatsMemoryRepo) UpdateGaugeValue(key string, value float64) error {
	return memStatsStorage.storage.Write(key, fmt.Sprintf("%v", value))
}

func (memStatsStorage MemStatsMemoryRepo) UpdateCounterValue(key string, value int64) error {
	//Чтение старого значения
	oldValue, err := memStatsStorage.storage.Read(key)
	if err != nil {
		oldValue = "0"
	}

	//Конвертация в число
	oldValueInt, err := strconv.ParseInt(oldValue, 10, 64)
	if err != nil {
		return errors.New("MemStats value is not int64")
	}

	newValue := fmt.Sprintf("%v", oldValueInt+value)
	memStatsStorage.storage.Write(key, newValue)

	return nil
}

func (memStatsStorage MemStatsMemoryRepo) ReadValue(key string) (string, error) {
	return memStatsStorage.storage.Read(key)
}

func (memStatsStorage MemStatsMemoryRepo) GetDBSchema() map[string]string {
	return memStatsStorage.storage.GetSchemaDump()
}

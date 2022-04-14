package cache

import (
	"fmt"
	"sync"

	"github.com/grt1st/wdproxy/storage"
)

type Storage struct {
	zone string
	M    map[interface{}]interface{}
	sync.RWMutex
}

func NewStorage(zone string) *Storage {
	s := Storage{
		zone: zone,
		M:    make(map[interface{}]interface{}, 0),
	}
	var err error
	if zone == "file" {
		var caches []storage.FileCache
		err = storage.Query(&caches, "")
		for _, v := range caches {
			s.M[v.Sum] = v
		}
	} else {
		var caches []storage.Document
		err = storage.Query(&caches, "")
		for _, v := range caches {
			s.M[v.Token] = v
		}
	}
	if err != nil {
		fmt.Println("[-] load caches failed: ", err)
	}
	return &s
}

func (s *Storage) Has(key interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.M[key]
	return ok
}

func (s *Storage) Get(key interface{}) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.M[key]
	return v, ok
}

func (s *Storage) Add(key, value interface{}) (err error) {
	s.Lock()
	defer s.Unlock()
	switch v := value.(type) {
	case storage.FileCache:
		s.M[v.Sum] = v
		err = storage.Create(&v)
	case storage.Document:
		s.M[v.Token] = v
		err = storage.Create(&v)
	default:
		err = fmt.Errorf("unknown cache type: %T", v)
	}
	return err
}

func (s *Storage) Update(key, value interface{}) (err error) {
	s.Lock()
	defer s.Unlock()
	switch v := value.(type) {
	case storage.FileCache:
		s.M[v.Sum] = v
		err = storage.Update(storage.FileCache{}, v, "sum = ?", key)
	case storage.Document:
		s.M[v.Token] = v
		err = storage.Update(storage.Document{}, v, "token = ?", key)
	default:
		err = fmt.Errorf("unknown cache type: %T", v)
	}
	return err
}

func (s *Storage) Del(key interface{}) error {
	s.Lock()
	defer s.Unlock()
	delete(s.M, key)
	if s.zone == "file" {
		return storage.DeleteByField(storage.FileCache{}, "sum", key)
	} else {
		return storage.DeleteByField(storage.Document{}, "token", key)
	}
}

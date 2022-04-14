package cache

import "sync"

type Memory struct {
	M map[interface{}]interface{}
	sync.RWMutex
}

func NewMemory() *Memory {
	M := Memory{
		M: make(map[interface{}]interface{}, 0),
	}
	return &M
}

func (m *Memory) Has(key interface{}) bool {
	m.Lock()
	defer m.Unlock()
	_, ok := m.M[key]
	if ok {
		return true
	}
	return false
}

func (m *Memory) Get(key interface{}) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.M[key]
	return v, ok
}

func (m *Memory) Add(key, value interface{}) error {
	m.Lock()
	defer m.Unlock()
	m.M[key] = value
	return nil
}

func (m *Memory) Update(key, value interface{}) error {
	return m.Add(key, value)
}

func (m *Memory) Del(key interface{}) error {
	m.Lock()
	defer m.Unlock()
	delete(m.M, key)
	return nil
}

package cache

import "sync"

type Memory struct {
	M map[int64][]byte
	sync.RWMutex
}

func NewMemory() *Memory {
	M := Memory{
		M: map[int64][]byte{},
	}
	return &M
}

func (m *Memory) Has(id int64) bool {
	m.Lock()
	defer m.Unlock()
	_, ok := m.M[id]
	if ok {
		return true
	}
	return false
}

func (m *Memory) Get(id int64) ([]byte, bool) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.M[id]
	return v, ok
}

func (m *Memory) Add(id int64, content []byte) error {
	m.Lock()
	defer m.Unlock()
	m.M[id] = content
	return nil
}

func (m *Memory) Del(id int64) error {
	m.Lock()
	defer m.Unlock()
	delete(m.M, id)
	return nil
}

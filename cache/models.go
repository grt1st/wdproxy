package cache

type Cacher interface {
	Has(id int64) bool
	Get(id int64) ([]byte, bool)
	Add(id int64, content []byte) error
	Del(id int64) error
}
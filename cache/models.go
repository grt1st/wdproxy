package cache

type Cacher interface {
	Has(key interface{}) bool
	Get(key interface{}) (interface{}, bool)
	Add(key, value interface{}) error
	Update(key, value interface{}) error
	Del(key interface{}) error
}

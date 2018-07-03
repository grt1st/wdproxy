package g

import "github.com/grt1st/wdproxy/cache"

var C cache.Cacher

func init() {
	C = cache.NewMemory()
}

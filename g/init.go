package g

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/grt1st/wdproxy/cache"
	"github.com/grt1st/wdproxy/storage"
)

var C cache.Cacher

func Init() {
	// 数据库等外部模块
	storage.Init()
	C = cache.NewMemory()
}

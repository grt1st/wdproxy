package g

import (
	"github.com/grt1st/wdproxy/cache"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strings"
)

// schema user password ip:port database
var Store = "mysql//root:123456@(127.0.0.1:3306)/cspp?charset=utf8mb4&parseTime=True&loc=Local"
var C cache.Cacher
var DB *gorm.DB

func init() {
	C = cache.NewMemory()
	initDB()
}

func initDB() {
	var err error
	parts := strings.SplitN(Store, "//", 2)
	DB, err = gorm.Open(parts[0], parts[1])
	if err != nil {
		log.Fatalln("[-] open database error.", err)
	}

	//DB.LogMode(true)
	DB.AutoMigrate(&WdproxyRecord{}, &WdproxyDomain{})
}

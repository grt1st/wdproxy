package storage

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

var Store = "mysql//linker:root123456@(127.0.0.1:3306)/mars?charset=utf8mb4&parseTime=True&loc=Local"
var DB *gorm.DB

func Init() {
	var err error
	parts := strings.SplitN(Store, "//", 2)
	DB, err = gorm.Open(parts[0], parts[1])
	if err != nil {
		panic(err)
	}

	DB.LogMode(true)
}

func Update(model, value interface{}, where string, args ...interface{}) error {
	return GetDB().Model(model).Where(where, args...).Updates(value).Error
}

func Create(value interface{}) error {
	return GetDB().Create(value).Error
}

func Query(model interface{}, where string, args ...interface{}) error {
	return GetDB().Where(where, args...).Find(model).Error
}

func Get(model interface{}, where string, args ...interface{}) error {
	return GetDB().Where(where, args...).Find(model).Error
}

func DeleteByField(model interface{}, field string, value interface{}) error {
	return GetDB().Delete(model, fmt.Sprintf("%s = ?", field), value).Error
}

func GetDB() *gorm.DB {
	return DB
}

func IsNotFound(err error) bool {
	return err != nil && (errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == gorm.ErrRecordNotFound.Error())
}

package g

import (
	"github.com/jinzhu/gorm"
)

type WdproxyRecord struct {
	gorm.Model
	WdproxyDomainID uint
	Url string `sql:"type:text(2000);"`
	Method string
	Status int
	Scheme string
	Path string
	ContentType string
	ContentLength uint
	RemoteAddr string
	Host string
	Port string
	Extension string
	HeaderRequest string `sql:"type:text(2000);"`
	HeaderResponse string `sql:"type:text(2000);"`
	BodyRequest string `sql:"type:text(60000);"`
	BodyResponse string `sql:"type:text(60000);"`
}

type WdproxyDomain struct {
	gorm.Model
	Value string `sql:"unique"`
	WdproxyRecords []WdproxyRecord
}
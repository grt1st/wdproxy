package storage

import "github.com/jinzhu/gorm"

type FileCache struct {
	gorm.Model
	Sum      string `gorm:"index:idx_file_sum"`
	Filename string
	Url      string
}

func (FileCache) TableName() string {
	return "proxy_file_cache"
}

type Document struct {
	gorm.Model
	Token        string `gorm:"index:idx_document_id"`
	Title        string `gorm:"type:varchar(1024);unique_index"`
	Owner        string `gorm:"type:varchar(512);unique_index"`
	LastEdited   int64
	FirstCreated int64
	URL          string `gorm:"type:varchar(1024);column:url"`
}

func (Document) TableName() string {
	return "proxy_document"
}

package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 创建一个模型，用于表示存储 IPFS 哈希的表
type FileHash struct {
	//gorm.Model
	FileName string `gorm:"not null"`
	IPFSHash string `gorm:"not null"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = connectDatabase()
	if err != nil {
		return
	}

}

// 连接 MySQL 数据库
func connectDatabase() (*gorm.DB, error) {
	// 修改为你的数据库连接信息
	db, err := gorm.Open("mysql", "wang:@tcp(localhost:3306)/ipfs?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		println(err)
		return nil, err
	}
	return db, nil
}
func GetDB() *gorm.DB {
	return db
}

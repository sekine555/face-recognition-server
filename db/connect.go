package db

import (
	"face-recognition/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB接続
func SqlConnect() (database *gorm.DB, err error) {
	dbConnectInfo := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		config.Config.DbUserName,
		config.Config.DbUserPassword,
		config.Config.DbHost,
		config.Config.DbPort,
		config.Config.DbName,
	)
	fmt.Println(dbConnectInfo)

	return gorm.Open(config.Config.DbDriverName, dbConnectInfo)
}
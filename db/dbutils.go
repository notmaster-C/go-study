package db

import (
	"fmt"
	"github.com/dzwvip/oracle"
	"go-study/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var prefix = "bigdata_"

// {"db_type":"oracle", "db_host":"202.115.154.76", "db_port":1521, "db_user":"ODS", "db_pass":"XjODS123","db_name":"XJSJZXDB"}
var (
	DB  *gorm.DB
	err error
)

// 获取对应的db
func GetDb(dbType string) *gorm.DB {
	var (
		pwd string
		dsn string
		db  *gorm.DB
	)
	if dbType == "mysql" {
		pwd = "mysql@123qaz"
		dsn = fmt.Sprintf("yunpan:%s@tcp(202.115.158.100:3306)/bigdata_dudao?charset=utf8&parseTime=true", pwd)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{TablePrefix: prefix},
		})
	} else if dbType == "oracle" {
		pwd = "XjODS123"
		dsn = fmt.Sprintf("ODS/%s@202.115.154.76:1521/XJSJZXDB", pwd)
		db, err = gorm.Open(oracle.Open(dsn))
	} else if dbType == "mysql2" {
		pwd = "Baiduyun@123"
		dsn = fmt.Sprintf("baiduyun:%s@tcp(192.168.10.47:3306)/bigdata_dudao?charset=utf8&parseTime=true", pwd)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{TablePrefix: prefix},
		})
		if err != nil {
			fmt.Errorf("failed to register callback: %v", err)
			return nil
		}
	}
	db.Logger.LogMode(0)
	return db.Debug()
}

func GetDB() *gorm.DB {
	return DB.Debug()
}
func Init() {
	dbType := config.Instance.Db.DbType
	var dsn string
	if dbType == 0 {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			config.Instance.Db.User, config.Instance.Db.Passwd,
			config.Instance.Db.Host, config.Instance.Db.Port,
			config.Instance.Db.DbName,
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{TablePrefix: config.Instance.Db.Prefix},
		})
		if err != nil {
			fmt.Errorf("failed to register callback: %v", err)
		}
	} else if dbType == 1 {
		dsn = fmt.Sprintf("%s/%s@%s:%d/%s",
			config.Instance.Db.User, config.Instance.Db.Passwd,
			config.Instance.Db.Host, config.Instance.Db.Port,
			config.Instance.Db.DbName,
		)
		DB, err = gorm.Open(oracle.Open(dsn))
	}
	DB.Logger.LogMode(0)

}
func setLogger() logger.Interface {
	// Silent、Error、Warn、Info
	var logMode logger.LogLevel
	logMode = logger.Info

	return logger.Default.LogMode(logMode)
}

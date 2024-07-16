package db

import (
	"fmt"
	"go-study/config"

	"github.com/dzwvip/oracle"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var prefix = "bigdata_"

var (
	DB  *gorm.DB
	err error
)

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
			fmt.Errorf("failed to register callback: %v, dsn:%v", err, dsn)
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

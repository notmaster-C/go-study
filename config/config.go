package config

import (
	"fmt"
	"strings"

	logging "go-study/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	log = logging.Logger("config")
)

// 结构定义从小到大省空间....不是必要

type Db struct {
	ShowSql     bool
	DbType      int //0 mysql 1 oracle
	MaxConn     int
	MaxIdleConn int
	Port        int
	User        string
	Passwd      string
	Host        string
	DbName      string
	Prefix      string
}

type S3 struct {
	Endpoint   string
	External   string
	BucketName string
	Ak         string
	Sk         string
	Ssl        bool
	Version    string
	PrefixPath string
}
type Redis struct {
	Addr       string
	Port       int
	Passwd     string
	Db         int
	MaxRetries int
	PoolSize   int
	Channel    string
}
type Log struct {
	Path    string
	Cutting string
	MaxSize int
	Console bool
	Level   string
}
type ConfigValue struct {
	Db    Db
	S3    S3
	Redis Redis
	Log   Log
}

var Instance = new(ConfigValue)

func Init() {
	v := viper.New()
	// 设置配置文件的名字
	v.SetConfigName("config")
	// 设置配置文件的类型
	v.SetConfigType("yml")
	// 添加配置文件的路径，指定 config 目录下寻找
	v.AddConfigPath("../config")
	// 寻找配置文件并读取
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	//监视配置文件，重新读取配置数据
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	//TODO: 后续优化看能不能直接绑定到结构体
	//var dc Db
	//if err = v.Unmarshal(&dc); err != nil {
	//	fmt.Printf("err:%s", err)
	//}
	//fmt.Println(dc)

	// db
	Instance.Db.User = v.GetString("db.user")
	Instance.Db.Passwd = v.GetString("db.passwd")
	Instance.Db.Host = v.GetString("db.host")
	Instance.Db.Port = v.GetInt("db.port")
	Instance.Db.DbName = v.GetString("db.dbName")
	Instance.Db.Prefix = v.GetString("db.prefix")
	Instance.Db.MaxConn = v.GetInt("db.maxConn")
	Instance.Db.MaxIdleConn = v.GetInt("db.maxIdleConn")
	Instance.Db.ShowSql = v.GetBool("db.showSql")

	// s3
	Instance.S3.Endpoint = v.GetString("s3.endpoint")
	Instance.S3.External = v.GetString("s3.external")
	Instance.S3.BucketName = v.GetString("s3.bucket")
	Instance.S3.Ak = v.GetString("s3.ak")
	Instance.S3.Sk = v.GetString("s3.sk")
	Instance.S3.Ssl = v.GetBool("s3.useSsl")
	Instance.S3.Version = strings.ToUpper(v.GetString("s3.version"))
	Instance.S3.PrefixPath = v.GetString("s3.path")

	// redis
	Instance.Redis.Addr = v.GetString("redis.addr")
	Instance.Redis.Port = v.GetInt("redis.port")
	Instance.Redis.Passwd = v.GetString("redis.passwd")
	Instance.Redis.PoolSize = v.GetInt("redis.poolSize")
	Instance.Redis.MaxRetries = v.GetInt("redis.maxRetries")

	// log
	Instance.Log.Path = v.GetString("log.path")
	Instance.Log.Cutting = v.GetString("log.cutting")
	Instance.Log.MaxSize = v.GetInt("log.maxsize")
	Instance.Log.Console = v.GetBool("log.console")
	Instance.Log.Level = v.GetString("log.level")
}

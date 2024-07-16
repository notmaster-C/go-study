package main

import (
	"go-study/test"
)

func init() {
	// config.Init()
	//todo:常量进行基础赋值初始化
	// constant.Init()
	// 设置日志规则
	// logFile := "go_study.log"
	// logFileName := strings.Join([]string{config.Instance.Log.Path, logFile}, "/")
	// logging.SetupLogFile(
	// 	logFileName,
	// 	config.Instance.Log.Cutting,
	// 	int64(config.Instance.Log.MaxSize),
	// 	config.Instance.Log.Console)
	// logging.SetLogLevel("*", config.Instance.Log.Level)
}

func main() {
	// db.Init()

	// cache.InitRedis(&config.Instance.Redis)
	// route.InitRoute()
	test.TestDefer()
}

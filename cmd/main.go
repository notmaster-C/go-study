package main

import (
	"context"
	"fmt"
	"go-study/algorithm"
	"go-study/config"
	"go-study/db"
	"go-study/file"
	"go-study/utils"
	"time"
)

func init() {
	config.Init()
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
	Test()
}
func Test() {
	fmt.Println("test...")
	//sTime := time.Now()
	algorithm.Test()
	// utils.GithubFlush()
	// db.EtcdWatch()
	//time.Sleep(100)
	//fmt.Println(time.Since(sTime).Seconds() * 1000)
	//fileExam()
	// 运行基准测试并报告结果
	// 创建一个父级 context，设置超时时间为 5 秒钟

	//sql:="set @total:=1,@res:=2;select @total+@res res"
	////找到最后一个分号，前面的sql需要exec，最后一句sql通过raw执行
	//var preSql string
	//if strings.Contains(sql,";")&&strings.Contains(sql,"set"){
	//	split := strings.Split(sql, ";")
	//	preSql = sql[:strings.LastIndex(sql,";")]
	//	sql=split[len(split)-1]
	//}
	////sql=preSql+sql
	//rows, _ := db.GetDB().Exec(preSql).Raw(sql).Rows()
	//if rows!=nil{
	//	defer rows.Close()
	//	res := utils.ScanRows2map(rows)
	//	fmt.Println(res)
	//}

}
func dbTest() {
	// 实测create和createInbatch都是一次insert多条
	// db.GetDB().AutoMigrate(&db.TestFile{})
	file := []*db.TestFile{}
	for i := 0; i < 10; i++ {
		file = append(file, &db.TestFile{
			Name: "test" + utils.ToString(i),
		})
	}
	err := db.GetDB().Model(&db.TestFile{}).CreateInBatches(file, 10).Error
	fmt.Println(err)
}
func isGraduated(xgh string) bool {
	return false
}
func ctxTest() {
	parentCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 创建一个子级 context，用于控制协程
	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel()

	costTime := 5 // 模拟耗时 5 秒钟

	// 启动一个协程
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				// 如果收到取消信号，退出协程
				fmt.Println("协程退出")
				return
			case <-time.After(15 * time.Second):
				fmt.Println("协程超时")
			default:
				timeConsuming(childCtx, costTime)
			}
		}
	}(childCtx)

	// 等待 3 秒钟，然后取消子级 context
	time.Sleep(3 * time.Second)
	fmt.Println("取消协程")
	childCancel()

	// 继续等待 3 秒钟，模拟主协程的一些其他操作
	time.Sleep(3 * time.Second)
	fmt.Println("主协程退出")
}
func timeConsuming(ctx context.Context, costTime int) {

	ctx.Done()

	for i := 1; i <= costTime; i++ {
		// 模拟一些耗时操3
		time.Sleep(5 * time.Second)
		fmt.Printf("协程正在运行第%v次...\n", i)
	}
}
func fileExam() {
	file.DeleteFiles([]string{"F1747902535533854729", "F1747902535533854721", "F1750715924568080385"})
	//now := time.Now()
	//file.HandleDeleteV2([]string{"F1747902535533854729", "F1747902535533854721", "F1750715924568080385"})
	//fmt.Println(time.Since(now))
	//ids, _ := file.DeleteFiles([]string{"F1747902535533854729", "F1747902535533854721", "F1750715924568080385"})
	//fmt.Println(ids)
	//db.GetFileListByPath("F2")
	//file.GetFileList("")
}

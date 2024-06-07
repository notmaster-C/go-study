package main

import (
	"context"
	"fmt"
	"go-study/algorithm"
	"go-study/file"
	"time"
)

func init() {
	//config.Init()
	//db.Init()
	//db.GetDB().AutoMigrate(&db.File{})
	//db.GetDB().Create(&db.Folder{
	//	Name:       "测试文件夹3",
	//	ParentId:   3,
	//	FolderPath: "F2,3",
	//})
	//route.InitRoute()
}

func main() {
	sTime := time.Now()
	algorithm.Test()
	time.Sleep(100)
	fmt.Println(time.Since(sTime).Seconds() * 1000)
	//fileExam()
	// 运行基准测试并报告结果
	// 创建一个父级 context，设置超时时间为 5 秒钟

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

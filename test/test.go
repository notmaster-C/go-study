package test

import (
	"fmt"
	"go-study/file"
)

const str = "test..."

func Test() {
	fmt.Println(str)
	// testGC()
	testChannel()

	//sTime := time.Now()
	// algorithm.Test()
	// utils.GithubFlush()
	// db.EtcdWatch()
	//time.Sleep(100)
	//fmt.Println(time.Since(sTime).Seconds() * 1000)
	//fileExam()

}
func Test2() {

	//var total int
	//for i := 0; i < 5; i++ {
	//	total += i
	//	fmt.Println(total)
	//}
	//s := "abcabcbb"
	//fmt.Println(lengthOfLongestSubstring(s))

	//executeStatDb(1261)
	//gin
	//r := gin.Default()
	//route.InitRoute(r)
	//r.Run()
	//gweb
	//sdb := db.GetDb()
	//var sds []*db.StatDataSource
	//rows := sdb.Where("`key`=?", "testapi").Find(&db.StatDataSource{})
	//rows.Scan(&sds)
	//utils.ExecuteStatApiDb(sds[0])
	//r := gin.Default()
	//route.InitRoute(r)
	//r.Run(":8080")
	// rows, _ := db.GetDb("mysql").Raw("delete from bigdata_middle_data where key_info = 'book_yjs_borrow';insert into bigdata_middle_data select ").Rows()
	// if rows != nil {
	// 	rowse := route.ScanRows2map2(rows)
	// 	jsonString, _ := json.Marshal(rowse)
	// 	result1 := string(jsonString)
	// 	fmt.Println(result1)
	// 	rows.Close()
	// }

	//生成一个字符串长度为8的最随机字符串数组ss，数组大小为len
	//len := 10
	//menus := map[int]string{}
	//fmt.Println(len)
	//rand.Seed(time.Now().UnixNano()) //初始化种子
	//for len > 0 {
	//	i := randomString(8)
	//	fmt.Println("len:", len, i) //打印
	//	menus[len-1] = i
	//	len--
	//}】
	//

	// start := time.Now()
	// wg.Wait()
	// go search("E:\\csl\\", true)
	// waitWork()
	// fmt.Println(matches)
	// fmt.Println(time.Since(start))

	//list := []db.AssitantBooks{}
	//db.GetDb("mysql").Model(&db.AssitantBooks{}).Limit(1).Find(&list)
	//info, _ := json.Marshal(list)
	//infos := string(info)
	//fmt.Println(infos)
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

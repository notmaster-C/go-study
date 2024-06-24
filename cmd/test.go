package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

// SuperviseCourseListenTemp 督导课程听课信息临时处理表
type SuperviseCourseListenTemp struct {
	Id             uint64 `gorm:"type:bigint(32);autoIncrement;primarykey" json:"id"`
	ClassNumber    string `gorm:"type:varchar(255);index" json:"classNumber"`                // 开课号
	TeacherNumber  string `gorm:"type:varchar(64);index" json:"teacherNumber"`               // 任课老师工号
	ListenTime     string `gorm:"type:varchar(255);not null" json:"listenTime"`              // 听课时间
	AddressCode    string `gorm:"type:varchar(255);not null" json:"AddressCode"`             // 听课地点
	Weeks          string `gorm:"type:varchar(255);not null" json:"weeks"`                   // 周次
	Week           int    `gorm:"type:varchar(255);not null" json:"week"`                    // 星期几
	Sessions       string `gorm:"type:varchar(100);not null;" json:"sessions"`               // 节次
	WeeksBinary    int64  `gorm:"type:varchar(255);not null;index" json:"weeksBinary"`       // 周次段
	SessionsBinary int64  `gorm:"type:varchar(100);not null;index" json:"sessionsBinary"`    // 节次段
	Status         bool   `gorm:"type:tinyint(1);not null;default:0;" json:"status"`         // 更新状态
	StartTime      int64  `gorm:"type:int(11);index:idx_address_start_end" json:"startTime"` // 听课开始时间
	EndTime        int64  `gorm:"type:int(11);index:idx_address_start_end" json:"endTime"`   // 听课结束时间
	Campus         string `gorm:"type:varchar(100);not null;" json:"campus"`                 // 校区
	TotalStudent   int64  `gorm:"type:int(64);" json:"totalStudent"`                         // 总学生
}

// 督导课程听课信息
type SuperviseCourseListen struct {
	ClassNumber   string `gorm:"type:varchar(255);index" json:"classNumber"`                  // 开课号
	TeacherNumber string `gorm:"type:varchar(64);index" json:"teacherNumber"`                 // 任课老师工号
	AddressCode   string `gorm:"type:int(11);index:idx_address_start_end" json:"addressCode"` // 听课地点编码
	Campus        string `gorm:"type:varchar(255)" json:"campus"`                             // 所属校区
	StartTime     int64  `gorm:"type:int(11);index:idx_address_start_end" json:"startTime"`   // 听课开始时间
	EndTime       int64  `gorm:"type:int(11);index:idx_address_start_end" json:"endTime"`     // 听课结束时间
	TotalStudent  int64  `gorm:"type:int(64);" json:"totalStudent"`                           // 总学生
}

const layout = "2006-01-02 15:04:05"
const (
	batchSize          = 1000
	firstDate          = 1708876800
	SCLParseStatusWait = iota
	SCLParseStatusDone
)

func getCurrentWeek() int {
	p := time.Now().Unix() - firstDate
	s := float64(p) / (24 * 60 * 60 * 7)
	return int(math.Ceil(s))
}

var wg sync.WaitGroup
var mutex sync.Mutex
var matches int
var query = "theme"

var workNum = 0
var maxWorkNum = 16

var searchRequest = make(chan string)
var foundMatch = make(chan bool)
var workDone = make(chan bool)

func lengthOfLongestSubstring(s string) int {

	rk, ans := -1, 0
	n := len(s)
	m := map[byte]int{}

	for i := 0; i < n; i++ {
		if i != 0 {
			delete(m, s[i-1])
		}
		for rk+1 < n && m[s[rk+1]] == 0 {
			val, ok := m[s[rk+1]]
			if !ok {
				// 如果该键不存在于映射中，则需要初始化
				m[s[rk+1]] = 1
			} else {
				// 如果键存在，则递增其值
				m[s[rk+1]] = val + 1
			}

			rk++
		}
		ans = max(ans, rk+1-i)
	}
	return ans
}

func test() {

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

	start := time.Now()
	wg.Wait()
	go search("E:\\csl\\", true)
	waitWork()
	fmt.Println(matches)
	fmt.Println(time.Since(start))

	//list := []db.AssitantBooks{}
	//db.GetDb("mysql").Model(&db.AssitantBooks{}).Limit(1).Find(&list)
	//info, _ := json.Marshal(list)
	//infos := string(info)
	//fmt.Println(infos)

}
func waitWork() {
	for {
		select {
		case <-foundMatch:
			matches++
		case path := <-searchRequest:
			go search(path, true)
		case <-workDone:
			workNum--
			if workNum <= 0 {
				return
			}
		}
	}
}
func search(path string, master bool) {

	files, err := os.ReadDir(path)
	if err == nil {
		for _, file := range files {
			name := file.Name()
			if strings.Contains(name, query) {
				foundMatch <- true
			}
			if file.IsDir() {
				if workNum < maxWorkNum {
					workNum++
					searchRequest <- path + name + "\\"
				}
				search(path+file.Name()+"\\", false)
			}
		}
	}
	if master {
		workDone <- true
		return
	}

}

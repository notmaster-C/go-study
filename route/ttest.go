package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-study/db"
	"go-study/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-resty/resty/v2"
)

type MSG struct {
	Id    int64  `json:"id" form:"id" xml:"id"`
	Name  string `json:"name" form:"name" xml:"name"`
	Value []byte `json:"value" form:"value" xml:"value"`
}

func handleTest(c *gin.Context) {

	req := resty.New().R()
	req.SetHeader("Authorization", c.Request.Header.Get("Authorization"))
	req.SetHeader("Content-Type", "application/json;charset=utf-8")
	req.SetBody("{\n    \"business_types\": \"BIZ\",\n    \"current\": 1,\n    \"dbTypes\": \"all\",\n    \"dept_ids\": \"all\",\n    \"manufacturer_ids\": \"all\",\n    \"name\": null,\n    \"size\": 1000,\n    \"sort_rule\": \"\",\n    \"state\": \"all\",\n    \"system_ids\": \"all\"\n}")
	resq, _ := req.Post("http://202.115.158.71:8770/warehouse/data_source/list_back_database/")
	var data db.XJT1
	err := json.Unmarshal(resq.Body(), &data)
	if err != nil {
		return
	}
	qbd := db.GetDB()
	//qbd.AutoMigrate(&db.ApiDbMiddleSource{})
	qbd.AutoMigrate(&db.Records{})

	for _, record := range data.Result.Records {
		db.GetDB().Table("Records").Create(record)
	}
	c.JSON(http.StatusOK, gin.H{"Result:": data.Result.Records[0]})

}
func handleTest5(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer func() {
		if err != nil {
			res["content"] = err
		}
		c.JSON(http.StatusOK, res)
	}()
	//form, err := c.MultipartForm()
	//if err != nil {
	//	return
	//}
	//parentId := form.Value["parentId"]
	var total int
	for i := 0; i < 5; i++ {
		total += i
	}
	res["content"] = total
}

type NetlogRecord struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"Id"`
	RecordInfo string    `gorm:"column:RecordInfo;" json:"RecordInfo"`
	CreateTime time.Time `gorm:"column:createtime;type:datetime(0);autoUpdateTime" json:"createtime"`
	UpdateTime time.Time `gorm:"column:updatetime;type:datetime(0);autoUpdateTime" json:"updatetime"`
}

//	func handleTest5(c *gin.Context) {
//		var (
//			err error
//			res = gin.H{}
//		)
//		defer func() {
//			c.JSON(http.StatusOK, res)
//		}()
//		var result string
//		sql := "select param from bigdata_stat_data_sources where id=1480;select param from bigdata_stat_data_sources where id=1479"
//		err = db.GetDb().Raw(sql).Scan(&result).Error
//		//err = db.GetDb().Raw(result).Scan(&result).Error
//		if err != nil {
//			fmt.Errorf("failed to statId:%v error:%v", sql, err)
//			res["sql"] = sql
//			return
//		}
//		res["sql"] = sql
//		res["result"] = result
//	}
func handleTest4(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer func() {
		c.JSON(http.StatusOK, res)
	}()
	dir := "E:\\csl\\netlog\\"  // 替换成你的目录路径
	pattern := "_online_detail" // 替换成你的文件名格式
	var str string
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		str = fmt.Sprintf("无法读取文件:", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		str = fmt.Sprintf("找不到符合条件的文件")
		os.Exit(0)
	}

	str = fmt.Sprintf("匹配到的文件:")
	for _, file := range files {
		str += fmt.Sprintf(file)
	}
	if err != nil {
		return
	}
	res["res"] = str
}
func handleTest2(c *gin.Context) {
	type DrillDownParams struct {
		Key         string   `form:"Key" json:"Key"`
		Fields      []string `form:"Fields" json:"Fields" xml:"Fields"`
		Groups      []string `form:"Groups" json:"Groups" xml:"Groups"`
		XueyuanIds  []string `form:"xueyuanIds" json:"xueyuanIds" xml:"xueyuanIds"`    // 学院代码 就是学院的code 数组
		NianjiName  []string `form:"nianjiName" json:"nianjiName" xml:"nianjiName"`    // 年级名称 比如2022
		ZhuanyeName []string `form:"zhuanyeName" json:"zhuanyeName" xml:"zhuanyeName"` // 专业代码 专业code数组
		Bjm         []string `form:"bjm" json:"bjm" xml:"bjm"`
		Dtype       int      `form:"dtype" json:"dtype" xml:"dtype"` //班级编码code 数组 0非查询 1查询非下钻
	}
	var (
		err    error
		result = gin.H{}
		rows   *sql.Rows
		sdb    = db.GetDB()
		sql    string
	)
	defer func() {
		if err != nil {
			result["error"] = err.Error()
			OutputReuslt(c, result)
		}
		if rows != nil {
			res := ScanRows2map2(rows)
			jsonString, _ := json.Marshal(res)
			result1 := string(jsonString)
			result["success"] = true
			result["code"] = 200
			result["content"] = result1
			fmt.Println(sql)
			rows.Close()
			OutputReuslt(c, result)
		}
	}()
	var req *DrillDownParams
	err = c.ShouldBindJSON(&req)
	// 构建 SQL 查询语句
	if len(req.Fields) == 0 {
		req.Fields = []string{"sum(result) as total", "count(DISTINCT ex3) as count1", "SUM(ex5) as sum2", "ex2", "ex2_name"}
	}
	if len(req.Groups) == 0 {
		req.Groups = []string{"ex2", "ex3", "ex4"}
	}
	// result ex1 前两列做参数返回或计算操作 ex2学院 ex3专业 ex4班级  三级下钻 后续需要再增加
	if req.Dtype == 0 {
		if len(req.XueyuanIds) == 0 {
			sql = "SELECT  " + strings.Join(req.Fields, ",") + " from bigdata_middle_data  WHERE key_info= ? GROUP BY " + req.Groups[0]
			rows, err = sdb.Raw(sql, req.Key).Rows()
			return
		}
		if len(req.XueyuanIds) > 0 && len(req.ZhuanyeName) == 0 {
			sql = "SELECT  " + strings.Join(req.Fields, ",") + " from bigdata_middle_data  WHERE key_info=? and ex2 in (?)  GROUP BY " + req.Groups[1]
			rows, err = sdb.Raw(sql, req.Key, req.XueyuanIds).Rows()
			return
		}
		if len(req.XueyuanIds) > 0 && len(req.ZhuanyeName) > 0 {
			sql = "SELECT  " + strings.Join(req.Fields, ",") + " from bigdata_middle_data  WHERE key_info=? and ex3 in (?)  GROUP BY " + req.Groups[2]
			rows, err = sdb.Raw(sql, req.Key, req.ZhuanyeName).Rows()
			return
		}
	}
	return
}

// 输出结果
func OutputReuslt(c *gin.Context, result gin.H) {
	if result == nil {
		result = gin.H{}
	}
	c.JSON(http.StatusOK, result)
}
func handleGet(c *gin.Context) {
	c.JSON(http.StatusOK, fmt.Sprintf("%d%s%s", time.Now().Year(), time.Now().Format("01"), time.Now().Format("02")))
}
func handleSetSswwg(c *gin.Context) {

	var req *db.SSWWG
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fmt.Printf("[xx] bind common_param error err=%+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sswwgt *db.SSWWG
	db.GetDB().Model(&db.SSWWG{}).Where("id=?", 1).First(&sswwgt)
	if req.OutbeginDate != "" {
		sswwgt.OutbeginDate = req.OutbeginDate
	}
	if req.OutendDate != "" {
		sswwgt.OutendDate = req.OutendDate
	}
	if req.LateendDate != "" {
		sswwgt.LateendDate = req.LateendDate
	}
	if req.LatebeginDate != "" {
		sswwgt.LatebeginDate = req.LatebeginDate
	}
	if req.StayHour != "" {
		sswwgt.StayHour = req.StayHour
	}
	if req.StayRate != "" {
		sswwgt.StayRate = req.StayRate
	}
	if req.OutRate != "" {
		sswwgt.OutRate = req.OutRate
	}
	if req.LateRate != "" {
		sswwgt.LateRate = req.LateRate
	}
	if req.Id != 0 {
		sswwgt.Id = req.Id
	}
	//if db.InitConfig(sswwgt) == 0 {
	//	c.JSON(http.StatusOK, gin.H{"message": "未做任何修改", "object": &sswwgt})
	//} else {
	//	c.JSON(http.StatusOK, gin.H{"message": sswwgt})
	//}

}
func ScanRows2map2(rows *sql.Rows) (res []map[string]interface{}) {
	defer rows.Close()
	cols, _ := rows.Columns()
	cache := make([]interface{}, len(cols))
	// 为每一列初始化一个指针
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	for rows.Next() {
		rows.Scan(cache...)
		row := make(map[string]interface{})
		for i, val := range cache {
			// 处理数据类型
			v := *val.(*interface{})
			switch v.(type) {
			case []uint8:
				v = string(v.([]uint8))
			case nil:
				v = ""
			}
			row[cols[i]] = v
		}
		res = append(res, row)
	}
	rows.Close()
	return res
}

// 宿舍晚未归接口-获取对应学号学生异常详细信息
func hanleGetSSWWGXY(c *gin.Context) {
	var (
		err    error
		result = gin.H{}
	)

	defer func() {
		c.JSON(http.StatusOK, result)
	}()
	var req db.SSWWG11
	if err = c.ShouldBindJSON(&req); err != nil {
		result["message"] = "err"
		return
	}
	var rows *sql.Rows
	defer func() {
		if rows != nil {
			res := ScanRows2map2(rows)
			jsonString, _ := json.Marshal(res)
			result1 := string(jsonString)
			result["content"] = result1
			rows.Close()
		}
	}()
	rows, _ = db.GetDB().Table("bigdata_middle_data").Where("key_info=? and ex8_name=?", req.Key).Group("ex3").Rows()
	//-- `key` result年级 ex1 进出异常  ex2 学号姓名 ex3 学院 ex4 校区 ex5 所住苑区 ex6 所住楼栋 ex7 所住房间号 ex8日期,ex8_name 是否研究生,ex9 学院总人数
	if len(req.Xh) == 0 {
		//这一判断表示非请求用户详细信息
		if len(req.XueyuanId) == 0 {
			if len(req.StartTime) == 0 {
				rows, _ = db.GetDB().Table("bigdata_middle_data").Where("key_info=? and ex8_name=?", req.Key, req.StuType).Group("ex3").Rows()
				//sql := "SELECT sum(ex1) as in_all,sum(ex1_name) as out_all,ex3_name as 'name',ex3 as 'code' from bigdata_middle_data WHERE key_info=? and ex8_name=? GROUP BY ex3"
				//rows, _ := db.GetDB().Raw(sql, req.Key, req.StuType).Rows()
				return
			} else if len(req.EndTime) == 0 {
				sql := "SELECT sum(ex1) as in_all,sum(ex1_name) as out_all,ex3_name as 'name',ex3 as 'code' from bigdata_middle_data WHERE key_info=? and ex8_name=? and ex8 =? GROUP BY ex3"
				rows, _ = db.GetDB().Raw(sql, req.Key, req.StuType, req.StartTime).Rows()
				return
			}
			if len(req.NianjiName) > 0 {
				sql := "SELECT sum(ex1) as in_all,sum(ex1_name) as out_all,ex3_name as 'name',ex3 as 'code' from bigdata_middle_data WHERE key_info=? and ex8_name=? and ex8 between ? and ? GROUP BY ex3"
				rows, _ = db.GetDB().Raw(sql, req.Key, req.StuType, req.StartTime, req.EndTime).Rows()
				return
			}

		}

	}

}

func handleFileToDb(c *gin.Context) {
	utils.Excute()
}

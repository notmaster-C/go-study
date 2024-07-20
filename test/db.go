package test

import (
	"fmt"
	"go-study/db"
	"go-study/utils"
	"strings"
)

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

func runSet() {
	sql := "set @total:=1,@res:=2;select @total+@res res"
	//找到最后一个分号，前面的sql需要exec，最后一句sql通过raw执行
	var preSql string
	if strings.Contains(sql, ";") && strings.Contains(sql, "set") {
		split := strings.Split(sql, ";")
		preSql = sql[:strings.LastIndex(sql, ";")]
		sql = split[len(split)-1]
	}
	//sql=preSql+sql
	rows, _ := db.GetDB().Exec(preSql).Raw(sql).Rows()
	if rows != nil {
		defer rows.Close()
		res := utils.ScanRows2map(rows)
		fmt.Println(res)
	}

}

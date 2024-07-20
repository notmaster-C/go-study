package test

import (
	"fmt"
	"go-study/db"
	"go-study/utils"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
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

type Optimistic struct {
	Id      int64   `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserId  string  `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 用户ID
	Amount  float32 `gorm:"column:amount;NOT NULL" json:"amount"`             // 金额
	Version int64   `gorm:"column:version;default:0;NOT NULL" json:"version"` // 版本
}

func (o *Optimistic) GetVersion() int64 {
	return o.Version
}

func (o *Optimistic) SetVersion(version int64) {
	o.Version = version
}
func testCAS() {
	db := db.GetDB()
	var out Optimistic
	db.AutoMigrate(Optimistic{})
	db.Create(Optimistic{UserId: "123", Amount: 12.3})
	db.First(&out, Optimistic{Id: 1})
	out.Amount = out.Amount + 10
	column := db.Model(&out).Where("id", out.Id).Where("version", out.Version).
		UpdateColumn("amount", out.Amount).
		UpdateColumn("version", gorm.Expr("version+1"))
	fmt.Printf("#######update %v line \n", column.RowsAffected)

}
func UpdateWithOptimistic(db *gorm.DB, model db.Lock, callBack func(model db.Lock) db.Lock, retryCount, currentRetryCount int32) (err error) {
	if currentRetryCount > retryCount {
		return errors.New("Maximum number of retries exceeded:" + strconv.Itoa(int(retryCount)))
	}
	currentVersion := model.GetVersion()
	model.SetVersion(currentVersion + 1)
	column := db.Model(model).Where("version", currentVersion).UpdateColumns(model)
	affected := column.RowsAffected
	if affected == 0 {
		if callBack == nil && retryCount == 0 {
			return errors.New("Concurrent optimistic update error")
		}
		time.Sleep(100 * time.Millisecond)
		db.First(model)
		bizModel := callBack(model)
		currentRetryCount++
		err := UpdateWithOptimistic(db, bizModel, callBack, retryCount, currentRetryCount)
		if err != nil {
			return err
		}
	}
	return column.Error

}
func BenchmarkUpdateWithOptimistic(b *testing.B) {
	tx := db.GetDB()
	b.RunParallel(func(pb *testing.PB) {
		var out Optimistic
		tx.First(&out, Optimistic{Id: 1})
		out.Amount = out.Amount + 10
		err := UpdateWithOptimistic(tx, &out, func(model db.Lock) db.Lock {
			bizModel := model.(*Optimistic)
			bizModel.Amount = bizModel.Amount + 10
			return bizModel
		}, 3, 0)
		if err != nil {
			fmt.Printf("%+v \n", err)
		}
	})
}

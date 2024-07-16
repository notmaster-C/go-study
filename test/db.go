package test

import (
	"fmt"
	"go-study/db"
	"go-study/utils"
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

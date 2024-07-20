package db

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Lock interface {
	SetVersion(version int64)
	GetVersion() int64
}

func UpdateWithOptimistic(db *gorm.DB, model Lock, callBack func(model Lock) Lock, retryCount, currentRetryCount int32) (err error) {
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

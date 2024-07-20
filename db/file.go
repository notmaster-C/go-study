package db

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TestFile struct {
	gorm.Model
	Name       string `gorm:"type:varchar(512);not null;index" json:"name"` // 文件名称
	FolderPath string `gorm:"type:varchar(768);index" json:"folderPath"`    // 文件路径
}
type File struct {
	Id         string         `gorm:"type:varchar(64);primarykey"  json:"id"`         // 文件id
	Name       string         `gorm:"type:varchar(512);not null;index" json:"name"`   // 文件名称
	Size       int64          `gorm:"type:int(32);default 0" json:"size"`             // 文件大小
	Hash       string         `gorm:"type:varchar(255)" json:"-"`                     // 文件唯一标识
	ParentId   string         `gorm:"type:varchar(64);index" json:"parentId"`         // 父级Id
	IsFolder   bool           `gorm:"type:tinyint(3);default 0" json:"isFolder"`      // 是否目录
	MineType   string         `gorm:"type:varchar(256)" json:"mineType"`              // 文件类型
	FolderPath string         `gorm:"type:varchar(768);index" json:"folderPath"`      // 文件路径
	DeletedAt  gorm.DeletedAt `gorm:"type:datetime;index" json:"deletedAt,omitempty"` // 删除时间 Unscoped()进行硬删除
	CreatedAt  time.Time      `gorm:"type:datetime" json:"createdAt,omitempty"`       // 创建时间
	UpdatedAt  time.Time      `gorm:"type:datetime" json:"updatedAt,omitempty"`       // 更新时间
}

func GetFileListByIds(ids []string) []*File {
	if DB == nil {
		return nil
	}

	files := []*File{}
	err = GetDB().
		Model(files).
		Where(fmt.Sprintf("id IN (\"%s\")", strings.Join(ids, "\",\""))).
		Find(&files).Error
	if err != nil {
		fmt.Errorf("Query files %+v failed with %+v", ids, err)
		return nil
	}

	return files
}

func GetFileById(id string) *File {
	if DB == nil {
		return nil
	}

	file := File{}
	err = GetDB().
		Model(file).
		Where("id = ?", id).
		First(&file).Error
	if err != nil {
		fmt.Errorf("Query file %v failed with %+v", id, err)
		return nil
	}

	return &file
}

func GetFileListByParent(id string) []*File {

	files := []*File{}
	err = GetDB().
		Model(files).
		Where("parent_id = ?", id).
		Find(&files).Error
	if err != nil {
		fmt.Errorf("no such parent files %+v failed with %+v", id, err)
		return nil
	}

	return files
}
func GetFileListByParentIds(id []string) []*File {

	files := []*File{}
	err = GetDB().
		Model(files).
		Where("parent_id in ?", id).
		Find(&files).Error
	if err != nil {
		fmt.Errorf("no such parent files %+v failed with %+v", id, err)
		return nil
	}

	return files
}
func GetFileOrder(order int64) string {
	switch order {
	case 1:
		return "created_at"
	case 2:
		return "size desc"
	case 3:
		return "size"
	default:
		return "created_at desc"
	}
}
func GetFilesByPath(path string) (file []*File) {

	err = GetDB().Where("folder_path like concat(?,'%')", path).
		Find(&file).Error

	if err != nil {
		fmt.Errorf("no such files %+v failed with %+v", path, err)
		return
	}
	return
}
func DeleteFiles(ids []string) error {

	deleteStr := fmt.Sprintf("id IN (\"%s\")", strings.Join(ids, "\",\""))
	err = GetDB().Delete(&File{}, deleteStr).Error
	if err != nil {
		return err
	}
	return nil
}

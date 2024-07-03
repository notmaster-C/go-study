package file

import (
	"fmt"
	"go-study/common"
	"go-study/db"
	"go-study/oss"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(inFileIds []string) {

}

// 递归批量删除文件
func DelFiles(inFileIds []string) {
	sTime := time.Now()
	files := db.GetFileListByIds(inFileIds)
	fileIds := []string{}
	fileObjects := []string{}
	for _, file := range files {
		fileIds = append(fileIds, file.Id)
		if file.IsFolder {
			ids, oid := deleteFolder(file.Id)
			fileIds = append(fileIds, ids...)
			fileObjects = append(fileObjects, oid...)
			continue
		}
		filePath := oss.GetObjectPath(fmt.Sprintf("%s/%s/%s", common.ATTACHMENT_FILE_PATH_PREFIX, generatePath(file.Hash), file.Id))
		fileObjects = append(fileObjects, filePath)
	}

	fmt.Println("常规删除：", fileIds)
	// 删除元数据
	//err := db.GetDB().Delete(&db.File{}, fmt.Sprintf("id IN (\"%s\")", strings.Join(fileIds, "\",\""))).Error
	//if err != nil {
	//	return
	//}
	fmt.Println("cost:", time.Since(sTime))
	//TODO:数据文件夹删除- oss文件夹删除流程....
}

// 队列删除
func DeleteFilesQueue(fIds []string) (fileIds []string, fileObjects []string) {
	sTime := time.Now()

	stack := fIds

	first := true
	for len(stack) > 0 {
		var files []*db.File
		if first {
			files = db.GetFileListByIds(stack)
			first = false
		} else {
			files = db.GetFileListByParentIds(stack)
		}
		stack = []string{}

		for _, file := range files {

			fileIds = append(fileIds, file.Id)
			if file.IsFolder {
				// 将子文件夹推入栈
				stack = append(stack, file.Id)
				continue
			}
			filePath := oss.GetObjectPath(fmt.Sprintf("%s/%s/%s", common.ATTACHMENT_FILE_PATH_PREFIX, generatePath(file.Hash), file.Id))
			fileObjects = append(fileObjects, filePath)
		}
	}
	fmt.Println("cost:", time.Since(sTime))
	return fileIds, fileObjects
}

// 递归删除文件夹下所有文件以及文件夹
func deleteFolder(folderId string) (fileIds []string, fileObjects []string) {
	files := db.GetFileListByParent(folderId)
	for _, file := range files {
		fileIds = append(fileIds, file.Id)
		if file.IsFolder {
			deleteFolder(file.Id)
			continue
		}
		filePath := oss.GetObjectPath(fmt.Sprintf("%s/%s/%s", common.ATTACHMENT_FILE_PATH_PREFIX, generatePath(file.Hash), file.Id))
		fileObjects = append(fileObjects, filePath)
	}
	fmt.Println("递归删除：", fileIds)
	return
}

// 生成路径
func generatePath(index string) string {
	// 16/85/96/28/1685962887317684224
	str := []string{}
	for i := 0; i <= 6; i += 2 {
		str = append(str, index[i:i+2])
	}
	str = append(str, index)

	return filepath.Join(str...)
}

func handleUpload(c *gin.Context) {
	var (
		err    error
		result = gin.H{}
	)
	defer func() {
		common.OutputResult(c, err, result)
	}()

	form, err := c.MultipartForm()
	if err != nil {
		return
	}
	parentId := form.Value["parentId"][0]
	uploadFiles := form.File["file"]
	files := []db.File{}
	parent := db.File{}
	//处理文件夹内上传
	if parentId != "" && parentId != "/" {
		err = db.GetDB().Where(db.File{Id: parentId}).First(&parent).Error
		if err != nil {
			result["content"] = "folder not exists"
			return
		}
	}
	var totalSize int64
	for _, file := range uploadFiles {
		// 处理上传逻辑
		//f, err1 := file.Open()
		//if err1 != nil {
		//	fmt.Errorf("upload %s to s3 with open error %s", file.Filename, err1)
		//	f.Close()
		//	continue
		//}
		//
		//// 获取content-type
		//contentType := file.Header.Get("Content-Type")
		//if contentType == "" {
		//	buffer := make([]byte, 512)
		//	if _, err1 := f.Read(buffer); err1 == nil {
		//		contentType = http.DetectContentType(buffer)
		//	}
		//	f.Seek(0, 0)
		//}
		//
		//hash := utils.GetSnowflakeString()
		//fid := "F" + utils.GetSnowflakeString()
		//filePath := fmt.Sprintf("%s/%s/%s", common.ATTACHMENT_FILE_PATH_PREFIX, generatePath(hash), fid)
		//objectPath := oss.GetObjectPath(filePath)
		//upInfo, err1 := oss.NewOssHandle().NewFiles().PutObject(objectPath, file.Size, f, contentType, privateMetaData)
		//if err1 != nil {
		//	fmt.Errorf("upload %s failed with %v", objectPath, err1)
		//	f.Close()
		//	continue
		//}
		//fmt.Errorf("upload %s to s3 with etag %s", objectPath, upInfo)
		//f.Close()
		//files = append(files, db.File{
		//	Id:       fid,
		//	Name:     file.Filename,
		//	Size:     file.Size,
		//	Hash:     hash,
		//	ParentId: parentId,
		//	MineType: contentType,
		//})
		files = append(files, db.File{
			Id:       "123",
			Name:     file.Filename,
			Size:     file.Size,
			Hash:     "hash",
			ParentId: parentId,
			MineType: "contentType",
		})
		totalSize += file.Size
	}

	err = db.GetDB().Create(files).Error
	if err != nil {
		return
	}
	//处理文件夹内上传
	if parentId != "" && parentId != "/" {
		parent.Size += totalSize
		err = db.GetDB().Updates(&parent).Error
		if err != nil {
			return
		}
	}

	result["content"] = files
}

// 文件列表
func GetFileList(parentId string) {

	var (
		list  []db.File
		total int64
	)
	order := db.GetFileOrder(1)
	if parentId == "" {
		parentId = "/"
	}
	err := db.GetDB().Model(db.File{}).
		Where(db.File{
			ParentId: parentId,
		}).
		Order(order).Count(&total).Limit(10).
		Find(&list).Error
	if err != nil {
		return
	}
	fmt.Println(list)

}

// 多线程删除
func DeleteFiles(ids []string) (err error) {

	files := db.GetFileListByIds(ids)
	if len(files) == 0 {
		fmt.Println("no such files")
		return
	}
	var wg sync.WaitGroup
	fileObjects := make([]string, len(files))
	fileIds := make([]string, len(files))
	for _, file := range files {
		fileIds = append(fileIds, file.Id)
		if file.IsFolder {
			folderPath := file.FolderPath + "/" + file.Id
			//获取该路径下所有文件
			folderFiles := db.GetFilesByPath(folderPath)
			//没有文件跳过
			if len(folderFiles) == 0 {
				continue
			}
			wg.Add(1)
			go delFolderFile(folderFiles, &wg)
			continue
		}
		filePath := oss.GetObjectPath(fmt.Sprintf("%s/%s/%s", common.ATTACHMENT_FILE_PATH_PREFIX, generatePath(file.Hash), file.Id))
		fileObjects = append(fileObjects, filePath)
	}

	// 删除元数据
	if err = db.DeleteFiles(fileIds); err != nil {
		return
	}
	// 删除流数据
	if err = delOssObjBatches(fileObjects, &wg, 100); err != nil {
		return
	}
	wg.Wait()
	return nil
}

/*
@param deleteBatchSize 批量删除的数量
*/
func delOssObjBatches(fileObjects []string, wg *sync.WaitGroup, deleteBatchSize int) error {
	// 计算批次数
	batchCount := (len(fileObjects) + deleteBatchSize - 1) / deleteBatchSize
	wg.Add(batchCount)
	// 分批次多线程删除流数据
	for i := 0; i < batchCount; i++ {
		start := i * deleteBatchSize
		end := start + deleteBatchSize
		if end > len(fileObjects) {
			end = len(fileObjects)
		}

		batch := fileObjects[start:end]
		go func(batch []string) {
			defer wg.Done()
			fmt.Println(batchCount, "batch:", batch)
			//res := oss.NewOssHandle().NewFiles().RemoveObjects(batch)
			//if res != nil {
			//	fmt.Printf("Failed to delete object %s: %v\n", batch, res)
			//	return
			//}
		}(batch)

	}

	return nil
}

// 处理文件夹的删除
func delFolderFile(files []*db.File, wg *sync.WaitGroup) {
	defer wg.Done()
	fileObjects := make([]string, len(files))
	fileIds := make([]string, len(files))
	for _, file := range files {
		fileIds = append(fileIds, file.Id)
		if file.IsFolder {
			//文件夹跳过,已被选中
			continue
		}
		filePath := oss.GetObjectPath(fmt.Sprintf("%s/%s/%s", common.ATTACHMENT_FILE_PATH_PREFIX, generatePath(file.Hash), file.Id))
		fileObjects = append(fileObjects, filePath)
	}
	// 删除元数据
	if err := db.DeleteFiles(fileIds); err != nil {
		return
	}
	// 删除流数据
	if err := delOssObjBatches(fileObjects, wg, 100); err != nil {
		return
	}
	return
}

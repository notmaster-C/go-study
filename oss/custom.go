package oss

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"io/ioutil"
)

type custom struct {
	oss ossHandle
}

//自定义扩展
func (oss *ossHandle) NewCustom() *custom {
	return &custom{
		oss: *oss,
	}
}

//根据Bytes计算filehash
func (cus *custom) FileHashBytesSum(fileByte []byte) (hex string, err error) {
	h := sha1.New()
	h.Write(fileByte)
	return fmt.Sprintf("%X", h.Sum(nil)), nil
}

//根据Render计算filehash
func (cus *custom) FileHashRenderSum(render io.Reader) (hex string, err error) {
	b, err := ioutil.ReadAll(render)
	if err != nil {
		return
	}
	return cus.FileHashBytesSum(b)
}

//生成key
func (cus *custom) GeneratePath(fileHash string) (path string) {
	path = cus.oss.path
	for i := 0; i < 4; i++ {
		path += fileHash[2*i:2*(i+1)] + "/"
	}
	return path + fileHash
}

//检测文件是否存在
func (cus *custom) CheckOssFile(key string) (stat minio.UploadInfo, err error) {
	statOpts := minio.StatObjectOptions{}
	res, err := cus.oss.client.StatObject(context.Background(), cus.oss.bucketName, key, statOpts)
	if res.Size <= 0 {
		err = errors.New(fmt.Sprintf("file size error,key:%v,err:%v", key, stat.Size))
		return
	}
	stat.Size = res.Size
	stat.Key = res.Key
	stat.ETag = res.ETag
	return
}

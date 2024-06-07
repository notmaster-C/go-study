package oss

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"go-study/config"
	"io"
	"net/http"
	"net/url"
	"time"
)

type files struct {
	oss ossHandle
}

// 文件操作
func (oss *ossHandle) NewFiles() *files {
	return &files{
		oss: *oss,
	}
}

// 获取流文件
func (f *files) GetObject(key string) ([]byte, error) {
	statOpts := minio.GetObjectOptions{}
	obj, err := f.oss.client.GetObject(context.Background(), f.oss.bucketName, key, statOpts)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = obj.Close()
	}()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 文件流上传
func (f *files) PutObject(key string, fileSize int64, fileStream io.Reader, ContentType string, userMetaData map[string]string) (uploadInfo minio.UploadInfo, err error) {
	//文件检测
	uploadInfo, err = f.oss.NewCustom().CheckOssFile(key)
	if err == nil && uploadInfo.Size == fileSize {
		return
	}
	//文件上传
	opts := minio.PutObjectOptions{}
	if ContentType != "" {
		opts.ContentType = ContentType
	}
	opts.UserMetadata = userMetaData
	// make it public
	// userMetaData := map[string]string{"x-amz-acl": "public-read"}
	uploadInfo, err = f.oss.client.PutObject(context.Background(), config.Instance.S3.BucketName, key, fileStream, fileSize, opts)
	return
}

// 获取文件url
// expires地址有效期，单位为秒，最长为604800（即 7 天），最短为 1 秒。
func (f *files) PreSignedGetObject(key, fileName string, expires int) (*url.URL, error) {
	if expires <= 0 {
		expires = 7 * 24 * 60 * 60
	}
	reqParams := make(url.Values)
	// 以附件的方式下载
	if fileName != "" {
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	}
	return f.oss.client.PresignedGetObject(context.Background(), f.oss.bucketName, key, time.Duration(expires)*time.Second, reqParams)
}

// 预签名
// expires地址有效期，单位为秒，最长为604800（即 7 天），最短为 1 秒。
func (f *files) PreSign(objectName string, expires int) (*url.URL, error) {
	if expires <= 0 {
		expires = 7 * 24 * 60 * 60
	}
	return f.oss.client.Presign(context.Background(), http.MethodPost, f.oss.bucketName, objectName, time.Duration(expires)*time.Second, url.Values{})
}

// 调用CopyObject接口拷贝同一地域下相同或不同存储空间（Bucket）之间的文件（Object）
func (f *files) CopyObject(destBucket, destKey, srcBucket, srcKey string) (minio.UploadInfo, error) {
	return f.oss.client.CopyObject(context.Background(),
		minio.CopyDestOptions{
			Bucket: destBucket,
			Object: destKey,
		},
		minio.CopySrcOptions{
			Bucket: srcBucket,
			Object: srcKey,
		},
	)

}

// 本地文件上传
// FPutObject 在单个 PUT 操作中上传小于 128MiB 的对象。对于大于 128MiB 的对象，
// 会根据实际文件大小以 128MiB 或更大的块无缝上传对象。对象的最大上传大小为 5TB。
func (f *files) FPutObject(key, filePath string) (info minio.UploadInfo, err error) {
	info, err = f.oss.NewCustom().CheckOssFile(key)
	if err == nil && info.Size != 0 {
		return
	}
	opts := minio.PutObjectOptions{}
	info, err = f.oss.client.FPutObject(context.Background(), f.oss.bucketName, key, filePath, opts)
	return
}

// 获取文件信息
func (f *files) StatObject(key string) (stat minio.ObjectInfo, err error) {
	stat, err = f.oss.client.StatObject(context.Background(), f.oss.bucketName, key, minio.StatObjectOptions{})
	if err == nil && stat.Size > 0 {
		err = errors.New(fmt.Sprintf("file size error:%v", stat.Size))
	}
	return
}

// 删除文件
func (f *files) RemoveObject(key string) error {
	return f.oss.client.RemoveObject(context.Background(), f.oss.bucketName, key, minio.RemoveObjectOptions{})
}

// 批量删除文件
// 删除从输入通道获得的对象列表。该调用一次向服务器发送最多 1000 个对象的删除请求。观察到的错误通过错误通道发送。
func (f *files) RemoveObjects(keys []string) (res []error) {
	objectsCh := make(chan minio.ObjectInfo)
	defer close(objectsCh)
	for _, v := range keys {
		object, err := f.StatObject(v)
		if err == nil {
			objectsCh <- object
		} else {
			res = append(res, err)
		}
	}
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	for rErr := range f.oss.client.RemoveObjects(context.Background(), f.oss.bucketName, objectsCh, opts) {
		fmt.Println("Error detected during deletion: ", rErr)
	}
	return
}

// 上传终止
func (f *files) RemoveIncompleteUpload(key string) error {
	return f.oss.client.RemoveIncompleteUpload(context.Background(), f.oss.bucketName, key)
}

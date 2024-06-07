package oss

import (
	"fmt"
	"go-study/config"
	"log"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
)

type ossHandle struct {
	client     *minio.Client
	bucketName string
	path       string
}

func NewOssHandle(path ...string) *ossHandle {
	ossPath := config.Instance.S3.PrefixPath
	if len(path) > 0 && path[0] != "" {
		ossPath = path[0]
	}
	return &ossHandle{
		client:     minioClient,
		bucketName: config.Instance.S3.BucketName,
		path:       ossPath,
	}
}

// 获取需要上传的对象名称
func GetObjectPath(fileName string) string {
	return filepath.Join(config.Instance.S3.PrefixPath, fileName)
}

// 获取已经上传后的URL地址
func GetObjectUrlAddress(fileName string) string {
	s3 := &config.Instance.S3
	endpoint := s3.Endpoint
	if s3.External != "" {
		endpoint = s3.External
	}
	if s3.Ssl {
		return fmt.Sprintf("https://%s/%s/%s", endpoint, s3.BucketName, fileName)
	}
	return fmt.Sprintf("http://%s/%s/%s", endpoint, s3.BucketName, fileName)
}

// 初始化S3
func InitOss() (err error) {
	var creds *credentials.Credentials
	ak := config.Instance.S3.Ak
	sk := config.Instance.S3.Sk
	s3 := &ossHandle{}
	s3.bucketName = config.Instance.S3.BucketName
	if config.Instance.S3.Version == "V4" {
		creds = credentials.NewStaticV4(ak, sk, "")
	} else {
		creds = credentials.NewStaticV2(ak, sk, "")
	}
	// 初始化minio客户端
	s3.client, err = minio.New(config.Instance.S3.Endpoint, &minio.Options{
		Creds:  creds,
		Secure: config.Instance.S3.Ssl,
	})
	minioClient = s3.client

	if err != nil {
		log.Panic(err)
	}
	//桶检查
	exists, err := s3.NewBucket().BucketExists(s3.bucketName)
	if exists && err == nil {
		return
	}
	//桶不存在则创建
	err = s3.NewBucket().MakeBucket(s3.bucketName)
	if err != nil {
		//TODO::桶创建错误后续操作
		log.Panicf("桶创建错误,err:%v", err)
	}
	return
}

package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
)

type bucket struct {
	oss ossHandle
}

//桶操作
func (oss *ossHandle) NewBucket() *bucket {
	return &bucket{
		oss: *oss,
	}
}

//创建桶
func (b *bucket) MakeBucket(bucket string) error {
	return b.oss.client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
}

//桶列表
func (b *bucket) BucketList() ([]minio.BucketInfo, error) {
	return b.oss.client.ListBuckets(context.Background())
}

//检测桶是否存在
func (b *bucket) BucketExists(bucket string) (bool, error) {
	return b.oss.client.BucketExists(context.Background(), bucket)
}

//删除桶 桶内文件应为空
func (b *bucket) RemoveBucket(bucket string) error {
	return b.oss.client.RemoveBucket(context.Background(), bucket)
}

//列出桶内对象
func (b *bucket) BucketListObjects(bucket string) []minio.ObjectInfo {
	obj := b.oss.client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{})
	res := make([]minio.ObjectInfo, len(obj))
	i := 0
	for object := range obj {
		res[i] = object
	}
	return res
}

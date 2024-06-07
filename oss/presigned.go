package oss

import (
	"context"
	"net/url"
	"time"
)

type preSigned struct {
	oss ossHandle
}

//预操作
func (oss *ossHandle) NewPreSigned() *preSigned {
	return &preSigned{
		oss: *oss,
	}
}

//获取文件url
//expires地址有效期，单位为秒，最长为604800（即 7 天），最短为 1 秒。
func (p *preSigned) PreSignedGetObject(key string, expires int) (*url.URL, error) {
	if expires <= 0 {
		expires = 7 * 24 * 60 * 60
	}
	reqParams := make(url.Values)
	return p.oss.client.PresignedGetObject(context.Background(), p.oss.bucketName, key, time.Duration(expires)*time.Second, reqParams)
}

//为 HTTP PUT 操作生成预签名 URL。浏览器/移动客户端可能指向此 URL 以将对象直接上传到存储桶，即使它是私有的。
//此预签名 URL 可以具有关联的过期时间（以秒为单位），在此之后它不再可操作。默认有效期设置为 7 天。
//注意：您只能使用指定的对象名称上传到 S3。
func (p *preSigned) PreSignedPutObject(key string, expires int) (*url.URL, error) {
	if expires <= 0 {
		expires = 7 * 24 * 60 * 60
	}
	return p.oss.client.PresignedPutObject(context.Background(), p.oss.bucketName, key, time.Duration(expires)*time.Second)
}

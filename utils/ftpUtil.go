package utils

import (
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"time"
)

func DownloadFtpFile(host string, username, password string, path, fileName string) error {
	// 建立连接，默认用21端口
	c, err := ftp.Dial(host+":21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	defer c.Quit()

	// 登录
	err = c.Login(username, password)
	if err != nil {
		return err
	}

	// 切换到path
	err = c.ChangeDir(path)
	if err != nil {
		return err
	}

	// 读取文件
	body, err := c.Retr(fileName)
	if err != nil {
		return err
	}
	defer body.Close()

	// 创建本地文件
	localFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// 下载到本地
	_, err = io.Copy(localFile, body)
	return nil
}

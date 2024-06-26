package log

import (
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	FILE_SLICE_DATE_NULL  = ""
	FILE_SLICE_DATE_YEAR  = "y"
	FILE_SLICE_DATE_MONTH = "m"
	FILE_SLICE_DATE_DAY   = "d"
	FILE_SLICE_DATE_HOUR  = "h"
)

// file writer
type FileWriter struct {
	lock      sync.RWMutex
	writer    *os.File
	startTime int64
	filename  string
	DateSlice string // 通过日期来进行文件分割
	MaxSize   int64  // 通过文件大小来进行分割，单位kb
}

// 新建文件日志模块
// fn 文件名称
// dateSlice 通过日期来进行文件分割
// maxSize 通过文件大小来进行分割，单位kb
func NewFileWrite(fn, dateSlice string, maxSize int64) *FileWriter {
	fw := &FileWriter{
		filename:  fn,
		DateSlice: dateSlice,
		MaxSize:   maxSize,
	}
	fw.initFile()
	return fw
}

// init file
func (fw *FileWriter) initFile() error {

	//check file exits, otherwise create a file
	ok, _ := PathExists(fw.filename)
	if ok == false {
		err := CreateFile(fw.filename)
		if err != nil {
			return err
		}
	}

	// get start time
	fw.startTime = time.Now().Unix()

	// get a file pointer
	file, err := fw.getFileObject(fw.filename)
	if err != nil {
		return err
	}
	fw.writer = file
	return nil
}

// write by config
// 写入文件内容
func (fw *FileWriter) Write(b []byte) (n int, err error) {
	// 上层已经加了锁，此处可以不加载锁了
	// fw.lock.Lock()
	// defer fw.lock.Unlock()

	if fw.DateSlice != "" {
		// file slice by date
		err = fw.sliceByDate(fw.DateSlice)
		if err != nil {
			return
		}
	}

	if fw.MaxSize != 0 {
		// file slice by size
		err = fw.sliceByFileSize(fw.MaxSize)
		if err != nil {
			return
		}
	}

	n, err = fw.writer.Write(b)
	return
}

//slice file by date (y, m, d, h, i, s), rename file is file_time.log and recreate file
func (fw *FileWriter) sliceByDate(dataSlice string) error {

	filename := fw.filename
	filenameSuffix := path.Ext(filename)
	startTime := time.Unix(fw.startTime, 0)
	nowTime := time.Now()

	oldFilename := ""
	isHaveSlice := false
	if (dataSlice == FILE_SLICE_DATE_YEAR) &&
		(startTime.Year() != nowTime.Year()) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("2006") + filenameSuffix
	}
	if (dataSlice == FILE_SLICE_DATE_MONTH) &&
		(startTime.Format("200601") != nowTime.Format("200601")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("200601") + filenameSuffix
	}
	if (dataSlice == FILE_SLICE_DATE_DAY) &&
		(startTime.Format("20060102") != nowTime.Format("20060102")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("20060102") + filenameSuffix
	}
	if (dataSlice == FILE_SLICE_DATE_HOUR) &&
		(startTime.Format("2006010215") != startTime.Format("2006010215")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("2006010215") + filenameSuffix
	}

	if isHaveSlice == true {
		//close file handle
		fw.writer.Close()
		err := os.Rename(fw.filename, oldFilename)
		if err != nil {
			return err
		}
		err = fw.initFile()
		if err != nil {
			return err
		}
	}

	return nil
}

//slice file by size, if maxSize < fileSize, rename file is file_size_maxSize_time.log and recreate file
func (fw *FileWriter) sliceByFileSize(maxSize int64) error {

	filename := fw.filename
	filenameSuffix := path.Ext(filename)
	nowSize, _ := fw.getFileSize(filename)

	if nowSize >= maxSize {
		//close file handle
		fw.writer.Close()
		timeFlag := time.Now().Format("2006-01-02-15.04.05.9999")
		oldFilename := strings.Replace(filename, filenameSuffix, "", 1) + "." + timeFlag + filenameSuffix
		err := os.Rename(filename, oldFilename)
		if err != nil {
			return err
		}
		err = fw.initFile()
		if err != nil {
			return err
		}
	}

	return nil
}

//get file object
//params : filename
//return : *os.file, error
func (fw *FileWriter) getFileObject(filename string) (file *os.File, err error) {
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0766)
	return file, err
}

//get file size
//params : filename
//return : fileSize(byte int64), error
func (fw *FileWriter) getFileSize(filename string) (fileSize int64, err error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return fileSize, err
	}

	return fileInfo.Size() / 1024, nil
}

// create file
func CreateFile(filename string) error {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	return err
}

// file or path is exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

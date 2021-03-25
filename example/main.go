package main

import (
	"log"
	"path/filepath"

	"github.com/UangDesign/filetype"
)

var (
	GZ_PATH     = filepath.Join(getCurrentPath("./"), "file/gz/test.gz")
	TGZ_PATH    = filepath.Join(getCurrentPath("./"), "file/tgz/test.tgz")
	TAR_GZ_PATH = filepath.Join(getCurrentPath("./"), "file/tgz/test.tar.gz")
	ZIP_PATH    = filepath.Join(getCurrentPath("./"), "file/zip/test.zip")
	TAR_PATH    = filepath.Join(getCurrentPath("./"), "file/tar/test.tar")
	JSON_PATH   = filepath.Join(getCurrentPath("./"), "file/json/test.json")
)

func main() {
	fileTypeObj := filetype.NewFileType()
	// 获取文件类型
	if fileTypeObj.GetFileType(GZ_PATH) != filetype.FILE_TYPE_GZ {
		log.Fatal("fileype should be gz")
	}
	if fileTypeObj.GetFileType(TGZ_PATH) != filetype.FILE_TYPE_TGZ {
		log.Fatal("fileype should be tgz")
	}
	if fileTypeObj.GetFileType(ZIP_PATH) != filetype.FILE_TYPE_ZIP {
		log.Fatal("fileype should be zip")
	}
	if fileTypeObj.GetFileType(TAR_PATH) != filetype.FILE_TYPE_TAR {
		log.Fatal("fileype should be tar")
	}
	// 传入非二进制文件，进行类型判断
	if fileTypeObj.GetFileType(JSON_PATH) != filetype.FILE_TYPE_JSON {
		log.Fatal("fileype should be json")
	}
	// 判断文件类型
	if !fileTypeObj.CheckFileType(GZ_PATH, filetype.FILE_TYPE_GZ) {
		log.Fatal("file Type should be gz")
	}
	if !fileTypeObj.CheckFileType(TGZ_PATH, filetype.FILE_TYPE_TGZ) {
		log.Fatal("file Type should be tgz")
	}
	if !fileTypeObj.CheckFileType(TAR_GZ_PATH, filetype.FILE_TYPE_TAR_GZ) {
		log.Fatal("file Type should be tar.gz")
	}
	if !fileTypeObj.CheckFileType("D:/test/tgz/code.tgz", filetype.FILE_TYPE_GZ) {
		log.Fatal("file Type should be gz")
	}
	if !fileTypeObj.CheckFileType(ZIP_PATH, filetype.FILE_TYPE_ZIP) {
		log.Fatal("file Type should be zip")
	}
	if !fileTypeObj.CheckFileType(TAR_PATH, filetype.FILE_TYPE_TAR) {
		log.Fatal("file Type should be tar")
	}
}

func getCurrentPath(dir string) (absPath string) {
	if curPath, err := filepath.Abs(dir); err == nil {
		absPath = curPath
	}
	return absPath
}

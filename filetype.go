package filetype

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type FILE_TYPE = string

const (
	// 可执行文件类型 PE
	FILE_TYPE_EXE FILE_TYPE = ".exe"
	FILE_TYPE_COM FILE_TYPE = ".com"
	FILE_TYPE_DLL FILE_TYPE = ".dll"
	// 可执行文件类型 ELF
	FILE_TYPE_BIN FILE_TYPE = ".bin"
	FILE_TYPE_SO  FILE_TYPE = ".so"
	// 文本类型
	FILE_TYPE_HTML FILE_TYPE = ".html"
	FILE_TYPE_HTM  FILE_TYPE = ".htm"
	FILE_TYPE_JSON FILE_TYPE = ".json"
	// 压缩类型
	FILE_TYPE_ZIP    FILE_TYPE = ".zip"
	FILE_TYPE_TAR    FILE_TYPE = ".tar"
	FILE_TYPE_GZ     FILE_TYPE = ".gz"
	FILE_TYPE_GZIP   FILE_TYPE = ".gz"
	FILE_TYPE_TGZ    FILE_TYPE = ".tgz"
	FILE_TYPE_TAR_GZ FILE_TYPE = ".tar.gz"

	// 视频类型
	FILE_TYPE_AVI  FILE_TYPE = ".avi"
	FILE_TYPE_MPG  FILE_TYPE = ".mpg"
	FILE_TYPE_RM   FILE_TYPE = ".rm"
	FILE_TYPE_WMV  FILE_TYPE = ".wmv"
	FILE_TYPE_FLV  FILE_TYPE = ".flv"
	FILE_TYPE_MP4  FILE_TYPE = ".mp4"
	FILE_TYPE_RMVB FILE_TYPE = ".rmvb"
	// 音频类型
	FILE_TYPE_MP3 FILE_TYPE = ".mp3"
	FILE_TYPE_WAV FILE_TYPE = ".wav"
	FILE_TYPE_WMA FILE_TYPE = ".wma"
	// 图片类型
	FILE_TYPE_JPG  FILE_TYPE = ".jpg"
	FILE_TYPE_BMP  FILE_TYPE = ".bmp"
	FILE_TYPE_GIF  FILE_TYPE = ".gif"
	FILE_TYPE_TIF  FILE_TYPE = ".tif"
	FILE_TYPE_PNG  FILE_TYPE = ".png"
	FILE_TYPE_JPEG FILE_TYPE = ".jpeg"
)

const (
	// TAR 打包文件最小 1536 个字节
	TAR_MIN_SIZE = 1536
	// TAR 文件标志偏移量起始位置
	TAR_OFFSET_POSITION = 257
	// TAR 文件标志偏移量
	TAR_OFFSET_SIZE = 5
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("504b0304", ".zip")   // zip 文件
	fileTypeMap.Store("1f8b", ".gz")        // gz 文件, tgz, tar.gz(tgz 和 tar.gz 为一个类型，需要单独处理)
	fileTypeMap.Store("7573746172", ".tar") // tar 文件，offset 以 275 为偏移，往后算5位
}

type IFileType interface {
	GetFileType(string) string
	CheckFileType(string, FILE_TYPE) bool
}

type FileType struct {
}

func NewFileType() IFileType {
	return &FileType{}
}

func (f *FileType) GetFileType(srcFile string) (fileType string) {
	buf := make([]byte, TAR_MIN_SIZE)
	file, err := os.Open(srcFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("open srcFaile failed, err is:%v", err)
	} else {
		n, err := file.Read(buf)
		if err != nil {
			log.Fatalf("read file failed, err is:%v", err)
		} else {
			fileType = f.getFileType(buf[:n])
			// 需要进一步判断类型是否为 tar.gz(tgz) 类型
			if fileType == FILE_TYPE_GZ {
				file.Seek(0, 0)
				tgzReader, err := gzip.NewReader(file)
				if err == nil {
					_, err := tar.NewReader(tgzReader).Next()
					if err == nil {
						fileType = FILE_TYPE_TGZ
					}
				}
			}
		}
	}
	// 如果如果正常识别文件类型， fileType 为"", 则需要通过其他方法进行判断
	if fileType == "" {
		fileType = f.getFileTypeByText(srcFile)
	}
	return fileType
}

func (f *FileType) getFileTypeByText(srcFile string) (fileType string) {
	file, err := os.Open(srcFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("open srcFaile failed, err is:%v", err)
	} else {
		if fileReader, err := ioutil.ReadAll(file); err != nil {
			log.Fatalf("read file from srcFaile failed, err is:%v", err)
		} else {
			// 判断是否为 json 格式
			if isJsonType(fileReader) {
				fileType = FILE_TYPE_JSON
			}
		}
	}
	return fileType
}

func isJsonType(buf []byte) (isJson bool) {
	return json.Unmarshal(buf, &map[string]interface{}{}) == nil
}

func (f *FileType) getFileType(buf []byte) (fileType string) {
	fileCode := byteToHexString(buf)
	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, k) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

func (f *FileType) CheckFileType(srcFile string, fileType FILE_TYPE) (is bool) {
	fType := f.GetFileType(srcFile)
	if fType == fileType {
		is = true
	} else if fType == FILE_TYPE_TGZ {
		if fileType == FILE_TYPE_TAR_GZ || fileType == FILE_TYPE_GZ {
			is = true
		}
	} else {
		is = false
	}
	return is
}

func byteToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) < 4 {
		res.Write([]byte{})
	} else {
		for _, v := range src[:8] {
			hv := hex.EncodeToString([]byte{v & 0xFF})
			if len(hv) < 2 {
				res.WriteString(strconv.FormatInt(int64(0), 10))
			}
			res.WriteString(hv)
		}
		// 判断是不是 tar 类型
		if len(src) >= TAR_MIN_SIZE {
			tarRes := bytes.Buffer{}
			for _, v := range src[TAR_OFFSET_POSITION : TAR_OFFSET_POSITION+6] {
				hv := hex.EncodeToString([]byte{v & 0xFF})
				if len(hv) < 2 {
					res.WriteString(strconv.FormatInt(int64(0), 10))
				}
				tarRes.WriteString(hv)
			}
			if strings.HasPrefix(tarRes.String(), "7573746172") {
				res = tarRes
			}
		}
	}
	return res.String()
}

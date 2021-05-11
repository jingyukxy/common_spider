package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/huichen/sego"
	"github.com/mozillazg/go-pinyin"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// 是否有中文字符
func IsChineseCharacter(input string) bool {
	for _, r := range input {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

// md5 []byte
func Md5(input []byte) string {
	digest := md5.New()
	digest.Write(input)
	return hex.EncodeToString(digest.Sum(nil))
}

// 三元表达
func If(condition bool, trueVal interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseValue
}

// md5 string
func Md5WithString(input string) string {
	return Md5([]byte(input))
}

// 获取拼音首字母
func GetFirstLetterOfPinYin(name string) string {
	if len(name) == 0 {
		return name
	}
	if !IsChineseCharacter(name) {
		return name[0:1]
	}
	pyArgs := pinyin.NewArgs()
	result := pinyin.Pinyin(name, pyArgs)
	if len(result[0][0]) > 0 {
		return string(result[0][0][0])
	}
	return ""
}

// 获取app 运行程序目录
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, "/bin")
	if index != -1 {
		return path[:index]
	}
	return path
}

// 路径是否存在
func PathExists(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SegmentTxt() {
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("../github.com/huichen/sego/data/dictionary.txt")
	text := []byte("使用它可以进行快速开发，同时它还是一个真正的编译语言，我们之所以现在将其开源，原因是我们认为它已经非常有用和强大")
	segments := segmenter.Segment(text)
	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToString(segments, false))
}

// 下载文件
func DownloadFile(filePath string, suffix string, url string) (name string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	uid := uuid.NewV4()
	name = strings.ReplaceAll(uid.String(), "-", "")
	realPath := fmt.Sprintf("%s/%s.%s", filePath, name, suffix)
	if err := ioutil.WriteFile(realPath, data, 0644); err != nil {
		return "", err
	}
	return
}

package library

import (
	"github.com/resize"
	"io/ioutil"
	"net/http"
	"os"
	//"strconv"
	//"time"
	"image/jpeg"
	//"log"
	//"fmt"
)

//文件夹是否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

//创建文件夹
func Mkdir(path string) {
	if IsDirExists(path) == false {
		os.Mkdir(path, os.ModePerm)
	}
}

//文件是否存在
func isFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//下载文件到picture目录
//imgurl:网络图片地址
func GetImage(imgUrl string) (string, bool) {
	Mkdir("picture")
	response, err := http.Get(imgUrl)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	file := "picture/tmp.jpg"
	image, err := os.Create(file)
	if err != nil {
		return "", false
	}
	image.Write(data)
	image.Close()
	resizeSucess := ImageResizeToSmallThen(file, 150*1000) //压缩到150K以下
	if resizeSucess == false {
		return "", false
	}
	return "picture/tmp.jpg", true
}

//下载文件到指定文件
//imgurl:网络图片地址
func GetImageNew(imgUrl string, file string) (string, bool) {
	Mkdir("picture")
	file = "picture/" + file
	response, err := http.Get(imgUrl)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	image, err := os.Create(file)
	if err != nil {
		return "", false
	}
	image.Write(data)
	image.Close()
	return file, true
}

//压缩到指定大小以下
func ImageResizeToSmallThen(filename string, size int64) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	info, err := file.Stat()
	file.Close()
	if err != nil {
		return false
	}
	var sizeX uint = 1000
	var sizeY uint = 1000

	for info.Size() > size {
		resizeSucess := ImageResize(filename, sizeX, sizeY)
		if resizeSucess == false {
			return false
		}
		file, _ = os.Open(filename)
		info, _ = file.Stat()
		file.Close()
		sizeX = sizeX - sizeX/5
		sizeY = sizeY - sizeY/5
	}
	return true
}

//图片压缩
func ImageResize(filename string, sizeX uint, sizeY uint) bool {
	in, err := os.Open(filename)
	if err != nil {
		return false
	}
	img, err := jpeg.Decode(in)
	in.Close()
	if err != nil {
		return false
	}
	m := resize.Thumbnail(sizeX, sizeY, img, resize.Lanczos3)
	out, err := os.Create(filename)
	if err != nil {
		return false
	}
	jpeg.Encode(out, m, nil)
	out.Close()
	return true
}

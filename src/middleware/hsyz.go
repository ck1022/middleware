package main

import (
	"encoding/json"
	"io/ioutil"
	"library"
	//"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	type stuinfo struct {
		Mid          int
		PersonId     int
		ImgFace      string
		IDCard       string
		Gender       string
		Name         string
		EducationNum string
	}

	type allStu struct {
		IsSuccess bool
		Retval    []stuinfo
	}

	library.Mkdir("log")
	library.Mkdir("picture")
	//获取照片信息

	b, e := Send("http://192.168.37.104:8080/api/member/getUserinfoAndpic", "")
	//b, e := ioutil.ReadFile("getUserinfoAndpic")
	//log.Println(string(b))
	if e {
		var info = new(allStu)
		json.Unmarshal(b, &info)
		code := info.IsSuccess
		if code {
			i := 0
			for _, v := range info.Retval {
				library.Println("log/l", "+", i, "+", v.Name, v.ImgFace)
				time.Sleep(100 * time.Millisecond)
				GetImageNew("http://192.168.37.104:8080"+v.ImgFace, v.Name+".jpg")
				i++
			}
		}
	}
}

/**
*提交请求
*url：请求地址
*param：请求参数
**/
func Send(url string, param string) ([]byte, bool) {

	req, err := http.NewRequest("GET", url, strings.NewReader(param))
	if err != nil {
		return nil, false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, true
}

//下载文件到指定文件
//imgurl:网络图片地址
func GetImageNew(imgUrl string, file string) (string, bool) {
	library.Mkdir("picture")
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

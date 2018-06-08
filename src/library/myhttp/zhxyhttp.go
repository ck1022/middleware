package myhttp

import (
	//"bytes"
	"config"
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"
	"library"
	"mime/multipart"
	"net/http"
	//"os"
	"bytes"
	"io"
	"model"
	"os"
	"strconv"
	"strings"
	"time"
)

type Zhxyhttp struct {
	Myconfig      *config.MyConfig
	ContentType   string
	RequestMethod string
	Token         string
	TokenExpire   int64
	LogFile       string
}

//获取token返回
type TokenReturn struct {
	Code    int
	Message string
	Data    struct {
		Token  string
		Expire string
	}
}

/***
*获取token
*return token,expire,code
 */
func (zhxy *Zhxyhttp) GetToken() bool {
	timestamp := time.Now().Unix()
	//提前两小时更新token
	if zhxy.TokenExpire-timestamp > 7200 {
		return true
	}
	zhxyConfig := zhxy.Myconfig.Zhxy
	param := fmt.Sprintf("appkey=%s&timestamp=%d&sign=%s", zhxyConfig.AppKey, timestamp, library.Mymd5(fmt.Sprintf("%s%d%s", zhxyConfig.AppKey, timestamp, zhxyConfig.AppSecret)))
	//log.Println(param)
	newurl := fmt.Sprintf("%s?%s", zhxyConfig.TokenUrl, param)
	resp, err := http.Get(newurl)
	defer resp.Body.Close()
	if err != nil {
		library.Panicln(zhxy.LogFile, "获取token失败", err.Error())
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		library.Panicln(zhxy.LogFile, "读取token失败", err.Error())
		return false
	}
	var tokenReturn = new(TokenReturn)
	json.Unmarshal(body, &tokenReturn)
	if tokenReturn.Code == 1 {
		zhxy.Token = tokenReturn.Data.Token
		expire := tokenReturn.Data.Expire
		tokenExpireTmp, _ := time.Parse("2006-01-02 15:04:05", expire)
		zhxy.TokenExpire = tokenExpireTmp.Unix()
		return true
	} else {
		library.Panicln(zhxy.LogFile, "获取token返回错误", err.Error())
		return false
	}
}

/***
**进程活动反馈
 */
func (zhxy *Zhxyhttp) TaskAlive(taskname string) {
	zhxyConfig := zhxy.Myconfig.Zhxy
	param := fmt.Sprintf("xxid=%s&xqid=%s&taskname=%s", zhxyConfig.SchoolId, zhxyConfig.CampusId, taskname)
	newurl := fmt.Sprintf("%s?%s", zhxyConfig.Domain+"/Api/MiddlewareMonitor/update", param)
	resp, _ := http.Get(newurl)
	resp.Body.Close()
}

/**
*获取未读消息队列
**/
func (zhxy *Zhxyhttp) GetMessageQueueList() ([]model.Message, bool) {
	var paramMap = make(map[string]string)
	zhxyConfig := zhxy.Myconfig.Zhxy
	messageQueueConfig := zhxy.Myconfig.MessageQueue
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	paramMap["topic"] = messageQueueConfig.Topic
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(messageQueueConfig.GetMessageUrl, param)
	if e == false {
		return nil, false
	} else {
		var info = new(model.MessageReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/**
*获取学生列表
*option：查询参数map，xsid：学生id
**/
func (zhxy *Zhxyhttp) GetStudentList(option map[string]string) ([]model.Student, bool) {
	var paramMap = option
	zhxyConfig := zhxy.Myconfig.Zhxy
	hk96Config := zhxy.Myconfig.Hk96
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(hk96Config.GetStudentUrl, param)
	if e == false {
		return nil, false
	} else {
		var info = new(model.StudentReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			library.Panicln(zhxy.LogFile, "获取学生列表失败：", string(b))
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/**
*获取教职工列表
*option：查询参数map，jzgid：教职工id，详见接口文档
**/
func (zhxy *Zhxyhttp) GetTeacherList(option map[string]string) ([]model.Teacher, bool) {
	var paramMap = option
	zhxyConfig := zhxy.Myconfig.Zhxy
	hk96Config := zhxy.Myconfig.Hk96
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(hk96Config.GetTeacherUrl, param)
	if e == false {
		return nil, false
	} else {
		var info = new(model.TeacherReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/**
*获取通校和走读学生列表
*date：日期
*in：进或出，true：进，false：出
**/
func (zhxy *Zhxyhttp) GetLeaveInStudentList(date string, in bool) ([]model.Student, bool) {
	var paramMap = make(map[string]string)
	zhxyConfig := zhxy.Myconfig.Zhxy
	hk96Config := zhxy.Myconfig.Hk96
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	paramMap["date"] = date
	paramMap["type"] = "2"
	if in {
		paramMap["type"] = "1"
	}
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(hk96Config.GetThoroughStudentUrl, param)
	if e == false {
		return nil, false
	} else {
		var info = new(model.StudentReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/***
*获取请假学生列表
*
**/
func (zhxy *Zhxyhttp) GetLeaveStudentList(date string) ([]model.Student, bool) {
	var paramMap = make(map[string]string)
	zhxyConfig := zhxy.Myconfig.Zhxy
	hk96Config := zhxy.Myconfig.Hk96
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	paramMap["kssj"] = date
	paramMap["jssj"] = date + " 23:59:59"
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(hk96Config.GetLeaveStudenturl, param)
	if e == false {
		return nil, false
	} else {
		var info = new(model.StudentReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/**
*获取空照片的考勤记录
*
**/
func (zhxy *Zhxyhttp) GetNullAttendanceFaceList() ([]model.Face, bool) {
	var paramMap = make(map[string]string)
	zhxyConfig := zhxy.Myconfig.Zhxy
	hk96Config := zhxy.Myconfig.Hk96
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	//paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	//paramMap["token"] = zhxy.Token
	//paramMap["date"] = date
	//paramMap["type"] = "2"
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(hk96Config.GetAttendanceFaceList, param)
	if e == false {
		return nil, false
	} else {
		var info = new(model.NullAttendanceFaceReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/**
*获取学生列表(cip)
*option：查询参数map，xsid：学生id
**/
func (zhxy *Zhxyhttp) GetStudentListCip(updateDate string) ([]model.StudentCip, bool) {
	type studentReturn struct {
		Code int
		Data struct {
			List []model.StudentCip
		}
	}
	var paramMap = make(map[string]string)
	zhxyConfig := zhxy.Myconfig.Zhxy
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	paramMap["gxsj"] = updateDate
	paramMap["sfsc"] = "2"
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(zhxyConfig.Domain+"/Api/Student/lists", param)
	if e == false {
		return nil, false
	} else {
		var info = new(studentReturn)
		json.Unmarshal(b, &info)
		code := info.Code
		if code == -1 {
			library.Panicln(zhxy.LogFile, "获取学生列表失败：", string(b))
			return nil, false
		}
		list := info.Data.List
		return list, true
	}
}

/**
*上传教室考勤记录（cip）
*
**/
func (zhxy *Zhxyhttp) UploadClassRoomAttendance(paramMap map[string]string) ([]byte, bool) {
	zhxyConfig := zhxy.Myconfig.Zhxy
	hkCipConfig := zhxy.Myconfig.HkCip
	paramMap["appkey"] = zhxyConfig.AppKey
	paramMap["xxid"] = zhxyConfig.SchoolId
	paramMap["xqid"] = zhxyConfig.CampusId
	paramMap["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	paramMap["token"] = zhxy.Token
	var paramArray []string
	for k, v := range paramMap {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}
	param := strings.Join(paramArray, "&")
	b, e := zhxy.Send(hkCipConfig.UploadStudentClassRoomAttendanceDataUrl, param)
	if e == false {
		return nil, false
	} else {
		return b, true
	}
}

//上传比对照片到美智(cip)
func (zhxy *Zhxyhttp) UpPic(filename string, ysjlid string) (string, bool) {
	type picUpReturn struct {
		Code int
		Data string
	}
	filename = "picture/" + filename
	zhxyConfig := zhxy.Myconfig.Zhxy
	hk96Config := zhxy.Myconfig.Hk96
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	body_writer.WriteField("appkey", zhxyConfig.AppKey)
	body_writer.WriteField("xxid", zhxyConfig.SchoolId)
	body_writer.WriteField("xqid", zhxyConfig.CampusId)
	body_writer.WriteField("ysjlid", ysjlid)
	boundary := body_writer.Boundary()
	url := hk96Config.UploadFaceImageUrl
	library.Println(zhxy.LogFile, url)
	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		library.Panicln(zhxy.LogFile, "图片写入缓存出错", err.Error())
		return "", false
	}
	fh, err := os.Open(filename)
	if err != nil {
		library.Panicln(zhxy.LogFile, "打开图片失败", err.Error())
		return "", false
	}

	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		library.Panicln(zhxy.LogFile, "读取图片失败", err.Error())
		return "", false
	}
	req, err := http.NewRequest("POST", url, request_reader)
	if err != nil {
		library.Panicln(zhxy.LogFile, "发送失败pic1", err.Error())
		return "", false
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Panicln(zhxy.LogFile, "发送失败pic1", err.Error())
		return "发送失败pic2：" + err.Error(), false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		library.Panicln(zhxy.LogFile, "解析返回数据出错pic", err.Error())
		return "", false
	}
	var returnObject = new(picUpReturn)
	json.Unmarshal(body, &returnObject)
	if returnObject.Code != 1 {
		library.Println(zhxy.LogFile, "上传图片失败", err.Error())
		return "", false
	}
	return returnObject.Data, true
}

//上传行为数据到美智
func (zhxy *Zhxyhttp) UploadWithFile(url string, filename string, paramMap map[string]string) (string, bool) {
	type uploadReturn struct {
		Code int
		Data string
	}

	filename = "picture/" + filename
	zhxyConfig := zhxy.Myconfig.Zhxy
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	body_writer.WriteField("appkey", zhxyConfig.AppKey)
	body_writer.WriteField("xxid", zhxyConfig.SchoolId)
	body_writer.WriteField("xqid", zhxyConfig.CampusId)
	for k, v := range paramMap {
		body_writer.WriteField(k, v)
	}
	boundary := body_writer.Boundary()

	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		library.Println(zhxy.LogFile, "图片写入缓存出错", err.Error())
		return "", false
	}
	fh, err := os.Open(filename)
	if err != nil {
		library.Println(zhxy.LogFile, "打开图片失败", err.Error())
		return "", false
	}

	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		library.Println(zhxy.LogFile, "读取图片失败", err.Error())
		return "", false
	}
	req, err := http.NewRequest("POST", url, request_reader)
	if err != nil {
		library.Println(zhxy.LogFile, "上传数据失败", err.Error())
		return "", false
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Println(zhxy.LogFile, "上传数据失败", err.Error())
		return "", false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		library.Println(zhxy.LogFile, "解析返回数据出错", err.Error())
		return "", false
	}
	var returnObject = new(uploadReturn)
	json.Unmarshal(body, &returnObject)
	library.Println(zhxy.LogFile, string(body))
	if returnObject.Code != 1 {
		library.Println(zhxy.LogFile, "上传数据失败")
		return "", false
	}
	return returnObject.Data, true
}

/**
*提交请求
*url：请求地址
*param：请求参数
**/
func (zhxy *Zhxyhttp) Send(url string, param string) ([]byte, bool) {
	library.Println(zhxy.LogFile, url)
	library.Println(zhxy.LogFile, string(param))
	if zhxy.ContentType == "" {
		zhxy.ContentType = "application/x-www-form-urlencoded"
	}
	if zhxy.RequestMethod == "" {
		zhxy.RequestMethod = "POST"
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		library.Panicln(zhxy.LogFile, "发起请求失败", err.Error())
		return nil, false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Panicln(zhxy.LogFile, "发送请求失败", err.Error())
		return nil, false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		library.Println(zhxy.LogFile, string(body), err.Error())
		return nil, false
	}
	library.Println(zhxy.LogFile, string(body))
	return body, true
}

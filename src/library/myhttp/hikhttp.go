package myhttp

import (
	"bytes"
	"config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"library"
	//"log"
	"mime/multipart"
	"net/http"
	//"net/url"
	"os"
	"strings"
	"time"
)

type Hikhttp struct {
	Myconfig      *config.MyConfig
	ContentType   string
	RequestMethod string
	LogFile       string
}

//添加或修改用户后更新目标库
func (t *Hikhttp) UpdateUserUpdateTarget(learnCode string) {
	t.DeleteTarget(t.Myconfig.Hk96.TestDormDevice, learnCode)
	t.AddTarget(t.Myconfig.Hk96.TestDormDevice, learnCode) //基本目标库
	if t.Myconfig.Hk96.LeaveInOutOpen != "on" {            //无进出校业务
		return
	}

	if t.Myconfig.Hk96.InOpen != "on" { //进校不控制，所有人都可以进
		t.DeleteTarget(t.Myconfig.Hk96.InDevice, learnCode)
		t.AddTarget(t.Myconfig.Hk96.InDevice, learnCode)
	}
	if t.Myconfig.Hk96.OutOpen != "on" { //出校不控制，所有人都可以进
		t.DeleteTarget(t.Myconfig.Hk96.OutDevice, learnCode)
		t.AddTarget(t.Myconfig.Hk96.OutDevice, learnCode)
	}
}

//添加或修改用户后更新目标库
func (t *Hikhttp) UpdateTeacherUpdateTarget(teacherCode string) {
	t.DeleteTarget(t.Myconfig.Hk96.TestDormDevice, teacherCode)
	t.AddTarget(t.Myconfig.Hk96.TestDormDevice, teacherCode) //基本目标库
	if t.Myconfig.Hk96.LeaveInOutOpen != "on" {              //无进出校业务
		return
	}
	t.DeleteTarget(t.Myconfig.Hk96.InDevice, teacherCode)
	t.AddTarget(t.Myconfig.Hk96.InDevice, teacherCode)
	t.DeleteTarget(t.Myconfig.Hk96.OutDevice, teacherCode)
	t.AddTarget(t.Myconfig.Hk96.OutDevice, teacherCode)

}

/**
*添加目标库
*deveiceCode：设备通道编号
*learncode：学生编号，多个用半角逗号分隔
**/
func (t *Hikhttp) AddTarget(deviceCode string, ids string) (string, bool) {
	map1 := make(map[string]string)
	map1["indexCode"] = ids
	map1["cameraIndex"] = deviceCode
	b, _ := json.Marshal(map1)
	r, e := t.Send(t.Myconfig.Hk96.TargetAddUrl, b)
	if e == false {
		library.Panicln(t.LogFile, "["+deviceCode+"]["+ids+"]添加到目标库失败：", r)
	} else {
		library.Println(t.LogFile, "["+deviceCode+"]["+ids+"]添加到目标库完成：", r)
	}
	return r, e
}

/**
*删除目标库
*deveiceCode：设备通道编号
*learncode：学生编号，多个用半角逗号分隔
**/
func (t *Hikhttp) DeleteTarget(deviceCode string, ids string) (string, bool) {
	map1 := make(map[string]string)
	map1["indexCode"] = ids
	map1["cameraIndex"] = deviceCode
	b, _ := json.Marshal(map1)
	r, e := t.Send(t.Myconfig.Hk96.TargetDeleteUrl, b)
	if e == false {
		library.Panicln(t.LogFile, "["+deviceCode+"]["+ids+"]删除目标库失败：", r)
	} else {
		library.Println(t.LogFile, "["+deviceCode+"]["+ids+"]删除目标库完成：", r)
	}
	return r, e
}

/**
*添加用户
*id：用户id，学生为学号
*name：姓名
*sex：性别
*pictureurl：照片地址，海康那边上传成功后的地址
*idno
*phone
**/
func (t *Hikhttp) AddUser(id string, name string, sex string, pictureurl string, idno string, phone string) (string, bool) {
	map1 := make(map[string]string)
	map1["userIndex"] = id
	map1["userName"] = name
	map1["sex"] = sex
	map1["picturePath"] = pictureurl
	map1["idNo"] = idno
	map1["mobilephone"] = phone
	library.Println(t.LogFile, map1)
	b, _ := json.Marshal(map1)
	return t.Send(t.Myconfig.Hk96.UserAddUrl, b)
}

/**
*修改用户
*id：用户id，学生为学号
*name：姓名
*sex：性别
*pictureurl：照片地址，海康那边上传成功后的地址
*idno
*phone
**/
func (t *Hikhttp) UpdateUser(id string, name string, sex string, pictureurl string, idno string, phone string) (string, bool) {
	map1 := make(map[string]string)
	map1["userIndex"] = id
	map1["userName"] = name
	map1["sex"] = sex
	map1["picturePath"] = pictureurl
	map1["idNo"] = idno
	map1["mobilephone"] = phone
	b, _ := json.Marshal(map1)
	return t.Send(t.Myconfig.Hk96.UserUpdateUrl, b)
}

/**
*删除用户
*ids：用户id，多个用半角逗号分隔
**/
func (t *Hikhttp) DeleteUser(ids []string) (string, bool) {
	var map1 = make(map[string][]string)
	map1["personNoList"] = ids
	b, _ := json.Marshal(map1)
	return t.Send(t.Myconfig.Hk96.UserDeleteUrl, b)
}

/**
**获取人脸行为抓拍数据
**/
func (t *Hikhttp) GetStudentBehaveData(startTime string, endTime string, page int, pageSize int) ([]byte, bool) {
	map1 := make(map[string]string)
	map1["page"] = fmt.Sprintf("%d", page)
	map1["size"] = fmt.Sprintf("%d", pageSize)
	map1["startTime"] = startTime
	map1["endTime"] = endTime
	return t.Get(t.Myconfig.HkCip.GetStudentBehaveDataUrl, map1)
}

/**
**获取人脸表情抓拍数据
**/
func (t *Hikhttp) GetStudentExpressionFaceData(startTime string, endTime string, page int, pageSize int) ([]byte, bool) {
	map1 := make(map[string]string)
	map1["page"] = fmt.Sprintf("%d", page)
	map1["size"] = fmt.Sprintf("%d", pageSize)
	map1["startTime"] = startTime
	map1["endTime"] = endTime
	map1["hasExpression"] = "true"
	return t.Get(t.Myconfig.HkCip.GetStudentExpressionFaceDataUrl, map1)
}

/**
**获取人脸识别抓拍数据
**/
func (t *Hikhttp) GetStudentClassRoomAttendanceData(startTime string, endTime string, page int, pageSize int) ([]byte, bool) {
	map1 := make(map[string]string)
	map1["page"] = fmt.Sprintf("%d", page)
	map1["size"] = fmt.Sprintf("%d", pageSize)
	map1["startDate"] = startTime
	map1["endDate"] = endTime
	return t.Get(t.Myconfig.HkCip.GetStudentClassRoomAttendanceUrl, map1)
}

/**
**添加或修改年级（cip）
**/

func (t *Hikhttp) AddGradeCip(name string, schoolOrgCode string) (string, bool) {
	type param struct {
		GradeName      string   `json:"gradeName"`
		ParentOrgCodes []string `json:"parentOrgCodes"`
	}
	param1 := param{
		GradeName:      name,
		ParentOrgCodes: []string{schoolOrgCode},
	}
	b, _ := json.Marshal(param1)
	return t.Send(t.Myconfig.HkCip.GradeAddUrl, b)
}

/**
**添加或修改班级（cip）
**/

func (t *Hikhttp) AddClassCip(name string, gradeOrgCode string) (string, bool) {
	map1 := make(map[string]string)
	map1["clazzName"] = name
	map1["gradeOrgCode"] = gradeOrgCode
	b, _ := json.Marshal(map1)
	return t.Send(t.Myconfig.HkCip.ClassAddUrl, b)
}

/**
*获取组织机构列表（cip）
 **/
func (t *Hikhttp) GetOrgTree() ([]byte, bool) {
	map1 := make(map[string]string)
	map1["orgCode"] = "000000"
	return t.Get(t.Myconfig.HkCip.OrgTreeUrl, map1)
}

/**
*添加或修改学生(cip)
**/
func (t *Hikhttp) AddOrUpdateStudentCip(personNo string, name string, sex string, orgCode string, pictureurl string) (string, bool) {
	type param struct {
		PersonNo    string   `json:"personNo"`
		Name        string   `json:"name"`
		Gender      string   `json:"gender"`
		OrgCode     string   `json:"orgCode"`
		PicturePath []string `json:"picturePath"`
	}
	param1 := param{
		PersonNo:    personNo,
		Name:        name,
		Gender:      sex,
		OrgCode:     orgCode,
		PicturePath: []string{pictureurl},
	}
	b, _ := json.Marshal(param1)
	return t.Send(t.Myconfig.HkCip.StudentAddOrUpdateUrl, b)
}

/**
*删除学生(cip)
 */
func (t *Hikhttp) DeleteStudentCip(ids []string) (string, bool) {
	map1 := make(map[string][]string)
	map1["personNoList"] = ids
	b, _ := json.Marshal(map1)
	return t.Send(t.Myconfig.HkCip.StudentDeleteUrl, b)
}

//上传头像到海康（96）
func (t *Hikhttp) UpPic(filename string) (string, bool) {
	type picUpReturn struct {
		Code int
		Data string
	}
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	url := t.Myconfig.Hk96.UserUploadPictrueUrl
	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		library.Println(t.LogFile, "图片写入缓存出错")
		return "", false
	}
	fh, err := os.Open(filename)
	if err != nil {
		library.Println(t.LogFile, "打开图片失败")
		return "", false
	}
	boundary := body_writer.Boundary()
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		library.Println(t.LogFile, "读取图片失败")
		return "", false
	}

	req, err := http.NewRequest("POST", url, request_reader)
	if err != nil {
		library.Println(t.LogFile, "发送给海康失败pic1："+err.Error())
		return "", false
	}
	timestemp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	wjwAuthorization := library.Mysha256(t.Myconfig.Hk.AppKey + t.Myconfig.Hk.AppCode + timestemp)
	req.Header.Set("Appkey", t.Myconfig.Hk.AppKey)
	req.Header.Set("AppSecret", timestemp)
	req.Header.Set("wjwAuthorization", string(wjwAuthorization))
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Println(t.LogFile, "发送给海康失败pic2："+err.Error())
		return "", false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	library.Println(t.LogFile, string(body))
	if err != nil {
		library.Panicln(t.LogFile, err.Error())
		return "", false
	}
	var returnObject = new(picUpReturn)
	json.Unmarshal(body, &returnObject)
	if returnObject.Code != 20000 {
		return "", false
	}
	return returnObject.Data, true
}

//上传文件到海康(cip))
func (t *Hikhttp) UpFile(filename string) (string, bool) {
	type picUpReturn struct {
		Status string
		Msg    string
		Data   string
	}
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	body_writer.WriteField("fileType", "图片")
	body_writer.WriteField("domainId", "11111")
	boundary := body_writer.Boundary()
	url := t.Myconfig.HkCip.UploadFileUrl
	// use the body_writer to write the Part headers to the buffer

	_, err := body_writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		library.Panicln(t.LogFile, "图片写入缓存出错：", err.Error())
		return "", false
	}

	fh, err := os.Open(filename)
	if err != nil {
		library.Panicln(t.LogFile, "读打开图片失败：", err.Error())
		return "", false
	}

	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		library.Panicln(t.LogFile, "读取图片出错：", err.Error())
		return "", false
	}

	req, err := http.NewRequest("POST", url, request_reader)
	if err != nil {
		library.Panicln(t.LogFile, "上传图片发起请求失败：", err.Error())
		return "", false
	}
	timestemp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	wjwAuthorization := library.Mysha256(t.Myconfig.Hk.AppKey + t.Myconfig.Hk.AppCode + timestemp)
	req.Header.Set("Appkey", t.Myconfig.Hk.AppKey)
	req.Header.Set("AppSecret", timestemp)
	req.Header.Set("wjwAuthorization", string(wjwAuthorization))
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Panicln(t.LogFile, "上传图片发送请求失败：", err.Error())
		return "", false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	library.Println(t.LogFile, string(body))
	if err != nil {
		library.Panicln(t.LogFile, "上传图片返回数据错误：", err.Error())
		return "", false
	}
	var returnObject = new(picUpReturn)
	json.Unmarshal(body, &returnObject)
	if returnObject.Status != "success" {
		library.Panicln(t.LogFile, "上传图片失败")
		return "", false
	}
	return returnObject.Data, true
}

//把数据提交给海康(96)
func (t *Hikhttp) Send(url string, param []byte) (string, bool) {
	library.Println(t.LogFile, url)
	library.Println(t.LogFile, string(param))
	if t.ContentType == "" {
		t.ContentType = "application/json"
	}
	if t.RequestMethod == "" {
		t.RequestMethod = "POST"
	}
	req, err := http.NewRequest(t.RequestMethod, url, strings.NewReader(string(param)))
	if err != nil {
		library.Panicln(t.LogFile, "发起请求出错", err.Error())
		return "", false
	}
	timestemp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	wjwAuthorization := library.Mysha256(t.Myconfig.Hk.AppKey + t.Myconfig.Hk.AppCode + timestemp)
	req.Header.Set("Appkey", t.Myconfig.Hk.AppKey)
	req.Header.Set("AppSecret", timestemp)
	req.Header.Set("wjwAuthorization", string(wjwAuthorization))
	req.Header.Set("Content-Type", t.ContentType)
	//library.Println(t.LogFile, "appkey:"+t.Myconfig.Hk.AppKey, "appsecret:"+timestemp, "appcode:"+t.Myconfig.Hk.AppCode, "wjwauthor:"+string(wjwAuthorization))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Panicln(t.LogFile, "发送请求出错", err.Error())
		return "", false
	}
	res, err := ioutil.ReadAll(resp.Body)
	library.Println(t.LogFile, string(res))
	resp.Body.Close()
	if err != nil {
		library.Panicln(t.LogFile, "解析海康返回数据出错", err.Error())
		return "", false
	}
	return string(res), true
}

//从海康get数据（cip）
func (t *Hikhttp) Get(url string, param map[string]string) ([]byte, bool) {
	library.Println(t.LogFile, url)
	library.Println(t.LogFile, param)
	t.ContentType = "application/json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		library.Panicln(t.LogFile, "发起请求失败：", err.Error())
		return nil, false
	}
	timestemp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	wjwAuthorization := library.Mysha256(t.Myconfig.Hk.AppKey + t.Myconfig.Hk.AppCode + timestemp)
	req.Header.Set("Appkey", t.Myconfig.Hk.AppKey)
	req.Header.Set("AppSecret", timestemp)
	req.Header.Set("wjwAuthorization", string(wjwAuthorization))
	req.Header.Set("Content-Type", t.ContentType)

	q := req.URL.Query()
	for k, v := range param {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		library.Panicln(t.LogFile, "发送请求出错：", err.Error())
		return nil, false
	}
	res, err := ioutil.ReadAll(resp.Body)
	library.Println(t.LogFile, string(res))
	resp.Body.Close()
	if err != nil {
		library.Panicln(t.LogFile, "读取返回数据出错：", err.Error())
		return nil, false
	}
	return res, true
}

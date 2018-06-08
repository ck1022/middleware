package task

/***海康用户同步程序（学生和教职工增、删、改）
*
 */
import (
	"config"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	//"log"
	"encoding/json"
	"library"
	"model"
	"net/http"
	"strings"
	"time"
)

type UserAsyn struct {
	//海康网址
	hkDomain string

	//海康用户添加地址
	hkUserAddUrl string

	//海康用户删除地址
	hkUserDeleteUrl string

	//海康用户更新地址
	hkUserUpdateUrl string

	//海康用户图片上传接口地址
	hkUserUploadPictrueUrl string

	//海康添加目标库地址
	targetAddUrl string

	//海康appkey
	hkappkey string

	//海康appcode
	hkappcode string

	//宿舍考勤机编号
	dormDevice string

	//学生同步是否开启
	studentAsynOpen string
	//教职工同步是否开启
	teacherAsynOpen string
	//智慧校园接口地址
	zhxyDomain string

	//token地址
	tokenUrl string
	//学生接口地址
	getStudentUrl string

	//教职工接口地址
	getTeacherUrl string

	//appkey
	appkey string

	//appsecret
	appsecret string

	//学校编号
	schoolCode string

	//校区编号
	campusCode string

	//校区名称
	campusName string

	//日志文件
	logFileName string
	//任务执行情况
	taskMap map[string]string
}

func (t *UserAsyn) Init() {
	t.logFileName = "log/userAsyn"
	t.taskMap = make(map[string]string)
	t.taskMap["studentAdd"] = "2006-01-02"
	t.taskMap["studentDelete"] = "2006-01-02"
	t.taskMap["studentUpdate"] = "2006-01-02"
	t.taskMap["teacherAdd"] = "2006-01-02"
	t.taskMap["teacherDelete"] = "2006-01-02"
	t.taskMap["teacherUpdate"] = "2006-01-02"
}

func (t *UserAsyn) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	//library.Println(t.logFileName, "用户信息同步程序启动")
	nowDate := time.Now().Local().Format("2006-01-02")
	t.readConfig()
	//获取token
	timestamp := time.Now().Unix()
	tokenParam := fmt.Sprintf("appkey=%s&timestamp=%d&sign=%s", t.appkey, timestamp, library.Mymd5(fmt.Sprintf("%s%d%s", t.appkey, timestamp, t.appsecret)))
	token, _, code := library.GetToken(t.tokenUrl, tokenParam)
	if code == -1 {
		library.Panicln(t.logFileName, "获取token失败")
		return
	}
	if t.studentAsynOpen == "on" {
		if t.taskMap["studentAdd"] != nowDate {
			t.taskMap["studentAdd"] = nowDate
			r := t.setStudent(token, nowDate, "add")
			if r == false {
				t.taskMap["studentAdd"] = "2006-01-02"
			}
		}
		if t.taskMap["studentDelete"] != nowDate {
			t.taskMap["studentDelete"] = nowDate
			r := t.setStudent(token, nowDate, "delete")
			if r == false {
				t.taskMap["studentDelete"] = "2006-01-02"
			}
		}
		if t.taskMap["studentUpdate"] != nowDate {
			t.taskMap["studentUpdate"] = nowDate
			r := t.setStudent(token, nowDate, "update")
			if r == false {
				t.taskMap["studentUpdate"] = "2006-01-02"
			}
		}
	}
	if t.teacherAsynOpen == "on" {
		if t.taskMap["teacherAdd"] != nowDate {
			t.taskMap["teacherAdd"] = nowDate
			r := t.setTeacher(token, nowDate, "add")
			if r == false {
				t.taskMap["teacherAdd"] = "2006-01-02"
			}
		}
		if t.taskMap["teacherDelete"] != nowDate {
			t.taskMap["teacherDelete"] = nowDate
			r := t.setTeacher(token, nowDate, "delete")
			if r == false {
				t.taskMap["teacherDelete"] = "2006-01-02"
			}
		}
		if t.taskMap["teacherUpdate"] != nowDate {
			t.taskMap["teacherUpdate"] = nowDate
			r := t.setTeacher(token, nowDate, "update")
			if r == false {
				t.taskMap["teacherUpdate"] = "2006-01-02"
			}
		}
	}
	//library.Println(t.logFileName, t.taskMap)
	*runing = false
	//library.Println(t.logFileName, "用户信息同步程序执行结束")
}

//学生和教职工全量更新
func (t *UserAsyn) AddAll() {
	t.logFileName = "log/userAsynAll"
	nowDate := "2006-01-02"
	t.readConfig()
	//获取token
	timestamp := time.Now().Unix()
	tokenParam := fmt.Sprintf("appkey=%s&timestamp=%d&sign=%s", t.appkey, timestamp, library.Mymd5(fmt.Sprintf("%s%d%s", t.appkey, timestamp, t.appsecret)))
	token, _, code := library.GetToken(t.tokenUrl, tokenParam)
	if code == -1 {
		library.Panicln(t.logFileName, "获取token失败")
		return
	}
	if t.studentAsynOpen == "on" {
		t.setStudent(token, nowDate, "add")
	}
	if t.teacherAsynOpen == "on" {
		t.setTeacher(token, nowDate, "add")
	}
}

//读取配置信息
func (t *UserAsyn) readConfig() {
	myConfig := new(config.Config)
	myConfig.InitConfig("config.txt")
	t.hkDomain = myConfig.Read("Hk", "hkDomain")
	t.hkUserAddUrl = t.hkDomain + "/eop/services/common/post/addPerson"
	if myConfig.Read("Hk", "hkUserAddUrl") != "" {
		t.hkUserAddUrl = t.hkDomain + myConfig.Read("Hk", "hkUserAddUrl")
	}
	t.hkUserDeleteUrl = t.hkDomain + "/eop/services/common/post/deletePerson"
	if myConfig.Read("Hk", "hkUserDeleteUrl") != "" {
		t.hkUserDeleteUrl = t.hkDomain + myConfig.Read("Hk", "hkUserDeleteUrl")
	}

	t.hkUserUpdateUrl = t.hkDomain + "/eop/services/common/post/updatePerson"
	if myConfig.Read("Hk", "hkUserUpdateUrl") != "" {
		t.hkUserUpdateUrl = t.hkDomain + myConfig.Read("Hk", "hkUserUpdateUrl")
	}
	t.hkUserUploadPictrueUrl = t.hkDomain + "/eop/services/common/post/uploadImage"
	if myConfig.Read("Hk", "hkUserUploadPictrueUrl") != "" {
		t.hkUserUploadPictrueUrl = t.hkDomain + myConfig.Read("Hk", "hkUserUploadPictrueUrl")
	}
	t.targetAddUrl = t.hkDomain + "/eop/services/common/post/addLeaveTarget"
	if myConfig.Read("Hk", "targetAddUrl") != "" {
		t.targetAddUrl = t.hkDomain + myConfig.Read("Hk", "targetAddUrl")
	}
	t.hkappkey = myConfig.Read("Hk", "hkappkey")
	t.hkappcode = myConfig.Read("Hk", "hkappcode")
	t.dormDevice = myConfig.Read("Hk", "dormDevice")
	t.studentAsynOpen = myConfig.Read("User", "studentAsynOpen")
	t.teacherAsynOpen = myConfig.Read("User", "teacherAsynOpen")
	t.zhxyDomain = myConfig.Read("Zhxy", "zhxyDomain")
	t.appkey = myConfig.Read("Zhxy", "appkey")
	t.appsecret = myConfig.Read("Zhxy", "appsecret")
	t.tokenUrl = t.zhxyDomain + "/api/Cert/getToken"
	if myConfig.Read("Zhxy", "tokenUrl") != "" {
		t.tokenUrl = t.zhxyDomain + myConfig.Read("Zhxy", "tokenUrl")
	}
	t.getStudentUrl = t.zhxyDomain + "/api/Student/ListsBaseInfo"
	if myConfig.Read("User", "getStudentUrl") != "" {
		t.getStudentUrl = t.zhxyDomain + myConfig.Read("User", "getStudentUrl")
	}
	t.getTeacherUrl = t.zhxyDomain + "/api/Teacher/ListsBaseInfo"
	if myConfig.Read("User", "getTeacherUrl") != "" {
		t.getTeacherUrl = t.zhxyDomain + myConfig.Read("User", "getTeacherUrl")
	}
	campus := strings.Split(myConfig.Read("Zhxy", "campus"), "|")
	t.schoolCode = campus[0]
	t.campusCode = campus[1]
	t.campusName = campus[2]
}

//设置学生
func (t *UserAsyn) setStudent(token string, date string, tag string) bool {
	library.Println(t.logFileName, "学生：", tag+"---------------------")
	theTime, _ := time.Parse("2006-01-02", date)
	yesterday := theTime.Add(-24 * time.Hour).Format("2006-01-02")
	list, err := t.getStudentList(token, yesterday, tag)
	if err {
		return false
	}
	var hklist []model.HkUser
	for _, student := range list {
		var user model.HkUser
		user.IndexCode = student.Xh
		user.Name = student.Xm
		if student.Xbm == "9" {
			student.Xbm = "0"
		}
		user.Sex = student.Xbm
		user.Picture = student.Zjzp
		user.PictureAddTime = student.Zjzptjsj
		hklist = append(hklist, user)
	}
	t.sendUserToHk(hklist, tag)
	return true
}

//设置教职工
func (t *UserAsyn) setTeacher(token string, date string, tag string) bool {
	library.Println(t.logFileName, "教职工：", tag+"---------------------")
	theTime, _ := time.Parse("2006-01-02", date)
	yesterday := theTime.Add(-24 * time.Hour).Format("2006-01-02")
	list, err := t.getTeacherList(token, yesterday, tag)
	if err {
		return false
	}
	var hklist []model.HkUser
	for _, teacher := range list {
		var user model.HkUser
		user.IndexCode = teacher.Gh
		user.Name = teacher.Xm
		if teacher.Xbm != "1" && teacher.Xbm != "2" {
			teacher.Xbm = "0"
		}
		user.Sex = teacher.Xbm
		user.Picture = teacher.Zjzp
		user.PictureAddTime = teacher.Zjzptjsj
		hklist = append(hklist, user)
	}
	library.Println(t.logFileName, hklist)
	t.sendUserToHk(hklist, tag)
	return true
}

//获取学生列表
func (t *UserAsyn) getStudentList(token string, date string, tag string) ([]model.Student, bool) {
	timestamp := time.Now().Unix()
	var timeKey = ""
	var sfsc = ""
	switch tag {
	case "add":
		timeKey = "tjsj"
		sfsc = "0"
		break
	case "delete":
		timeKey = "scsj"
		sfsc = "1"
		break
	case "update":
		timeKey = "gxsj"
		sfsc = "0"
		break
	}
	library.Println(t.logFileName, "读取学生列表："+timeKey+"="+date)
	param := fmt.Sprintf("appkey=%s&xxid=%s&xqid=%s&timestamp=%d&token=%s&sfsc=%s&%s=%s", t.appkey, t.schoolCode, t.campusCode, timestamp, token, sfsc, timeKey, date)
	req, err := http.NewRequest("POST", t.getStudentUrl, strings.NewReader(param))
	if err != nil {
		library.Panicln(t.logFileName, "读取教职工列表失败："+err.Error())
		return nil, true
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Panicln(t.logFileName, "读取学生列表失败："+err.Error())
		return nil, true
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var info = new(model.StudentReturn)
	if err != nil {
		library.Panicln(t.logFileName, "读取学生列表失败："+err.Error())
		return nil, true
	}
	json.Unmarshal(body, &info)
	code := info.Code
	list := info.Data.List
	if code == -1 {
		library.Panicln(t.logFileName, "读取学生列表失败："+string(body))
		return nil, true
	}
	library.Println(t.logFileName, "读取学生列表成功：(共", len(list), "人)：", list)
	return list, false
}

//获取教职工列表
func (t *UserAsyn) getTeacherList(token string, date string, tag string) ([]model.Teacher, bool) {
	timestamp := time.Now().Unix()
	var timeKey = ""
	var sfsc = ""
	switch tag {
	case "add":
		timeKey = "tjsj"
		sfsc = "0"
		break
	case "delete":
		timeKey = "scsj"
		sfsc = "1"
		break
	case "update":
		timeKey = "gxsj"
		sfsc = "0"
		break
	}
	library.Println(t.logFileName, "读取教职工列表："+timeKey+"="+date)
	param := fmt.Sprintf("appkey=%s&xxid=%s&xqid=%s&timestamp=%d&token=%s&sfsc=%s&%s=%s", t.appkey, t.schoolCode, t.campusCode, timestamp, token, sfsc, timeKey, date)
	req, err := http.NewRequest("POST", t.getTeacherUrl, strings.NewReader(param))
	if err != nil {
		library.Panicln(t.logFileName, "读取教职工列表失败："+err.Error())
		return nil, true
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		library.Panicln(t.logFileName, "读取教职工列表失败："+err.Error())
		return nil, true
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var info = new(model.TeacherReturn)
	if err != nil {
		library.Panicln(t.logFileName, "读取教职工列表失败："+err.Error())
		return nil, true
	}
	json.Unmarshal(body, &info)
	code := info.Code
	list := info.Data.List
	if code == -1 {
		library.Panicln(t.logFileName, "读取教职工列表失败："+string(body))
		return nil, true
	}
	library.Println(t.logFileName, "读取教职工列表成功：(共", len(list), "人)：", list)
	return list, false
}

//发送给海康
func (t *UserAsyn) sendUserToHk(list []model.HkUser, tag string) {
	type hikReturn struct {
		Code int
	}
	targetUrl := ""
	switch tag {
	case "add":
		targetUrl = t.hkUserAddUrl
		break
	case "delete":
		targetUrl = t.hkUserDeleteUrl
		break
	case "update":
		targetUrl = t.hkUserUpdateUrl
		break
	}
	if tag == "delete" { //删除人员，批量删除
		var learnCodes []string
		for _, user := range list {
			learnCodes = append(learnCodes, user.IndexCode)
		}
		if len(learnCodes) == 0 {
			library.Println(t.logFileName, "没有要删除的用户")
			return
		}
		library.Println(t.logFileName, "提交删除用户给海康："+"targeturl="+targetUrl+",indexcode="+strings.Join(learnCodes, " "))
		var map1 = make(map[string][]string)
		map1["personNoList"] = learnCodes
		b, _ := json.Marshal(map1)
		hkht := hk.Hkhttp{Appkey: t.hkappkey, Appcode: t.hkappcode}
		result, success := hkht.Send(targetUrl, b)
		if success == false {
			library.Panicln(t.logFileName, result)
		} else {
			library.Println(t.logFileName, result)
		}
	} else { //添加和修改人员，逐个修改
		for _, user := range list {
			library.Println(t.logFileName, "["+user.IndexCode+"]同步用户："+"targeturl="+targetUrl+",indexcode="+user.IndexCode+",xm="+user.Name+",zjzp="+user.Picture)
			hkPictrueUrl := ""
			if user.Picture == "" && tag == "update" { //更新时，没有设置头像，那么不更新
				library.Panicln(t.logFileName, "["+user.IndexCode+"]放弃更新：用户无照片")
				continue
			} else if tag == "update" { //更新时，如果头像没有变化或没有头像，那么不进行更新
				theTime, _ := time.Parse("2006-01-02", time.Now().Local().Format("2006-01-02"))
				yesterday := theTime.Add(-24 * time.Hour).Format("2006-01-02")
				if user.PictureAddTime < yesterday {
					library.Panicln(t.logFileName, "["+user.IndexCode+"]放弃更新：用户头像没变")
					continue
				}
			}
			if user.Picture != "" {
				picture, err := library.GetImage(user.Picture)
				if err == false {
					library.Panicln(t.logFileName, "["+user.IndexCode+"]获取用户照片失败")
				} else {
					//上传头像
					library.Println(t.logFileName, "["+user.IndexCode+"]上传照片")
					hkht := hk.Hkhttp{Appkey: t.hkappkey, Appcode: t.hkappcode}
					result, err := hkht.UpPic(t.hkUserUploadPictrueUrl, picture)
					if err == false {
						library.Panicln(t.logFileName, "["+user.IndexCode+"]"+result)
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+user.IndexCode+"]上传照片成功")
					}
				}
			}
			//上传用户信息
			library.Println(t.logFileName, "["+user.IndexCode+"]同步基本信息")
			map1 := make(map[string]string)
			map1["userIndex"] = user.IndexCode
			map1["userName"] = user.Name
			map1["sex"] = user.Sex
			map1["picturePath"] = hkPictrueUrl
			map1["idNo"] = ""
			map1["mobilephone"] = ""
			b, _ := json.Marshal(map1)
			hkht := hk.Hkhttp{Appkey: t.hkappkey, Appcode: t.hkappcode}
			result, success := hkht.Send(targetUrl, b)
			if success == false {
				library.Panicln(t.logFileName, "["+user.IndexCode+"]", map1)
				library.Panicln(t.logFileName, "["+user.IndexCode+"]"+result)
			} else {
				var returnObject = new(hikReturn)
				json.Unmarshal([]byte(result), &returnObject)
				if returnObject.Code != 20000 {
					library.Panicln(t.logFileName, "基本信息同步失败："+result)
				} else {
					library.Println(t.logFileName, "["+user.IndexCode+"]基本信息同步成功")
					t.addDeviceTarget(user.IndexCode)
				}

			}
		}
	}
}

//添加目标库
func (t *UserAsyn) addDeviceTarget(learnCode string) {

	library.Println(t.logFileName, "把名单传给海康："+"targeturl="+t.targetAddUrl+",learncode="+learnCode+",devicecode="+t.dormDevice)
	map1 := make(map[string]string)
	map1["indexCode"] = learnCode
	map1["cameraIndex"] = t.dormDevice
	b, _ := json.Marshal(map1)
	hkht := hk.Hkhttp{Appkey: t.hkappkey, Appcode: t.hkappcode}
	result, success := hkht.Send(t.targetAddUrl, b)
	if success == false {
		library.Panicln(t.logFileName, result)
	} else {
		library.Println(t.logFileName, result)
	}
}

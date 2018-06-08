package task

/***海康用户同步程序（学生和教职工增、删、改）
*
 */
import (
	"config"
	"encoding/json"
	"library"
	"library/myhttp"
	"model"
	"strconv"
	"time"
)

type StudentCip struct {
	//配置文件信息
	Config *config.MyConfig

	//日志文件
	logFileName string

	//任务执行情况
	taskMap map[string]string

	//组织机构队列
	orgMap map[string]Org

	zhxyHttp myhttp.Zhxyhttp

	hikHttp myhttp.Hikhttp
}
type OrgReturn struct {
	Code string
	Msg  string
	Data *Org
}
type Org struct {
	OrgCode    string
	OrgName    string
	ParentCode string
	OrgDefine  int
	OrgSubTree []Org
}

func (t *StudentCip) Init() {
	t.logFileName = "log/cip_studentAsyn"
	t.taskMap = make(map[string]string)
	t.taskMap["studentAsyn"] = "2006-01-02"
}

func (t *StudentCip) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	//library.Println(t.logFileName, "用户信息同步程序启动")
	nowDate := time.Now().Local().Format("2006-01-02")
	//要更新的时间，更新的是昨天的数据
	theTime, _ := time.Parse("2006-01-02", nowDate)
	date := theTime.Add(-24 * time.Hour).Format("2006-01-02")

	if t.Config.Hk96.StudentAsynOpen != "on" {
		return
	}
	t.zhxyHttp = myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	t.hikHttp = myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	//当天没有同步过
	if t.taskMap["studentAsyn"] != nowDate {
		//获取班级列表
		if t.getOrgList() == false {
			library.Panicln(t.logFileName, "获取班级列表失败")
			return
		}
		t.taskMap["studentAsyn"] = nowDate
		r := t.setStudent(date)
		if r == false {
			t.taskMap["studentAsyn"] = "2006-01-02"
		}
	}
	//library.Println(t.logFileName, t.taskMap)
	*runing = false
	//library.Println(t.logFileName, "用户信息同步程序执行结束")
}

//学生和教职工全量更新
func (t *StudentCip) AddAll() {
	t.logFileName = "log/cip_userAddAll"
	t.zhxyHttp = myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	t.hikHttp = myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	nowDate := "2006-01-02"
	//获取班级列表
	if t.getOrgList() == false {
		library.Panicln(t.logFileName, "获取班级列表失败")
		return
	}
	if t.Config.Hk96.StudentAsynOpen == "on" {
		t.setStudent(nowDate)
	}
}

//设置学生
func (t *StudentCip) setStudent(date string) bool {
	library.Println(t.logFileName, "学生：---------------------")
	//获取token
	ret := t.zhxyHttp.GetToken()
	if ret == false {
		library.Panicln(t.logFileName, "获取token失败")
		return false
	}
	list, err := t.zhxyHttp.GetStudentListCip(date)
	if err == false {
		library.Panicln(t.logFileName, "获取学生列表失败")
		return false
	}
	library.Println(t.logFileName, "获取学生列表成功，共 "+strconv.Itoa(len(list))+" 人")
	for _, student := range list {
		if student.Xbm != "1" || student.Xbm != "2" {
			student.Xbm = "1"
		}
		orgCode := t.getOrgCode(student)
		hkPictrueUrl := ""
		switch student.Isdel {
		case "0":
			if student.Zjzp != "" {
				picture, err := library.GetImage(student.Zjzp)
				if err == false {
					library.Panicln(t.logFileName, "["+student.Xnxh+"]获取用户照片失败")
				} else {
					library.Println(t.logFileName, "["+student.Xnxh+"]上传照片")
					result, err := t.hikHttp.UpFile(picture)
					if err == false {
						library.Panicln(t.logFileName, "["+student.Xnxh+"]"+result)
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+student.Xnxh+"]上传照片成功："+result)
					}
				}
			}
			res, e := t.hikHttp.AddOrUpdateStudentCip(student.Xnxh, student.Xm, student.Xbm, orgCode, hkPictrueUrl)
			if e {
				library.Println(t.logFileName, "["+student.Xnxh+"]请求成功："+res)
				//hikHttp.UpdateUserUpdateTarget(student.Xh) //更新目标库
			} else {
				library.Panicln(t.logFileName, "["+student.Xnxh+"]添加学生失败："+res)
			}
			break
		case "1":
			var ids = []string{student.Xnxh}
			t.hikHttp.DeleteStudentCip(ids)
			break
		}
	}
	return true
}

//获取系统所有org列表
func (t *StudentCip) getOrgList() bool {
	res2, err := t.hikHttp.GetOrgTree()
	if err == false {
		return false
	} else {
		var info = new(OrgReturn)
		json.Unmarshal(res2, &info)
		if info.Code != "0" {
			return false
		}
		t.orgMap = make(map[string]Org)
		t.tree2list(info.Data.OrgSubTree)
		return true
	}
	return true
}

//获取班级编号
func (t *StudentCip) getOrgCode(student model.StudentCip) string {
	if _, ok := t.orgMap[student.Njmc+"("+student.Bjmc+")班"]; !ok { //班级不存在
		//年级是否存在
		if _, ok := t.orgMap[student.Njmc]; !ok { //年级不存在
			//添加年级
			t.addGrade(student.Njmc)
		}
		//添加班级
		t.addClass(student.Njmc+"("+student.Bjmc+")班", student.Njmc)
	}
	return t.orgMap[student.Njmc+"("+student.Bjmc+")班"].OrgCode
}

//把组织机构队列转换成数组形式
func (t *StudentCip) tree2list(orgList []Org) {
	for _, v := range orgList {
		sorg := new(Org)
		sorg.OrgCode = v.OrgCode
		sorg.OrgDefine = v.OrgDefine
		sorg.OrgName = v.OrgName
		sorg.ParentCode = v.ParentCode
		t.orgMap[v.OrgName] = *sorg
		if len(v.OrgSubTree) > 0 {
			t.tree2list(v.OrgSubTree)
		}
	}
}

//创建年级
func (t *StudentCip) addGrade(gradeName string) {
	schoolOrgCode := ""
	for _, v := range t.orgMap {
		if v.OrgDefine == 1 {
			schoolOrgCode = v.OrgCode
		}
	}
	t.hikHttp.AddGradeCip(gradeName, schoolOrgCode)
	t.getOrgList()
}

//创建班级
func (t *StudentCip) addClass(className string, gradeName string) {
	t.hikHttp.AddClassCip(className, t.orgMap[gradeName].OrgCode)
	t.getOrgList()
}

//设置教职工
/*
func (t *Student) setTeacher(date string, tag string) bool {
	library.Println(t.logFileName, "教职工：", tag+"---------------------")
	zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	hikHttp := myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	//获取token
	ret := zhxyHttp.GetToken()
	if ret == false {
		library.Panicln(t.logFileName, "获取token失败")
		return false
	}
	option := make(map[string]string)
	switch tag {
	case "add":
		option["tjsj"] = date
		option["sfsc"] = "0"
		break
	case "delete":
		option["scsj"] = date
		option["sfsc"] = "1"
		break
	case "update":
		option["gxsj"] = date
		option["sfsc"] = "0"
		break
	}
	list, err := zhxyHttp.GetTeacherList(option)
	if err == false {
		return false
	}
	for _, teacher := range list {
		teacher.Gh = "G" + teacher.Gh
		if teacher.Xbm != "1" && teacher.Xbm != "2" {
			teacher.Xbm = "0"
		}
		hkPictrueUrl := ""
		switch tag {
		case "add":
			if teacher.Zjzp != "" {
				picture, err := library.GetImage(teacher.Zjzp)
				if err == false {
					library.Panicln(t.logFileName, "["+teacher.Gh+"]获取用户照片失败")
				} else {
					library.Println(t.logFileName, "["+teacher.Gh+"]上传照片")
					result, err := hikHttp.UpPic(picture)
					if err == false {
						library.Panicln(t.logFileName, "["+teacher.Gh+"]"+result)
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+teacher.Gh+"]上传照片成功："+result)
					}
				}
			}
			hikHttp.AddUser(teacher.Gh, teacher.Xm, teacher.Xbm, hkPictrueUrl, "", "")
			break
		case "delete":
			var ids = []string{teacher.Gh}
			hikHttp.DeleteUser(ids)
			break
		case "update":
			library.Println(t.logFileName, "["+teacher.Gh+"]"+teacher.Zjzp)
			if teacher.Zjzp != "" {
				picture, err := library.GetImage(teacher.Zjzp)
				if err == false {
					library.Panicln(t.logFileName, "["+teacher.Gh+"]获取用户照片失败")
				} else {
					library.Println(t.logFileName, "["+teacher.Gh+"]上传照片")
					result, err := hikHttp.UpPic(picture)
					if err == false {
						library.Panicln(t.logFileName, "["+teacher.Gh+"]"+result)
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+teacher.Gh+"]上传照片成功："+result)
						hikHttp.UpdateUser(teacher.Gh, teacher.Xm, teacher.Xbm, hkPictrueUrl, "", "")
					}
				}
			} else {
				library.Println(t.logFileName, "["+teacher.Gh+"]用户没有照片，放弃更新")
			}
			break
		}
	}
	return true
}
*/

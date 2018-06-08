package task

/***海康用户同步程序（学生和教职工增、删、改）
*
 */
import (
	"config"
	"library"
	"library/myhttp"
	"strconv"
	"time"
)

type UserAsyn struct {
	//配置文件信息
	Config *config.MyConfig

	//日志文件
	logFileName string

	//任务执行情况
	taskMap map[string]string
}

func (t *UserAsyn) Init() {
	t.logFileName = "log/96_userAsyn"
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
	//要更新的时间，更新的是昨天的数据
	theTime, _ := time.Parse("2006-01-02", nowDate)
	date := theTime.Add(-24 * time.Hour).Format("2006-01-02")

	if t.Config.Hk96.StudentAsynOpen == "on" {
		if t.taskMap["studentAdd"] != nowDate {
			t.taskMap["studentAdd"] = nowDate
			r := t.setStudent(date, "add")
			if r == false {
				t.taskMap["studentAdd"] = "2006-01-02"
			}
		}
		if t.taskMap["studentDelete"] != nowDate {
			t.taskMap["studentDelete"] = nowDate
			r := t.setStudent(date, "delete")
			if r == false {
				t.taskMap["studentDelete"] = "2006-01-02"
			}
		}
		if t.taskMap["studentUpdate"] != nowDate {
			t.taskMap["studentUpdate"] = nowDate
			r := t.setStudent(date, "update")
			if r == false {
				t.taskMap["studentUpdate"] = "2006-01-02"
			}
		}
	}
	if t.Config.Hk96.TeacherAsynOpen == "on" {
		if t.taskMap["teacherAdd"] != nowDate {
			t.taskMap["teacherAdd"] = nowDate
			r := t.setTeacher(date, "add")
			if r == false {
				t.taskMap["teacherAdd"] = "2006-01-02"
			}
		}
		if t.taskMap["teacherDelete"] != nowDate {
			t.taskMap["teacherDelete"] = nowDate
			r := t.setTeacher(date, "delete")
			if r == false {
				t.taskMap["teacherDelete"] = "2006-01-02"
			}
		}
		if t.taskMap["teacherUpdate"] != nowDate {
			t.taskMap["teacherUpdate"] = nowDate
			r := t.setTeacher(date, "update")
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
	if t.Config.Hk96.StudentAsynOpen == "on" {
		t.setStudent(nowDate, "add")
	}
	if t.Config.Hk96.TeacherAsynOpen == "on" {
		t.setTeacher(nowDate, "add")
	}
}

//设置学生
func (t *UserAsyn) setStudent(date string, tag string) bool {
	library.Println(t.logFileName, "学生：", tag+"---------------------")
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
	list, err := zhxyHttp.GetStudentList(option)
	if err == false {
		library.Panicln(t.logFileName, "获取学生列表失败")
		return false
	}
	library.Println(t.logFileName, "获取学生列表成功，共 "+strconv.Itoa(len(list))+" 人")
	for _, student := range list {
		if student.Xbm != "1" || student.Xbm != "2" {
			student.Xbm = "1"
		}
		hkPictrueUrl := ""
		switch tag {
		case "add":
			if student.Zjzp != "" {
				picture, err := library.GetImage(student.Zjzp)
				if err == false {
					library.Panicln(t.logFileName, "["+student.Xh+"]获取用户照片失败")
				} else {
					library.Println(t.logFileName, "["+student.Xh+"]上传照片")
					result, err := hikHttp.UpPic(picture)
					if err == false {
						library.Panicln(t.logFileName, "["+student.Xh+"]上传照片失败")
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+student.Xh+"]上传照片成功")
					}
				}
			}
			_, e := hikHttp.AddUser(student.Xh, student.Xm, student.Xbm, hkPictrueUrl, "", "")
			if e {
				library.Println(t.logFileName, "["+student.Xh+"]添加用户成功")
				time.Sleep(1 * time.Second)
				hikHttp.UpdateUserUpdateTarget(student.Xh) //更新目标库
			} else {
				library.Panicln(t.logFileName, "["+student.Xh+"]添加用户失败")
			}

			break
		case "delete":
			var ids = []string{student.Xh}
			hikHttp.DeleteUser(ids)
			break
		case "update":
			if student.Zjzp != "" {
				picture, err := library.GetImage(student.Zjzp)
				if err == false {
					library.Panicln(t.logFileName, "["+student.Xh+"]获取用户照片失败")
				} else {
					library.Println(t.logFileName, "["+student.Xh+"]上传照片")
					result, err := hikHttp.UpPic(picture)
					if err == false {
						library.Panicln(t.logFileName, "["+student.Xh+"]上传照片失败")
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+student.Xh+"]上传照片成功")
						_, e := hikHttp.UpdateUser(student.Xh, student.Xm, student.Xbm, hkPictrueUrl, "", "")
						if e {
							library.Println(t.logFileName, "["+student.Xh+"]修改用户成功")
							time.Sleep(1 * time.Second)
							hikHttp.UpdateUserUpdateTarget(student.Xh) //更新目标库
						} else {
							library.Panicln(t.logFileName, "["+student.Xh+"]修改用户失败")
						}
					}
				}
			} else {
				library.Println(t.logFileName, "["+student.Xh+"]用户没有照片，放弃更新")
			}
			break
		}
	}
	return true
}

//设置教职工
func (t *UserAsyn) setTeacher(date string, tag string) bool {
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
			teacher.Xbm = "1"
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
						library.Panicln(t.logFileName, "["+teacher.Gh+"]上传照片失败")
					} else {
						hkPictrueUrl = result
						library.Println(t.logFileName, "["+teacher.Gh+"]上传照片成功")
					}
				}
			}
			hikHttp.AddUser(teacher.Gh, teacher.Xm, teacher.Xbm, hkPictrueUrl, "", "")
			time.Sleep(1 * time.Second)
			hikHttp.UpdateTeacherUpdateTarget(teacher.Gh) //更新目标库
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
						library.Println(t.logFileName, "["+teacher.Gh+"]上传照片成功")
						hikHttp.UpdateUser(teacher.Gh, teacher.Xm, teacher.Xbm, hkPictrueUrl, "", "")
						time.Sleep(1 * time.Second)
						hikHttp.UpdateTeacherUpdateTarget(teacher.Gh) //更新目标库
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

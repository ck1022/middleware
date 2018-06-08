package main

import (
	"config"
	//"fmt"
	//"io/ioutil"
	//"log"
	"encoding/json"
	"library"
	"library/myhttp"
	"model"
	//"os"
	//"strings"
	"strconv"
	"time"
)

//配置信息
var Config *config.MyConfig

var zhxyHttp myhttp.Zhxyhttp
var hikHttp myhttp.Hikhttp

//日志文件
var logFileName = "log/96_MessageQueue"

type messageReturn struct {
	Code int
	Data struct {
		List []message
	}
}
type message struct {
	Queueid string
	Topic   string
	Tag     string
	Body    string
}

func main() {
	//日志目录是否存在
	library.Mkdir("log")
	library.Println(logFileName, "开始处理消息队列")
	taskAliveLastTime := time.Now().Unix() - 10000
	for {
		if time.Now().Unix()-taskAliveLastTime > 60 {
			Config = new(config.MyConfig)
			Config.ReadConfig()
			zhxyHttp = myhttp.Zhxyhttp{Myconfig: Config, LogFile: logFileName}
			hikHttp = myhttp.Hikhttp{Myconfig: Config, LogFile: logFileName}

			zhxyHttp.TaskAlive("message96")
			//获取token
			ret := zhxyHttp.GetToken()
			if ret == false {
				library.Panicln(logFileName, "获取token失败")
				continue
			}
			taskAliveLastTime = time.Now().Unix()
		}
		library.Println(logFileName, "--------------------------------------------------------------------------------")

		list, _ := zhxyHttp.GetMessageQueueList()
		if list != nil {
			if len(list) > 0 {
				library.Println(logFileName, "新的消息("+strconv.Itoa(len(list))+")", list)
			}
			for _, m := range list {
				library.Println(logFileName, "处理消息：id="+m.Queueid+",topic="+m.Topic+",tag="+m.Tag+",body="+m.Body)
				switch m.Topic {
				case "student": //学生资料变更
					student(m)
					break
				case "studentLeave": //学生请假申请
					qingjia(m)
					break
				case "studentResidentThrough": //学生通校申请
					tongxiao(m)
					break
				case "studentEarlyBack": //学生提前返校申请
					fanxiao(m)
					break
				case "studentSchoolDoorSwingCard": //学生进出校门记录
					enterLeaveSchoolDoor(m)
					break
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

//学生资料变更处理
func student(m model.Message) {
	var body struct {
		Xsid string
	}
	json.Unmarshal([]byte(m.Body), &body)
	switch m.Tag {
	case "edit": //修改学生
		option := make(map[string]string)
		option["xsid"] = body.Xsid
		list, err := zhxyHttp.GetStudentList(option)
		library.Println(logFileName, "学生更新：", list)
		if err == false {
			break
		}
		for _, student := range list {
			if student.Zjzp != "" {
				if student.Xbm == "9" {
					student.Xbm = "0"
				}
				picture, err := library.GetImage(student.Zjzp)
				if err == false {
					library.Panicln(logFileName, "["+student.Xh+"]获取用户照片失败")
				} else {
					library.Println(logFileName, "["+student.Xh+"]上传照片")
					result, err := hikHttp.UpPic(picture)
					if err == false {
						library.Panicln(logFileName, "["+student.Xh+"]上传照片失败："+result)
					} else {
						hkPictrueUrl := result
						library.Println(logFileName, "["+student.Xh+"]上传照片成功："+result)
						res, e := hikHttp.UpdateUser(student.Xh, student.Xm, student.Xbm, hkPictrueUrl, "", "")
						if e {
							library.Println(logFileName, "["+student.Xh+"]更新用户成功："+res)
							hikHttp.UpdateUserUpdateTarget(student.Xh) //更新目标库
						} else {
							library.Panicln(logFileName, "["+student.Xh+"]更新用户失败："+res)
						}
					}
				}
			} else {
				library.Println(logFileName, "["+student.Xh+"]用户没有照片，放弃更新")
			}
		}
		break
	}
}

//请假消息处理，进和出都加
func qingjia(m model.Message) {
	var body struct {
		LearnCode string
	}
	json.Unmarshal([]byte(m.Body), &body)
	res := ""
	switch m.Tag {
	case "add": //添加请假
		if Config.Hk96.InOpen == "on" { //进校不控制，那么不进行处理
			library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.InDevice+"]添加请假进校名单")
			res, _ = hikHttp.AddTarget(Config.Hk96.InDevice, body.LearnCode)
			library.Println(logFileName, "["+body.LearnCode+"]添加结果："+res)
		}
		if Config.Hk96.OutOpen == "on" { //出校不控制，那么不进行处理
			library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.OutDevice+"]添加请假出校名单")
			res, _ = hikHttp.AddTarget(Config.Hk96.OutDevice, body.LearnCode)
			library.Println(logFileName, "["+body.LearnCode+"]添加结果："+res)
		}
		break
	case "delete": //删除请假
		if Config.Hk96.InOpen == "on" { //进校不控制，那么不进行处理
			library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.OutDevice+"]删除请假进校名单")
			res, _ = hikHttp.DeleteTarget(Config.Hk96.InDevice, body.LearnCode)
			library.Println(logFileName, "["+body.LearnCode+"]删除结果："+res)
		}
		if Config.Hk96.OutOpen == "on" { //出校不控制，那么不进行处理
			library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.OutDevice+"]删除请假出校名单")
			res, _ = hikHttp.DeleteTarget(Config.Hk96.OutDevice, body.LearnCode)
			library.Println(logFileName, "["+body.LearnCode+"]删除结果："+res)
		}
		break
	}
}

//通校消息处理，只加出名单
func tongxiao(m model.Message) {
	var body struct {
		LearnCode string
	}
	json.Unmarshal([]byte(m.Body), &body)
	res := ""
	switch m.Tag {
	case "add": //添加通校
		library.Println(logFileName, "["+body.LearnCode+"] topic=studentResidentThrough tag=add")
		library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.OutDevice+"]添加通校出校名单")
		res, _ = hikHttp.AddTarget(Config.Hk96.OutDevice, body.LearnCode)
		library.Println(logFileName, "["+body.LearnCode+"]添加结果："+res)
		break
	}
}

//返校消息处理，只加进名单
func fanxiao(m model.Message) {
	var body struct {
		LearnCode string
	}
	json.Unmarshal([]byte(m.Body), &body)
	res := ""
	switch m.Tag {
	case "add": //添加返校
		library.Println(logFileName, "["+body.LearnCode+"] topic=studentEarlyBack tag=add")
		library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.InDevice+"]添加提前返校名单")
		res, _ = hikHttp.DeleteTarget(Config.Hk96.InDevice, body.LearnCode)
		library.Println(logFileName, "["+body.LearnCode+"]添加结果："+res)
		break
	}
}

//学生进出校门
func enterLeaveSchoolDoor(m model.Message) {
	var body struct {
		LearnCode string
	}
	json.Unmarshal([]byte(m.Body), &body)
	res := ""
	switch m.Tag {
	case "enter": //进
		if Config.Hk96.InOpen == "on" { //进校不控制，那么不进行处理
			library.Println(logFileName, "["+body.LearnCode+"] topic=studentSchoolDoorSwingCard tag=enter")
			library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.InDevice+"删除进校名单")
			res, _ = hikHttp.DeleteTarget(Config.Hk96.InDevice, body.LearnCode)
			library.Println(logFileName, "["+body.LearnCode+"]删除结果："+res)
		}
		break
	case "leave": //出
		if Config.Hk96.OutOpen == "on" { //出校不控制，那么不进行处理
			library.Println(logFileName, "["+body.LearnCode+"] topic=studentSchoolDoorSwingCard tag=leave")
			library.Println(logFileName, "["+body.LearnCode+"]["+Config.Hk96.OutDevice+"]删除出校名单")
			res, _ = hikHttp.DeleteTarget(Config.Hk96.OutDevice, body.LearnCode)
			library.Println(logFileName, "["+body.LearnCode+"]删除结果："+res)
			break
		}
	}
}

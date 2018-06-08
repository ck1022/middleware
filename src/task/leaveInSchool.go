package task

/** 海康进出校门走读和临时通校名单设置程序同步
*
 */
import (
	"config"
	//"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	//"io/ioutil"
	//"log"
	//"model"
	//"encoding/json"
	"library"
	"library/myhttp"
	//"model"
	//"net/http"
	//"os"
	"strconv"
	"strings"
	"time"
)

type LeaveInSchool struct {
	//配置文件信息
	Config *config.MyConfig

	//日志文件
	logFileName string

	//任务执行情况
	taskMap map[string]string
}

//构造函数
func (t *LeaveInSchool) Init() {
	t.logFileName = "log/96_leaveInSchool"
	t.taskMap = make(map[string]string)
	t.taskMap["inAdd"] = "2006-01-02"
	t.taskMap["inDelete"] = "2006-01-02"
	t.taskMap["outAdd"] = "2006-01-02"
	t.taskMap["outDelete"] = "2006-01-02"
}

//开始走读和临时通校的学生同步程序
func (t *LeaveInSchool) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	//library.Println(t.logFileName, "走读和通校学生目标库维护程序启动")
	nowDate := time.Now().Local().Format("2006-01-02")
	nowMinite := time.Now().Local().Format("15:04")
	weekDay := strconv.Itoa(t.getWeekDay())
	if strings.Contains(t.Config.Hk96.InWeekDay, weekDay) && t.Config.Hk96.InOpen == "on" {
		//允许进校名单设置
		if nowMinite >= t.Config.Hk96.InAddTime && t.taskMap["inAdd"] != nowDate {

			if t.setDeviceStudent(t.Config.Hk96.InDevice, true, true) {
				t.taskMap["inAdd"] = nowDate
			}
		}
		//删除进校名单
		if nowMinite >= t.Config.Hk96.InDeleteTime && t.taskMap["inDelete"] != nowDate {
			if t.setDeviceStudent(t.Config.Hk96.InDevice, false, true) {
				t.taskMap["inDelete"] = nowDate
			}
		}
	}
	library.Println(t.logFileName, "--------------------", t.Config.Hk96.OutAddTime, "---------------------", t.Config.Hk96.OutWeekDay, "-----", t.Config.Hk96.OutOpen)

	if strings.Contains(t.Config.Hk96.OutWeekDay, weekDay) && t.Config.Hk96.OutOpen == "on" {
		library.Println(t.logFileName, "--------------------", t.Config.Hk96.OutAddTime, "---------------------")
		//允许离校名单设置
		if nowMinite >= t.Config.Hk96.OutAddTime && t.taskMap["outAdd"] != nowDate {

			if t.setDeviceStudent(t.Config.Hk96.OutDevice, true, false) {
				t.taskMap["outAdd"] = nowDate
			}
		}
		//删除离校名单
		if nowMinite >= t.Config.Hk96.OutDeleteTime && t.taskMap["outDelete"] != nowDate {
			if t.setDeviceStudent(t.Config.Hk96.OutDevice, false, false) {
				t.taskMap["outDelete"] = nowDate
			}
		}
	}
	//library.Println(t.logFileName, "走读和通校学生目标库维护程序执行结束")
	*runing = false
}

//设置设备名单,isAdd：true表示添加，false表示删除
func (t *LeaveInSchool) setDeviceStudent(deviceCode string, add bool, in bool) bool {
	library.Println(t.logFileName, "设置名单：deviceCode="+deviceCode+",add="+strconv.FormatBool(add)+",in="+strconv.FormatBool(in))
	zhxyhttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	hikHttp := myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	//获取token
	ret := zhxyhttp.GetToken()
	if ret == false {
		library.Panicln(t.logFileName, "获取token失败")
		return false
	}
	date := time.Now().Format("2006-01-02")
	//请假人员名单
	var qingjiaStudentLearnCodes []string
	if add == false {
		qjlist, ret := zhxyhttp.GetLeaveStudentList(date)
		if ret == false {
			return false
		}
		for _, s := range qjlist {
			qingjiaStudentLearnCodes = append(qingjiaStudentLearnCodes, s.Xh)
		}
	}
	//读取走读和通校人员列表
	zouduList, ret := zhxyhttp.GetLeaveInStudentList(date, in)
	if zouduList == nil {
		return false
	}
	//逐个发送给海康设备
	var targetLearnCodes []string
	for i, s := range zouduList {
		//删除名单时,如果学生在请假名单中，那么不删除名单
		if add == false && library.InArray(s.Xh, qingjiaStudentLearnCodes) {
			library.Println(t.logFileName, "学号="+s.Xh+"，姓名="+s.Xm+",在请假名单中，不处理")
			continue
		}
		targetLearnCodes = append(targetLearnCodes, s.Xh)
		if (i+1)%10 == 0 || (i+1) == len(zouduList) {
			library.Println(t.logFileName, "学号列表："+strings.Join(targetLearnCodes, ","))
			if add {
				hikHttp.AddTarget(deviceCode, strings.Join(targetLearnCodes, ","))
			} else {
				hikHttp.DeleteTarget(deviceCode, strings.Join(targetLearnCodes, ","))
			}
			targetLearnCodes = nil
		}
	}
	return true
}

//获取当前星期几（数字）
func (t *LeaveInSchool) getWeekDay() int {
	weekDayString := time.Now().Weekday().String()
	weekDayMap := map[string]int{
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}
	return weekDayMap[weekDayString]
}

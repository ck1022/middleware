package task

/***海康用户同步程序（学生和教职工增、删、改）
*
 */
import (
	"config"
	"library"
	"library/myhttp"
	//"strconv"
	"encoding/json"
	//"fmt"
	//"strings"
	"time"
)

type StudentClassRoomAttendance struct {
	//配置文件信息
	Config *config.MyConfig

	//日志文件
	logFileName string

	//任务执行情况
	taskMap map[string]string
	//时间跨度
	timeStep string
}

//海康学生行为接口返回
type StudentClassRoomAttendanceReturn struct {
	Code    string          //返回编号
	Success bool            //是否成功
	Data    *AttendanceData //行为
}

//海康数据区内容
type AttendanceData struct {
	Content          []Attendance //行为数据
	TotalElements    int          //总共多少条
	Last             bool
	TotalPages       int //总共页数
	Size             int //每页几条
	NumberOfElements int //多少条
}

//海康返回识别信息
type Attendance struct {
	IndexCode  string //设备编号
	ReportTime string //记录时间
	PersonCode string //人员编号
	PersonName string //人员姓名
	Similarity string //人员姓名
}

func (t *StudentClassRoomAttendance) Init() {
	t.logFileName = "log/cip_studentClassRoomAttendance"
	t.timeStep = "+1m"
}

func (t *StudentClassRoomAttendance) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	pageSize := 100 //每页几条
	/**是否开启同步**/
	if t.Config.HkCip.ClassRoomAttendanceOpen != "on" {
		*runing = false
		return
	}
	m, _ := time.ParseDuration("-3m") //从1小时前开始读取
	startTime := time.Now().Add(m)
	m1, _ := time.ParseDuration(t.timeStep) //每次读取1分钟数据
	endTime := startTime.Add(m1)
	library.Println(t.logFileName, startTime.Format("2006-01-02 15:04:00"), endTime.Format("2006-01-02 15:04:00"))
	for {
		/**如果时间没到截止时间，那么等待一分钟，再执行**/
		if endTime.Format("2006-01-02 15:04:00") >= time.Now().Format("2006-01-02 15:04:00") {
			time.Sleep(60 * time.Second)
			continue
		}
		zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
		zhxyHttp.TaskAlive("student_classroom_attendance")
		for page := 1; page <= 10000; page++ {
			data, e := t.GetAttendanceList(startTime.Format("2006-01-02 15:04:00"), endTime.Format("2006-01-02 15:04:00"), page, pageSize)
			if !e { //出错，可能是断网
				library.Println(t.logFileName, startTime.Format("2006-01-02 15:04:00"), endTime.Format("2006-01-02 15:04:00"), "获取海康行为数据失败")
				time.Sleep(60 * time.Second)
				break
			}
			list := data.Content
			if len(list) > 0 {
				library.Println(t.logFileName, "共(", page, "/", data.TotalPages, ")页")
				library.Println(t.logFileName, "共(", len(list), "/", data.TotalElements, ")条")
			} else {
				break
			}
			t.UploadDeal(list)
			if data.Last {
				break
			}
		}
		m, _ := time.ParseDuration("-1s") //多读10秒的数据，避免遗漏
		startTime = endTime.Add(m)
		m1, _ := time.ParseDuration(t.timeStep) //每次读取跨度2分钟的数据
		endTime = endTime.Add(m1)
	}
	*runing = false
	//library.Println(t.logFileName, "用户信息同步程序执行结束")
}

//获取行为列表
func (t *StudentClassRoomAttendance) GetAttendanceList(starTime string, endTime string, page int, pageSize int) (*AttendanceData, bool) {
	/**获取学生行为数据列表**/
	hikHttp := myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	b, e := hikHttp.GetStudentClassRoomAttendanceData(starTime, endTime, page, pageSize)
	library.Println(t.logFileName, string(b))
	if e == false {
		return nil, false
	}
	var info = new(StudentClassRoomAttendanceReturn)
	json.Unmarshal(b, &info)
	if info.Code != "0" {
		return nil, false
	}
	data := info.Data
	return data, true
}

/**
*上传到美智
 */
func (t *StudentClassRoomAttendance) UploadDeal(list []Attendance) {
	zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	zhxyHttp.GetToken()
	for _, v := range list {
		library.Println(t.logFileName, "开始处理")
		//上传到美智
		param := make(map[string]string)
		param["cameraIndex"] = v.IndexCode
		param["learnCode"] = v.PersonCode
		param["studentName"] = v.PersonName
		param["recordTime"] = v.ReportTime
		s, e := zhxyHttp.UploadClassRoomAttendance(param)
		if e {
			library.Println(t.logFileName, "上传到服务器成功", string(s))
		} else {
			library.Panicln(t.logFileName, "上传到服务器失败", string(s))
		}
	}
}

package task

import (
	"config"
	"encoding/json"
	"fmt"
	"library"
	"library/myhttp"
	"time"
)

type StudentBehavea struct {
	Config      *config.MyConfig
	logFileName string
	timeStep    string
}

//海康学生行为接口返回
type BehaveaReturn struct {
	Code    string       //返回编号
	Success bool         //是否成功
	Data    *BehaveaData //行为
}

//海康数据区内容
type BehaveaData struct {
	Content          []Behavea //行为数据
	TotalElements    int       //总共多少条
	Last             bool
	TotalPages       int //总共页数
	Size             int //每页几条
	NumberOfElements int //多少条
}

//海康返回行为信息
type Behavea struct {
	CameraIndex   string //设备编号
	RecordTime    string //记录时间
	BehaviourType string //行为类型
	PartPicUrl    string //图片
	X1            int
	Y1            int
	X2            int
	Y2            int
}

func (t *StudentBehavea) Init() {
	t.logFileName = "log/cip_studentBehavea"
	t.timeStep = "+1m"
}

func (t *StudentBehavea) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	pageSize := 100 //每页几条
	/**是否开启同步**/
	if t.Config.HkCip.BehaveaOpen != "on" {
		*runing = false
		return
	}
	m, _ := time.ParseDuration("-30m") //从1小时前开始读取
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
		zhxyHttp.TaskAlive("student_behavea")
		for page := 1; page <= 50; page++ {
			data, e := t.GetBehaveaList(startTime.Format("2006-01-02 15:04:00"), endTime.Format("2006-01-02 15:04:00"), page, pageSize)
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
}

//获取行为列表
func (t *StudentBehavea) GetBehaveaList(starTime string, endTime string, page int, pageSize int) (*BehaveaData, bool) {
	/**获取学生行为数据列表**/
	hikHttp := myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	b, e := hikHttp.GetStudentBehaveData(starTime, endTime, page, pageSize)
	if e == false {
		return nil, false
	}
	var info = new(BehaveaReturn)
	json.Unmarshal(b, &info)
	if info.Code != "0" {
		return nil, false
	}
	data := info.Data
	return data, true
}

/**上传到美智
*没有图片证据的不处理
 */
func (t *StudentBehavea) UploadDeal(list []Behavea) {
	zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	for _, v := range list {
		library.Println(t.logFileName, "开始处理")
		pic := ""
		if v.PartPicUrl == "" { //无图无真相
			continue
		}
		//下载识别比对图片
		_, e := library.GetImageNew(v.PartPicUrl, "StudentBehavea.jpg")
		if e {
			library.Println(t.logFileName, "从海康下载图片成功")
			pic = "StudentBehavea.jpg"
		} else {
			library.Panicln(t.logFileName, "从海康下载图片失败")
			continue
		}
		//上传到美智
		param := make(map[string]string)
		param["cameraIndex"] = v.CameraIndex
		param["recordTime"] = v.RecordTime
		param["behaviour"] = v.BehaviourType
		param["x1"] = fmt.Sprintf("%d", v.X1)
		param["y1"] = fmt.Sprintf("%d", v.Y1)
		param["x2"] = fmt.Sprintf("%d", v.X2)
		param["y2"] = fmt.Sprintf("%d", v.Y2)
		//return
		s, e := zhxyHttp.UploadWithFile(t.Config.HkCip.UploadStudentBehaveDataUrl, pic, param)
		if e {
			library.Println(t.logFileName, "上传到服务器成功")
		} else {
			library.Panicln(t.logFileName, "上传到服务器失败", string(s))
		}
	}
}

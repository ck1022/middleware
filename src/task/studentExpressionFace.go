package task

import (
	"config"
	//"fmt"
	"encoding/json"
	"fmt"
	"library"
	"library/myhttp"
	"time"
)

type StudentExpressionFace struct {
	Config      *config.MyConfig
	logFileName string
	timeStep    string
}

//海康学生表情接口返回
type ExpressionFaceReturn struct {
	Code    string              //返回编号
	Success bool                //是否成功
	Data    *ExpressionFaceData //行为
}

//海康数据区内容
type ExpressionFaceData struct {
	Content          []ExpressionFace //行为数据
	TotalElements    int              //总共多少条
	Last             bool
	TotalPages       int //总共页数
	Size             int //每页几条
	NumberOfElements int //多少条
}

//海康返回表情信息
type ExpressionFace struct {
	CameraIndex    string //设备编号
	StudentNo      string //学号
	StudentName    string //学生姓名
	RecordTime     string //记录时间
	ExpressionType string //表情类型
	FacePtz        string
	PartPicUrl     string //图片
	X1             int
	Y1             int
	X2             int
	Y2             int
}

func (t *StudentExpressionFace) Init() {
	t.logFileName = "log/cip_studentExpressionFace"
	t.timeStep = "+2m"
}
func (t *StudentExpressionFace) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	pageSize := 10 //每页几条
	if t.Config.HkCip.ExpressionOpen != "on" {
		*runing = false
		return
	}
	m, _ := time.ParseDuration("-40m") //新启动的从40分钟前开始执行
	startTime := time.Now().Add(m)
	m1, _ := time.ParseDuration(t.timeStep) //每次读取1分钟数据
	endTime := startTime.Add(m1)
	library.Println(t.logFileName, startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))
	for {
		if endTime.Format("2006-01-02 15:04") >= time.Now().Format("2006-01-02 15:04") {
			time.Sleep(60 * time.Second)
			continue
		}
		zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
		zhxyHttp.TaskAlive("student_expression")
		for page := 1; page <= 10000; page++ {
			data, e := t.GetExpressionFaceList(startTime.Format("2006-01-02 15:04:00"), endTime.Format("2006-01-02 15:04:00"), page, pageSize)
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
		m1, _ := time.ParseDuration(t.timeStep) //每次读取跨度1分钟的数据
		endTime = endTime.Add(m1)
	}
	*runing = false
}

//获取表情列表
func (t *StudentExpressionFace) GetExpressionFaceList(starTime string, endTime string, page int, pageSize int) (*ExpressionFaceData, bool) {
	hikHttp := myhttp.Hikhttp{Myconfig: t.Config, LogFile: t.logFileName}
	b, e := hikHttp.GetStudentExpressionFaceData(starTime, endTime, page, pageSize)
	if e == false {
		return nil, false
	}
	var info = new(ExpressionFaceReturn)
	json.Unmarshal(b, &info)
	if info.Code != "0" {
		return nil, false
	}
	data := info.Data
	return data, true
}

func (t *StudentExpressionFace) UploadDeal(list []ExpressionFace) {
	zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	for _, v := range list {
		library.Println(t.logFileName, "开始处理")
		pic := ""
		if v.PartPicUrl == "" { //无图无真相
			continue
		}
		//下载识别比对图片
		_, e := library.GetImageNew(v.PartPicUrl, "StudentExpressionFace.jpg")
		if e {
			library.Println(t.logFileName, "从海康下载图片成功")
			pic = "StudentExpressionFace.jpg"
		} else {
			library.Panicln(t.logFileName, "从海康下载图片失败")
			continue
		}
		//上传到美智
		param := make(map[string]string)
		param["cameraIndex"] = v.CameraIndex
		param["learnCode"] = v.StudentNo
		param["studentName"] = v.StudentName
		param["recordTime"] = v.RecordTime
		param["expressionType"] = v.ExpressionType
		param["x1"] = fmt.Sprintf("%d", v.X1)
		param["y1"] = fmt.Sprintf("%d", v.Y1)
		param["x2"] = fmt.Sprintf("%d", v.X2)
		param["y2"] = fmt.Sprintf("%d", v.Y2)
		library.Println(t.logFileName, param)
		s, e := zhxyHttp.UploadWithFile(t.Config.HkCip.UploadStudentExpressionFaceDataUrl, pic, param)
		if e {
			library.Println(t.logFileName, "上传到服务器成功")
		} else {
			library.Panicln(t.logFileName, "上传到服务器失败", string(s))
		}
	}
}

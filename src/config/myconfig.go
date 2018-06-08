package config

import (
	"strings"
)

type MyConfig struct {
	Hk           *Hk
	Zhxy         *Zhxy
	Hk96         *Hk96
	MessageQueue *MessageQueue
	HkCip        *HkCip
}

//海康配置
type Hk struct {
	//海康appkey
	AppKey string
	//海康code
	AppCode string
}

//基础CIP信息配置
type HkCip struct {
	Domain                                  string //海康cip接口域名
	StudentAsynOpen                         string //学生同步资料是否开启
	TeacherAsynOpen                         string //教职工同步资料是否开启
	BehaveaOpen                             string //行为同步开启
	ExpressionOpen                          string //表情同步开启
	ClassRoomAttendanceOpen                 string //学生班级考勤信息同步
	UploadFileUrl                           string //上传文件地址（cip）
	GetStudentExpressionFaceDataUrl         string //获取学生表情数据地址
	UploadStudentExpressionFaceDataUrl      string //上传行为同步数据地址
	GetStudentBehaveDataUrl                 string //获取行为同步数据地址
	UploadStudentBehaveDataUrl              string //上传行为同步数据地址
	GetStudentClassRoomAttendanceUrl        string //获取考勤同步数据地址
	UploadStudentClassRoomAttendanceDataUrl string //上传考勤同步数据地址

	StudentAddOrUpdateUrl string //学生添加或修改地址
	StudentDeleteUrl      string //学生删除
	TeacherAddOrUpdateUrl string //教职工添加或修改
	TeacherDeleteUrl      string //教职工删除
	OrgTreeUrl            string //获取组织树地址
	GradeAddUrl           string //添加年级
	ClassAddUrl           string //添加班级
	TestDormDevice        string //测试用宿舍明眸设备，用于重点目标库更新
}

//海康96信息配置
type Hk96 struct {
	Domain                string //96平台接口域名
	StudentAsynOpen       string //学生同步是否开启
	TeacherAsynOpen       string //教职工同步是否开启
	UploadFaceOpen        string //刷脸照片是否上传
	LeaveInOutOpen        string //是否有进出校设备
	InOpen                string //进校是否开启
	OutOpen               string //进校是否开启
	GetStudentUrl         string //获取学生资料地址
	GetTeacherUrl         string //获取教职工资料地址
	GetThoroughStudentUrl string //获取走读和同校学生名单
	GetLeaveInStudenturl  string //获取进出校学生名单
	GetLeaveStudenturl    string //获取请假学生名单
	GetAttendanceFaceList string //人脸比对图片未处理记录列表
	UploadFaceImageUrl    string //上传人脸比对图片

	TargetAddUrl         string //目标库添加接口
	TargetDeleteUrl      string //目标库删除接口
	UserUploadPictrueUrl string //上传照片地址(原用户上传地址)
	UserAddUrl           string //用户添加接口
	UserDeleteUrl        string //用户删除接口
	UserUpdateUrl        string //用户更新接口
	InWeekDay            string //出校考勤周几，多个用半角逗号分隔，如1,2,3,4
	OutWeekDay           string //进校考勤周几，多个用半角逗号分隔，如1,2,3,4
	InDevice             string //进校明眸设备编号
	OutDevice            string //出校明眸设备编号
	InAddTime            string //进校名单添加时间
	InDeleteTime         string //进校名单删除时间
	OutAddTime           string //出校名单添加时间
	OutDeleteTime        string //出校名单删除时间
	TestDormDevice       string //测试用宿舍明眸设备，用于重点目标库更新
}

//智慧校园配置
type Zhxy struct {
	AppKey     string //appkey
	AppSecret  string //appsecret
	Domain     string //接口域名地址
	SchoolId   string //学校id
	CampusId   string //校区id
	CampusName string //校区名称
	TokenUrl   string //token获取地址
}

//消息队列配置
type MessageQueue struct {
	Open          string //消息队列服务是否开启，on为开启，off为关闭
	GetMessageUrl string //获取未读消息地址
	Topic         string //要获取的消息
}

func (m *MyConfig) ReadConfig() {
	myConfig := new(Config)
	myConfig.InitConfig("config.txt")
	//海康基础部分
	hkConfig := new(Hk)
	hkConfig.AppKey = myConfig.Read("Hk", "AppKey")
	hkConfig.AppCode = myConfig.Read("Hk", "AppCode")

	//智慧校园基础部分
	zhxyConfig := new(Zhxy)
	zhxyConfig.Domain = myConfig.Read("Zhxy", "Domain")
	zhxyConfig.AppKey = myConfig.Read("Zhxy", "AppKey")
	zhxyConfig.AppSecret = myConfig.Read("Zhxy", "AppSecret")
	campus := strings.Split(myConfig.Read("Zhxy", "Campus"), "|")
	zhxyConfig.SchoolId = campus[0]
	zhxyConfig.CampusId = campus[1]
	zhxyConfig.CampusName = campus[2]
	zhxyConfig.TokenUrl = zhxyConfig.Domain + "/api/Cert/getToken"

	//海康96配置信息
	hk96 := new(Hk96)
	hk96.Domain = myConfig.Read("Hk96", "domain")
	hk96.StudentAsynOpen = myConfig.Read("Hk96", "studentAsynOpen")
	hk96.TeacherAsynOpen = myConfig.Read("Hk96", "teacherAsynOpen")
	hk96.LeaveInOutOpen = myConfig.Read("Hk96", "leaveInOutOpen")
	hk96.InOpen = myConfig.Read("Hk96", "InOpen")
	hk96.OutOpen = myConfig.Read("Hk96", "OutOpen")
	hk96.UploadFaceOpen = myConfig.Read("Hk96", "UploadFaceOpen")
	hk96.InWeekDay = myConfig.Read("Hk96", "InWeekDay")
	hk96.OutWeekDay = myConfig.Read("Hk96", "OutWeekDay")
	hk96.InDevice = myConfig.Read("Hk96", "InDevice")
	hk96.OutDevice = myConfig.Read("Hk96", "OutDevice")
	hk96.InAddTime = myConfig.Read("Hk96", "InAddTime")
	hk96.InDeleteTime = myConfig.Read("Hk96", "InDeleteTime")
	hk96.OutAddTime = myConfig.Read("Hk96", "OutAddTime")
	hk96.OutDeleteTime = myConfig.Read("Hk96", "OutDeleteTime")
	hk96.TestDormDevice = myConfig.Read("Hk96", "TestDormDevice")
	hk96.GetStudentUrl = zhxyConfig.Domain + "/api/Student/ListsBaseInfo"                        //学生基础资料获取
	hk96.GetTeacherUrl = zhxyConfig.Domain + "/api/Teacher/ListsBaseInfo"                        //教职工基础资料获取
	hk96.GetThoroughStudentUrl = zhxyConfig.Domain + "/api/AttendStudent/getThoroughStudentList" //通校学生
	hk96.GetLeaveInStudenturl = zhxyConfig.Domain + "/api/AttendStudent/getRecordStudentList"
	hk96.GetLeaveStudenturl = zhxyConfig.Domain + "/api/StudentLeave/lists"                 //请假学生
	hk96.GetAttendanceFaceList = zhxyConfig.Domain + "/api/HikAttendance/getUndealFaceList" //未上传照片刷脸信息
	hk96.UploadFaceImageUrl = zhxyConfig.Domain + "/api/HikAttendance/uploadFaceImage"      //上传刷脸照片

	hk96.UserAddUrl = hk96.Domain + "/eop/services/common/post/addPerson"
	hk96.UserDeleteUrl = hk96.Domain + "/eop/services/common/post/deletePerson"
	hk96.UserUpdateUrl = hk96.Domain + "/eop/services/common/post/updatePerson"
	hk96.UserUploadPictrueUrl = hk96.Domain + "/eop/services/common/post/uploadImage"
	hk96.TargetAddUrl = hk96.Domain + "/eop/services/common/post/addLeaveTarget"
	hk96.TargetDeleteUrl = hk96.Domain + "/eop/services/common/post/deleteLeaveTarget"

	//海康cip配置信息
	hkCip := new(HkCip)
	hkCip.Domain = myConfig.Read("HkCip", "domain")
	hkCip.StudentAsynOpen = myConfig.Read("HkCip", "studentAsynOpen")
	hkCip.TeacherAsynOpen = myConfig.Read("HkCip", "teacherAsynOpen")
	hkCip.BehaveaOpen = myConfig.Read("HkCip", "behaveaOpen")
	hkCip.ExpressionOpen = myConfig.Read("HkCip", "expressionOpen")
	hkCip.ClassRoomAttendanceOpen = myConfig.Read("HkCip", "classRoomAttendanceOpen")
	hkCip.TestDormDevice = myConfig.Read("hkCip", "TestDormDevice")
	hkCip.UploadFileUrl = hkCip.Domain + "/eop/services/common/post/uploadFile"
	hkCip.StudentAddOrUpdateUrl = hkCip.Domain + "/eop/services/common/post/addOrUpdateStudent"
	hkCip.StudentDeleteUrl = hkCip.Domain + "/eop/services/common/post/deleteStudent"
	hkCip.TeacherAddOrUpdateUrl = hkCip.Domain + "/eop/services/common/post/addOrUpdateTeacher"
	hkCip.TeacherDeleteUrl = hkCip.Domain + "/eop/services/common/post/deleteTeacherUrl"
	hkCip.OrgTreeUrl = hkCip.Domain + "/eop/services/common/get/orgTreeEx"
	hkCip.GradeAddUrl = hkCip.Domain + "/eop/services/common/post/addGradeOrg"
	hkCip.ClassAddUrl = hkCip.Domain + "/eop/services/common/post/addOrUpdateClass"

	hkCip.GetStudentBehaveDataUrl = hkCip.Domain + "/eop/services/common/get/studentOriginalBehaveData"
	hkCip.UploadStudentBehaveDataUrl = zhxyConfig.Domain + "/api/HikStudentBehave/add"

	hkCip.GetStudentExpressionFaceDataUrl = hkCip.Domain + "/eop/services/common/get/studentOriginalFace"
	hkCip.UploadStudentExpressionFaceDataUrl = zhxyConfig.Domain + "/api/HikStudentExpression/add"

	hkCip.GetStudentClassRoomAttendanceUrl = hkCip.Domain + "/eop/services/common/get/studentAtteDetailEx"
	hkCip.UploadStudentClassRoomAttendanceDataUrl = zhxyConfig.Domain + "/Api/AttendanceClassRoom/add"

	//智慧校园消息队列部分
	mqConfig := new(MessageQueue)
	mqConfig.Open = myConfig.Read("MessageQueue", "Open")
	mqConfig.Topic = myConfig.Read("MessageQueue", "Topic")
	mqConfig.GetMessageUrl = zhxyConfig.Domain + "/api/AppMessageQueue/lists"

	m.Hk = hkConfig
	m.HkCip = hkCip
	m.Hk96 = hk96

	m.Zhxy = zhxyConfig
	m.MessageQueue = mqConfig

}

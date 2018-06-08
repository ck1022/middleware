package main

import (
	"config"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"log"
	//"model"
	"encoding/json"
	"library"
	"net/http"
	"os"
	"strings"
	"time"
)

type TeacherInfo struct {
	Jzgid string
	Uid   string
	Gh    string
	Xm    string
	Xbm   string
	Xb    string
	Sfzjh string
	Lxdh  string
}
type TeacherListReturn struct {
	Code int
	Data struct {
		List []TeacherInfo
	}
}

var complete chan int = make(chan int)

func main() {
	myConfig := new(config.Config)
	myConfig.InitConfig("config.txt")
	var isdebug = true
	var tokenUrl = myConfig.Read("ZkAsynBaseData", "tokenUrl")
	var teacherInfoUrl = myConfig.Read("ZkAsynBaseData", "teacherInfoUrl")
	var appkey = myConfig.Read("ZkAsynBaseData", "appkey")
	var appsecret = myConfig.Read("ZkAsynBaseData", "appsecret")
	var server = myConfig.Read("localdb", "server")
	var port = myConfig.Read("localdb", "port")
	var user = myConfig.Read("localdb", "user")
	var password = myConfig.Read("localdb", "password")
	var database = myConfig.Read("localdb", "database")
	var logFileName = "log/ZkAsynBaseData"
	var tmpSchool = ""
	var school = ""
	var tmpCampus = ""
	var campus = ""
	var campusName = ""
	fmt.Println("请输入学校编号:")
	fmt.Scanln(&tmpSchool)
	fmt.Println("再次输入学校编号:")
	fmt.Scanln(&school)
	if tmpSchool != school {
		log.Println("两次输入的学校编号不一致,5秒后退出")
		time.Sleep(5 * time.Second)
		return
	}
	fmt.Println("请输入校区编号:")
	fmt.Scanln(&tmpCampus)
	fmt.Println("再次输入校区编号:")
	fmt.Scanln(&campus)
	if tmpCampus != campus {
		log.Println("两次输入的学校编号不一致，5秒后退出")
		time.Sleep(5 * time.Second)
		return
	}
	fmt.Println("请输入学校名称:")
	fmt.Scanln(&campusName)
	log.Println(school, campus, campusName)
	//日志目录是否存在
	if library.IsDirExists("log") == false {
		os.Mkdir("log", os.ModePerm)
	}
	//连接数据库
	connString := fmt.Sprintf("server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=disable", server, port, database, user, password)
	if isdebug {
		log.Println(connString)
	}
	conn, err := sql.Open("mssql", connString)
	defer conn.Close()
	if err != nil { //打开数据库出错
		library.Panicln(logFileName, "打开数据库失败:", err.Error())
		time.Sleep(30 * time.Second)
		return
	}
	err = conn.Ping()
	if err != nil { //连接数据库出错
		library.Panicln(logFileName, "PING数据库 错误:", err.Error())
		time.Sleep(30 * time.Second)
		return
	}
	//校区是否存在，不存在则添加一个
	var count int32
	err = conn.QueryRow(fmt.Sprintf("select count(*) from [personnel_area] where areaid='%s'", campus)).Scan(&count)
	if err != nil {
		library.Panicln(logFileName, "sql 查询错误:", err.Error())
	}
	//校区不存在，添加一个区域（学校校区）和部门（校区）
	if count == 0 {
		library.Println(logFileName, "校区不存在，添加校区：", campusName)
		_, err := conn.Exec("insert into [System_Area] (AreaCode,AreaName,Status,OperationTime,Flag) values(?,?,1,?,0)", campus, campusName, time.Now().Local().Format("2006-01-02 15:04:05"))
		if err != nil {
			library.Println(logFileName, "添加校区失败(System_Area)", err.Error())
			return
		}
		_, err1 := conn.Exec("insert into [System_Department] (DepartmentCode,DepartmentName,FatherDepartmentCode,Status,OperationTime,Flag) values(?,?,0,1,?,0)", campus, campusName, time.Now().Local().Format("2006-01-02 15:04:05"))
		if err1 != nil {
			library.Println(logFileName, "添加校区部门失败(System_Department)", err.Error())
			return
		}
	}
	//获取token
	timestamp := time.Now().Unix()
	tokenParam := fmt.Sprintf("appkey=%s&timestamp=%d&sign=%s", appkey, timestamp, library.Mymd5(fmt.Sprintf("%s%d%s", appkey, timestamp, appsecret)))
	token, _, code := library.GetToken(tokenUrl, tokenParam)
	if code == -1 {
		library.Panicln(logFileName, "获取token失败")
	}
	//获取服务器的教职工信息，只更新7天内的
	addTime := time.Unix((time.Now().Unix() - 3000*24*60*60), 0)
	param := fmt.Sprintf("appkey=%s&xxid=%s&xqid=%s&timestamp=%d&token=%s&sfsc=2&tjsj=%s", appkey, school, campus, timestamp, token, addTime.Format("2006-01-02"))
	log.Println(teacherInfoUrl, "\n", param)

	teacherList, code := getTeacherList(teacherInfoUrl, param)
	if code == -1 {
		library.Println(logFileName, campusName, "读取教职工列表失败")
	} else {
		for _, ar := range teacherList {
			library.Println(logFileName, "教职工信息："+ar.Jzgid, "\t", ar.Uid, "\t", ar.Xm)
			//教职工是否已经存在
			var tcount int
			conn.QueryRow(fmt.Sprintf("select count(*) from [userinfo] where Cuser1='%s'", ar.Uid)).Scan(&tcount)
			if tcount > 0 {
				library.Println(logFileName, "教职工已经存在")
				continue
			}

			_, err := conn.Exec("insert into [System_Users] (UserCode,UserName,DepartmentCode,UserStatus,OperationTime,Flag) values(?,?,?,1,?,0)", ar.Uid, ar.Xm, campus, time.Now().Local().Format("2006-01-02 15:04:05"))
			if err != nil {
				library.Println(logFileName, "添加教职工失败System_Users)", err.Error())
				continue
			}
			_, err1 := conn.Exec("insert into [HR_WorkCheckUserArea] (UserCode,AreaCode,Status,OperationTime,Flag) values(?,?,1,?,0)", ar.Uid, campus, time.Now().Local().Format("2006-01-02 15:04:05"))
			if err1 != nil {
				library.Println(logFileName, "添加教职工失败(HR_WorkCheckUserArea)", err.Error())
				continue
			}
			library.Println(logFileName, "添加教职工成功")

		}
	}
	time.Sleep(30 * time.Second)
}

//获取教职工列表
func getTeacherList(url string, param string) ([]TeacherInfo, int) {
	resp, _ := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(param))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var info = new(TeacherListReturn)
	if err != nil {
		return nil, -1
	} else {
		json.Unmarshal(body, &info)
		code := info.Code
		list := info.Data.List
		return list, code
	}
}

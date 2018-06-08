/***
*中控人脸识别结果上传到服务器
***/
package main

import (
	"config"
	"database/sql"
	//"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"library"
	"log"
	"model"
	"net/http"
	"os"
	"strings"
	"time"
)

var complete chan int = make(chan int)

func main() {
	//上传人脸识别数据线程
	go upload()
	<-complete
}

//上传人脸识别结果数据到智慧校园平台
func upload() {
	myConfig := new(config.Config)
	myConfig.InitConfig("config.txt")
	var isdebug = true
	var apiUrl = myConfig.Read("ZkAttendanceUpload", "attendanceAddUrl")
	var server = myConfig.Read("localdb", "server")
	var port = myConfig.Read("localdb", "port")
	var user = myConfig.Read("localdb", "user")
	var password = myConfig.Read("localdb", "password")
	var database = myConfig.Read("localdb", "database")
	var logFileName = "log/ZkAttendanceUpload"
	//日志目录是否存在
	if library.IsDirExists("log") == false {
		os.Mkdir("log", os.ModePerm)
	}
	library.Println(logFileName, "程序启动（中控人脸识别结果上传到服务器）...")
	//连接字符串
	connString := fmt.Sprintf("server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=disable", server, port, database, user, password)
	if isdebug {
		log.Println(connString)
	}

	conn, err := sql.Open("mssql", connString)
	if err != nil { //打开数据库出错
		library.Panicln(logFileName, "打开数据库失败:", err.Error())
		complete <- 0
	}
	defer conn.Close()
	for {

		err = conn.Ping()
		if err != nil { //检测连接数据库出错
			library.Panicln(logFileName, "PING 错误:", err.Error())
			complete <- 0
		}

		//产生查询语句的Statement
		currentTime := time.Now().Local()

		rows, err := conn.Query(fmt.Sprintf("select top 50 a.ID,a.UserCode,CheckTime,SN,a.Flag,name from [HR_WorkCheckRecord] as a inner join [userinfo] as b on a.UserCode=b.Cuser1 where (a.Flag=0 or a.Flag=2) and a.CheckTime>'%s' order by ID asc", currentTime.Format("2006-01-02")))
		if err != nil {
			library.Panicln(logFileName, "sql 查询错误:", err.Error())
			complete <- 0
		}
		//建立一个列数组
		var rowsData []*model.Attendance
		//遍历每一行
		for rows.Next() {
			var row = new(model.Attendance)
			rows.Scan(&row.ID, &row.UserCode, &row.CheckTime, &row.SN, &row.Flag, &row.UserName)
			rowsData = append(rowsData, row)
		}
		if len(rowsData) == 0 {

			time.Sleep(3 * time.Second)
			continue
		}
		//上传数据
		library.Println(logFileName, "总共", len(rowsData), "条记录")
		for _, ar := range rowsData {
			library.Println(logFileName, ar.ID, ",", ar.UserCode, ",", ar.UserName, ",", ar.CheckTime.Format("2006-01-02 15:04:05"), ",", ar.SN, ",", ar.Flag)
			param := fmt.Sprintf("sbbh=%s&userCode=%s&userName=%s&checkTime=%s", ar.SN, ar.UserCode, ar.UserName, ar.CheckTime)
			resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(param))
			//上传出错，第一次状态更新为2，第二次更新为1
			if err != nil {
				library.Println(logFileName, "上传失败", err.Error())
				flag := 2
				if ar.Flag == 2 {
					flag = 1
				}
				conn.Query("update [HR_WorkCheckRecord] set Flag=? where ID=?", flag, ar.ID)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				library.Println(logFileName, "上传数据返回读取失败", apiUrl, err.Error())
			} else {
				library.Println(logFileName, "上传数据返回:", string(body))
			}
			conn.Query("update [HR_WorkCheckRecord] set Flag=? where ID=?", 1, ar.ID)
			resp.Body.Close()
		}
		rows.Close()
		time.Sleep(3 * time.Second)
	}
}

/*
//打印线程日志
func Println(logFileName string, v ...interface{}) {

	logFileName = fmt.Sprintf("%s_%s.log", logFileName, time.Now().Local().Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {

	}
	defer logFile.Close()
	infoDebug := log.New(logFile, "[info]", log.Ldate|log.Ltime)
	infoDebug.Println(v)
	log.Println(v)
}

//发生严重错误，要结束线程
func Panicln(logFileName string, v ...interface{}) {
	logFileName = fmt.Sprintf("%s_%s.log", logFileName, time.Now().Local().Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {

	}
	defer logFile.Close()
	errorDebug := log.New(logFile, "[error]", log.Ldate|log.Ltime)
	errorDebug.Println(v)
	log.Println("10秒后关闭线程")
	time.Sleep(30 * time.Second)
	complete <- 0
	log.Panicln(v)
}
*/

package main

/**
*数据同步进程，
*	task1：进出校重点名单设置
*	task2：学生、教职工信息同步
*	task3：人脸比对照片上传
*	task4：课堂行为数据上传
**/
import (
	"library"
	//"log"
	//"os"
	"config"
	"library/myhttp"
	"task"
	"time"
)

var complete chan int = make(chan int)

func main() {
	//日志目录是否存在
	library.Mkdir("log")
	library.Mkdir("picture")
	minite()
}

//每分钟检查执行一遍的任务
func minite() {
	logFileName := "log/cip_minite"

	//学生同步任务
	var m2 chan int = make(chan int, 1)
	task2Runing := false
	task2 := new(task.StudentCip)
	task2.Init()

	//课堂考勤结果
	var m3 chan int = make(chan int, 1)
	task3Runing := false
	task3 := new(task.StudentClassRoomAttendance)
	task3.Init()

	//课堂行为同步任务
	var m4 chan int = make(chan int, 1)
	task4Runing := false
	task4 := new(task.StudentBehavea)
	task4.Init()

	//课堂表情同步任务
	var m5 chan int = make(chan int, 1)
	task5Runing := false
	task5 := new(task.StudentExpressionFace)
	task5.Init()

	i := 0
	for {
		//读取配置文件
		c := new(config.MyConfig)
		c.ReadConfig()
		zhxyHttp := myhttp.Zhxyhttp{Myconfig: c, LogFile: logFileName}
		zhxyHttp.TaskAlive("datacip")

		i++

		//学生同步
		if task2Runing == false {
			task2.Config = c
			go task2.Start(m2, &task2Runing)
			<-m2
		}

		//课堂考勤结果
		if task3Runing == false {
			task3.Config = c
			go task3.Start(m3, &task3Runing)
			<-m3
		}

		//课堂行为结果
		if task4Runing == false {
			task4.Config = c
			go task4.Start(m4, &task4Runing)
			<-m4
		}

		//课堂表情结果
		if task5Runing == false {
			task5.Config = c
			go task5.Start(m5, &task5Runing)
			<-m5
		}

		library.Println(logFileName, "main:", i)
		time.Sleep(60 * time.Second)
	}
}

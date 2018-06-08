package main

/**
*数据同步进程，
*	task1：进出校重点名单设置
*	task2：学生、教职工信息同步
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
	logFileName := "log/96_minite"
	//学生通校和走读生任务
	var m1 chan int = make(chan int, 1)
	task1Runing := false
	task1 := new(task.LeaveInSchool)
	task1.Init()

	//学生、教职工同步任务
	var m2 chan int = make(chan int, 1)
	task2Runing := false
	task2 := new(task.UserAsyn)
	task2.Init()

	//刷脸照片同步任务
	var m3 chan int = make(chan int, 1)
	task3Runing := false
	task3 := new(task.Face)
	task3.Init()

	i := 0
	for {
		//读取配置文件
		c := new(config.MyConfig)
		c.ReadConfig()

		zhxyHttp := myhttp.Zhxyhttp{Myconfig: c, LogFile: logFileName}
		zhxyHttp.TaskAlive("data96")

		i++
		//学生进出校名单
		if task1Runing == false {
			task1.Config = c
			go task1.Start(m1, &task1Runing)
			<-m1
		}

		//学生，教职工同步
		if task2Runing == false {
			task2.Config = c
			go task2.Start(m2, &task2Runing)
			<-m2
		}

		//人脸比对照片上传
		if task3Runing == false {
			task3.Config = c
			go task3.Start(m3, &task3Runing)
			<-m3
		}

		library.Println(logFileName, "main:", i)
		time.Sleep(60 * time.Second)
	}
}

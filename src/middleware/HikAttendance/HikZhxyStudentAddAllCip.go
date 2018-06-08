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
	"task"
)

var complete chan int = make(chan int)

func main() {
	//日志目录是否存在
	library.Mkdir("log")
	library.Mkdir("picture")
	task := new(task.StudentCip)
	c := new(config.MyConfig)
	c.ReadConfig()
	task.Config = c
	task.AddAll()
}

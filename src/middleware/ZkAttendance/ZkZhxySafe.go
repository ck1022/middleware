package main

/**
**与海康同步程序守护进程
 */
import (
	"library"
	"os"
	"os/exec"
	//"syscall"
	"time"
	//"log"
	"regexp"
	"runtime"
	"strings"
)

var logFileName string

func main() {
	logFileName = "log/ZkZhxySafe.log"
	if library.IsDirExists("log") == false {
		os.Mkdir("log", os.ModePerm)
	}

	for {
		checkThread()
		time.Sleep(1800 * time.Second)
	}
}

func checkThread() {
	osString := runtime.GOOS
	if strings.ToLower(osString) == "windows" {
		//检测数据同步程序
		windowCheck("ZkAttendanceUpload")
	} else {
		//检测数据同步程序
		linuxCheck("ZkAttendanceUpload")
	}
}
func linuxCheck(threadName string) {
	cmd := exec.Command("ps", "aux")
	out, err := cmd.CombinedOutput()
	if err != nil {
		library.Panicln(logFileName, "查询进程命令执行失败，"+err.Error())
	}
	processList := string(out)
	r1, _ := regexp.Compile(threadName)
	if r1.MatchString(processList) == false {
		library.Println(logFileName, threadName+"程序没有运行")
		cmd1 := exec.Command("nohup", "./"+threadName, "&")
		err1 := cmd1.Start()
		if err1 != nil {
			library.Panicln(logFileName, threadName+"程序启动失败;"+err1.Error())
		} else {
			library.Println(logFileName, threadName+"程序启动成功")
		}
	} else {
		library.Println(logFileName, threadName+"程序正常运行")
	}
}
func windowCheck(threadName string) {
	cmd := exec.Command("tasklist")
	out, err := cmd.Output()
	if err != nil {
		library.Panicln(logFileName, "查询进程命令执行失败，"+err.Error())
	}
	processList := string(out)
	r1, _ := regexp.Compile(threadName)
	if r1.MatchString(processList) == false {
		library.Println(logFileName, threadName+"程序没有运行")
		cmd1 := exec.Command(threadName + ".exe")
		err1 := cmd1.Start()
		if err1 != nil {
			library.Panicln(logFileName, threadName+"程序启动失败;"+err1.Error())
		} else {
			library.Println(logFileName, threadName+"程序启动成功")
		}
	} else {
		library.Println(logFileName, threadName+"程序正常运行")
	}
}

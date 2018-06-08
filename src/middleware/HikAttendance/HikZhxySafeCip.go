package main

/**
**与海康同步程序守护进程
 */
import (
	"library"
	//"os"
	"os/exec"
	//"syscall"
	"time"
	//"log"
	"config"
	"library/myhttp"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var logFileName string
var exedir string

type ulong int32
type ulong_ptr uintptr

type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte
}

func main() {
	logFileName = "log/cip_Safe"

	for {

		c := new(config.MyConfig)
		c.ReadConfig()
		zhxyHttp := myhttp.Zhxyhttp{Myconfig: c, LogFile: logFileName}
		zhxyHttp.TaskAlive("safecip")
		checkThread()
		time.Sleep(1800 * time.Second)
	}
}

func checkThread() {
	osString := runtime.GOOS
	if strings.ToLower(osString) == "windows" {
		//检测数据同步程序
		//windowCheck("HikZhxyMessageQueueCip")
		//检测消息队列处理程序
		windowCheck("HikZhxyDataAsynCip")
	} else {
		//检测数据同步程序
		//linuxCheck("HikZhxyMessageQueue96")
		//检测消息队列处理程序
		linuxCheck("HikZhxyDataAsyn96")
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
	processString := getProcessString()
	library.Println(logFileName, processString)
	r1, _ := regexp.Compile(threadName)
	if r1.MatchString(processString) == false {
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

func getProcessString() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		return ""
	}
	Process32Next := kernel32.NewProc("Process32Next")
	var processList = []string{}
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			processList = append(processList, strings.Trim(string(proc.szExeFile[0:]), ""))
		} else {
			break
		}
	}
	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, _ = CloseHandle.Call(pHandle)
	processString := strings.Join(processList, " ")
	return processString
}

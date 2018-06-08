package library

import (
	"fmt"
	"log"
	"os"
	"time"
)

//打印线程日志
func Println(logFileName string, v ...interface{}) {
	Writeln(logFileName, "[info]", v)
}

//发生严重错误，要结束线程
func Panicln(logFileName string, v ...interface{}) {
	Writeln(logFileName, "[error]", v)
}
func Writeln(logFileName string, level string, v ...interface{}) {
	logFileName = fmt.Sprintf("%s_%s.log", logFileName, time.Now().Local().Format("02"))
	b, size := getFileByteSize(logFileName)
	if b && size > 100*1000*1000 {
		emptiedFile(logFileName)
	}
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println(err.Error())
	}
	defer logFile.Close()
	errorDebug := log.New(logFile, level, log.Ldate|log.Ltime)
	errorDebug.Println(v)
	log.Println(level, v)
}

func getFileByteSize(filename string) (bool, int64) {
	if !isFileIsExist(filename) {
		return false, 0
	}
	fhandler, _ := os.Stat(filename)
	return true, fhandler.Size()
}
func emptiedFile(filename string) bool {
	FN, err := os.Create(filename)
	defer FN.Close()
	if err != nil {
		return false
	}
	fmt.Fprint(FN, "")
	return true
}

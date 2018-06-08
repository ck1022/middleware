package main

import (
	//"github.com/resize"
	//"library"
	//"task"
	//"io/ioutil"
	//"os"
	//"strconv"
	//"time"
	//"image/jpeg"
	"log"
	//"strings"
	"config"
	"fmt"
	"library/myhttp"
)

func main() {
	c := new(config.MyConfig)
	c.ReadConfig()
	/*
		hikhttp := myhttp.Hikhttp{Myconfig: c}
		//hikhttp := new(myhttp.Hikhttp)
		//hikhttp.Myconfig = c
		//var learnCodes = []string{"111", "222"}
		result, flag := hikhttp.DeleteTarget("001108", "111")
		log.Println(flag)
		log.Println(result)
	*/
	zhxyhttp := myhttp.Zhxyhttp{Myconfig: c}
	//hikHttp := myhttp.Hikhttp{Myconfig: c}
	//上传到美智
	param := make(map[string]string)
	param["cameraIndex"] = "001111"
	param["recordTime"] = "2018-05-22 15:00:00"
	param["behaviour"] = "reg"
	param["x1"] = fmt.Sprintf("%d", 1)
	param["y1"] = fmt.Sprintf("%d", 1)
	param["x2"] = fmt.Sprintf("%d", 1)
	param["y2"] = fmt.Sprintf("%d", 1)
	pic := "StudentExpressionFace.jpg"
	//return
	s, _ := zhxyhttp.UploadWithFile(c.HkCip.UploadStudentBehaveDataUrl, pic, param)
	log.Println("返回", string(s))
	//hikHttp.AddTarget("001155", "20151328,20171141,20171134,20171285,20151324")
	//returnData, _ := hikHttp.UpPicNew("picture/tmp.jpg")
	//log.Println(returnData)
	//addReturnData, _ := hikHttp.AddOrUpdateStudent("0001", "测试111", "1", "000000", "")
	//log.Println(addReturnData)
	/**获取学生行为**/
	//res, _ := hikHttp.GetStudentOriginalBehaveData("2018-05-09 08:30:00", "2018-05-09 18:25:00", 1, 100)
	//log.Println("行为调用返回", string(res))
	/**获取学生表情**/
	//res1, _ := hikHttp.GetStudentOriginalFaceData("2018-05-03 11:00:00", "2018-05-09 19:00:00", 1, 100)
	//log.Println("表情调用返回", string(res1))
	/*
		zhxyhttp.GetToken()
		var map1 = make(map[string]string)
		//map1["xsid"] = "111"
		map1["date"] = "2018-01-02"
		map1["in"] = "1"
		//b, _ := zhxyhttp.GetLeaveStudentList("2018-01-02")
		//log.Println(b)
	*/
}

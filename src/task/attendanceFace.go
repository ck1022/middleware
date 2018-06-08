package task

import (
	"config"
	//"fmt"
	"library"
	"library/myhttp"
	"regexp"
)

type Face struct {
	Config      *config.MyConfig
	logFileName string
}

func (t *Face) Init() {
	t.logFileName = "log/96_attandanceFace"
}

func (t *Face) Start(complete chan<- int, runing *bool) {
	*runing = true
	complete <- 0
	/**是否开启同步**/
	if t.Config.Hk96.UploadFaceOpen != "on" {
		*runing = false
		return
	}
	/**获取未同步照片数据列表**/
	zhxyHttp := myhttp.Zhxyhttp{Myconfig: t.Config, LogFile: t.logFileName}
	list, _ := zhxyHttp.GetNullAttendanceFaceList()
	if len(list) > 0 {
		library.Println(t.logFileName, "共(", len(list), ")条")
	} else {
		*runing = false
		return
	}
	for _, v := range list {
		library.Println(t.logFileName, "["+v.Ysjlid+"]开始处理")
		r1, _ := regexp.Compile("http://")
		e := true
		if r1.MatchString(v.FaceUrl) == false { //没有http
			_, e = library.GetImageNew(t.Config.Hk96.Domain+":8080/kms/services/rest/dataInfoService/downloadFile?id="+v.FaceUrl, "face.jpg")
		} else {
			_, e = library.GetImageNew(v.FaceUrl, "face.jpg")
		}
		//_, e := library.GetImageNew(t.Config.Hk96.Domain+":8080/kms/services/rest/dataInfoService/downloadFile?id="+v.FaceUrl, "face.jpg")
		//_, e := library.GetImageNew(v.FaceUrl, "face.jpg")
		if e {
			library.Println(t.logFileName, "["+v.Ysjlid+"]从海康下载图片成功")
		} else {
			library.Panicln(t.logFileName, "["+v.Ysjlid+"]从海康下载图片失败")
			continue
		}
		s, e := zhxyHttp.UpPic("face.jpg", v.Ysjlid)
		if e {
			library.Println(t.logFileName, "["+v.Ysjlid+"]上传到服务器成功")
		} else {
			library.Panicln(t.logFileName, "["+v.Ysjlid+"]上传到服务器失败", string(s))
		}
	}
	*runing = false
}

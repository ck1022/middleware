package model

//学生列表返回
type StudentReturn struct {
	Code int
	Data struct {
		List []Student
	}
}

//教职工列表返回
type TeacherReturn struct {
	Code int
	Data struct {
		List []Teacher
	}
}

//消息队列返回
type MessageReturn struct {
	Code int
	Data struct {
		List []Message
	}
}

//
type NullAttendanceFaceReturn struct {
	Code int
	Data struct {
		List []Face
	}
}

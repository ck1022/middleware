package model

import (
	"time"
)

type Attendance struct {
	ID             int32
	UserCode       string
	CheckTime      time.Time
	SN             string
	Flag           int32
	sn_name        string
	OperationTime  time.Time
	processingTime time.Time
	Temark         string
	UserName       string
}

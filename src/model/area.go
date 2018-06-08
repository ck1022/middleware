package model

import (
	"time"
)

type Area struct {
	ID             int32
	AreaCode       string
	FatherAreaCode string
	Status         string
	OperationTime  time.Time
	Flag           int32
}

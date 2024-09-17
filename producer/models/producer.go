package models

import (
	"strings"
	"time"
)

type LogMessage struct {
	AppName     string     `json:"app_name"`
	ServiceName string     `json:"service_name"`
	LogLevel    string     `json:"log_level"`
	Message     string     `json:"message"`
	Timestamp   CustomTime `json:"timestamp"`
}

type CustomTime struct {
	time.Time
}

const customTimeLayout = "2006-01-02T15:04:05.999999"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	ct.Time, err = time.Parse(customTimeLayout, s)
	return
}

package utils

import (
	"fmt"
	"os"
	"time"
)

func getHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost:2566"
	}
	return host
}

func Href(key, value string) string {
	return fmt.Sprintf("%s/%s/%s", getHost(), key, value)
}

func ConvertTimeBangkok(dataTime string) string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	t, _ := time.Parse("2006-01-02T15:04:05Z07:00", dataTime)
	return t.In(loc).Format("2006-01-02T15:04:05Z07:00")
}

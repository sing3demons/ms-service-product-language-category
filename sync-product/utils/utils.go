package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func GetHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost:2566"
	}
	return host
}

func Href(key, value string) string {
	return fmt.Sprintf("%s/%s/%s", GetHost(), strings.ToLower(key), value)
}

func ConvertTimeBangkok(dataTime string) string {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("error load location")
		fmt.Println(err)
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", dataTime)
	if err != nil {
		fmt.Println("error parse time")
		fmt.Println(err)
	}
	return t.In(loc).Format("2006-01-02T15:04:05Z07:00")
}

func GetTransactionID() string {
	return time.Now().Format("20060102150405")
}

package utils

import (
	"log"
	"time"
)

//时间格式（2006-01-02 15:04:05）的字符串转换为时间戳
func StringToTimestamp(s string) int64 {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Println(err)
		return 0
	}
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	if err != nil {
		log.Println(err)
		return 0
	}
	return t1.Unix()
}

//时间戳转为时间格式（2006-01-02 15:04:05）的字符串
func TimestampToString(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func GetNewTimestamp() int64 {
	return time.Now().Unix()
}

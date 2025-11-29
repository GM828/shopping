package util

import (
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

const timeLayout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	str := string(data)
	if len(str) > 2 {
		str = str[1 : len(str)-1] // 去除引号
	}
	t, err := time.Parse(timeLayout, str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + ct.Time.Format(timeLayout) + `"`), nil
}

var (
	DateUtil dateUtil
)

type dateUtil struct {
}

func (d *dateUtil) TimePtr(t time.Time) *time.Time {
	return &t
}

// FormatDateByCustomLayout 将time格式化为指定的格式
func (d *dateUtil) FormatDateByCustomLayout(time *time.Time, layout string) string {
	if time == nil {
		return ""
	}
	return time.Format(layout)
}

func (d *dateUtil) ParseStandard(dateStr *string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}
	date, err := d.ParseByLayout(dateStr, DateLayout.YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		return nil, err
	}
	return date, nil
}

func (d *dateUtil) ParseByLayout(dateStr *string, layout string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date, err := time.ParseInLocation(layout, *dateStr, loc)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

// 日期格式常量
var DateLayout = struct {
	YYYY_MM_DD_HH_MM_SS string //yyyy-MM-dd HH:mm:ss
	YYYY_MM_DD          string //yyyy-MM-dd
	YYYYMMDD            string //yyyyMMdd
	YYYYMMDDHHMM        string //yyyyMMddHHmm
	YYYYMMDDHHMMSS      string //yyyyMMddHHmmss
	YYYYMM              string //yyyyMM
	YYYY_MM_DD_HH_MM    string //yyyy-MM-dd HH:mm
	YYYY_MM_DD_CN       string //yyyy-MM-dd Chinese format
	YYYY                string // yyyy
}{
	YYYY_MM_DD_HH_MM_SS: "2006-01-02 15:04:05", //yyyy-MM-dd HH:mm:ss
	YYYY_MM_DD:          "2006-01-02",          //yyyy-MM-dd
	YYYYMMDD:            "20060102",            //yyyyMMdd
	YYYYMMDDHHMM:        "200601021504",        //yyyyMMddHHmm
	YYYYMMDDHHMMSS:      "20060102150405",      //yyyyMMddHHmmss
	YYYYMM:              "200601",              //yyyyMM
	YYYY_MM_DD_HH_MM:    "2006-01-02 15:04",    //yyyy-MM-dd HH:mm
	YYYY_MM_DD_CN:       "2006年01月02日",         //yyyy-MM-dd Chinese format
	YYYY:                "2006",                // yyyy
}

func Int64ToStr(i int64) string {
	return fmt.Sprintf("%d", i)
}

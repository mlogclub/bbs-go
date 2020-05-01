package simple

import (
	"strconv"
	"time"
)

const (
	FMT_DATE_TIME    = "2006-01-02 15:04:05"
	FMT_DATE         = "2006-01-02"
	FMT_TIME         = "15:04:05"
	FMT_DATE_TIME_CN = "2006年01月02日 15时04分05秒"
	FMT_DATE_CN      = "2006年01月02日"
	FMT_TIME_CN      = "15时04分05秒"
)

// 秒时间戳
func NowUnix() int64 {
	return time.Now().Unix()
}

// 毫秒时间戳
func NowTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// 毫秒时间戳
func Timestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// 秒时间戳转时间
func TimeFromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// 毫秒时间戳转时间
func TimeFromTimestamp(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

// 时间格式化
func TimeFormat(time time.Time, layout string) string {
	return time.Format(layout)
}

// 字符串时间转时间类型
func TimeParse(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// return yyyyMMdd
func GetDay(time time.Time) int {
	ret, _ := strconv.Atoi(time.Format("20060102"))
	return ret
}

// 返回指定时间当天的开始时间
func WithTimeAsStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

/**
 * 将时间格式换成 xx秒前，xx分钟前...
 * 规则：
 * 59秒--->刚刚
 * 1-59分钟--->x分钟前（23分钟前）
 * 1-24小时--->x小时前（5小时前）
 * 昨天--->昨天 hh:mm（昨天 16:15）
 * 前天--->前天 hh:mm（前天 16:15）
 * 前天以后--->mm-dd（2月18日）
 */
func PrettyTime(timestamp int64) string {
	_time := TimeFromTimestamp(timestamp)
	_duration := (NowTimestamp() - timestamp) / 1000
	if _duration < 60 {
		return "刚刚"
	} else if _duration < 3600 {
		return strconv.FormatInt(_duration/60, 10) + "分钟前"
	} else if _duration < 86400 {
		return strconv.FormatInt(_duration/3600, 10) + "小时前"
	} else if Timestamp(WithTimeAsStartOfDay(time.Now().Add(-time.Hour*24))) <= timestamp {
		return "昨天 " + TimeFormat(_time, FMT_TIME)
	} else if Timestamp(WithTimeAsStartOfDay(time.Now().Add(-time.Hour*24*2))) <= timestamp {
		return "前天 " + TimeFormat(_time, FMT_TIME)
	} else {
		return TimeFormat(_time, FMT_DATE)
	}
}

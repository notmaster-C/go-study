package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mrand "math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	// MilliSecond 毫秒
	MilliSecond = 1000
	// MinuteMilliSecond 一分钟的毫秒
	MinuteMilliSecond = int64(60000)
	// HourMilliSecond 一小时的毫秒
	HourMilliSecond = 3600000
	// DayMilliSecond 一天的毫秒
	DayMilliSecond = 24 * HourMilliSecond
	// DaySecond 一天的秒
	DaySecond = 24 * 3600
	// WeekMilliSecond 一周的毫秒
	WeekMilliSecond = 7 * DayMilliSecond
	// WeekSecond 一周的秒
	WeekSecond = WeekMilliSecond / 1000
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	startTime   = Now()
)

func ScanRows2map(rows *sql.Rows) (res []map[string]interface{}) {
	defer rows.Close()
	cols, _ := rows.Columns()
	cache := make([]interface{}, len(cols))
	// 为每一列初始化一个指针
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	for rows.Next() {
		rows.Scan(cache...)
		row := make(map[string]interface{})
		for i, val := range cache {
			// 处理数据类型
			v := *val.(*interface{})
			switch v.(type) {
			case []uint8:
				v = string(v.([]uint8))
			case nil:
				v = ""
			}
			row[cols[i]] = v
		}
		res = append(res, row)
	}
	rows.Close()
	return res
}

func ParseNestedJSON(jsonData []byte, keys ...string) (any, error) {
	var result any
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if m, ok := result.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				result = val
			} else {
				return nil, fmt.Errorf("Key '%s' not found", key)
			}
		} else {
			return nil, fmt.Errorf("Value is not a map[string]interface{}")
		}
	}

	return result, nil
}

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[mrand.Intn(len(charset))])
	}
	return sb.String()
}
func If2float64(value interface{}) (float64, error) {
	// 尝试直接转换为float64
	if floatVal, ok := value.(float64); ok {
		return floatVal, nil
	}

	// 如果不是float64，尝试将字符串转换为float64
	if strVal, ok := value.(string); ok {
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return 0, errors.New("无法将字符串转换为float64")
		}
		return floatVal, nil
	}

	// 如果既不是float64也不是字符串，返回错误
	return 0, errors.New("无法将值转换为float64")
}

// ToInt 获取转换后的json字段
func ToInt(mi interface{}) int {
	i := int(0) // your final value

	switch t := mi.(type) {
	case int:
		i = t
	case int8:
		i = int(t) // standardizes across systems
	case int16:
		i = int(t) // standardizes across systems
	case int32:
		i = int(t) // standardizes across systems
	case int64:
		i = int(t) // standardizes across systems
	case float32:
		i = int(t) // standardizes across systems
	case float64:
		i = int(t) // standardizes across systems
	case uint8:
		i = int(t) // standardizes across systems
	case uint16:
		i = int(t) // standardizes across systems
	case uint32:
		i = int(t) // standardizes across systems
	case uint64:
		i = int(t) // standardizes across systems
	default:
		v, _ := strconv.ParseInt(fmt.Sprintf("%v", t), 10, 0)
		i = int(v)
	}

	return i
}

// ToInt64 转换成64位整数
func ToInt64(mi interface{}) int64 {
	i := int64(0) // your final value

	switch t := mi.(type) {
	case int:
		i = int64(t)
	case int8:
		i = int64(t) // standardizes across systems
	case int16:
		i = int64(t) // standardizes across systems
	case int32:
		i = int64(t) // standardizes across systems
	case int64:
		i = int64(t) // standardizes across systems
	case float32:
		i = int64(t) // standardizes across systems
	case float64:
		i = int64(t) // standardizes across systems
	case uint8:
		i = int64(t) // standardizes across systems
	case uint16:
		i = int64(t) // standardizes across systems
	case uint32:
		i = int64(t) // standardizes across systems
	case uint64:
		i = int64(t) // standardizes across systems
	default:
		i, _ = strconv.ParseInt(fmt.Sprintf("%v", t), 10, 0)
	}

	return i
}

// ToFloat64 转换双浮点数
func ToFloat64(mi interface{}) float64 {
	i := float64(0) // your final value

	switch t := mi.(type) {
	case int:
		i = float64(t)
	case int8:
		i = float64(t) // standardizes across systems
	case int16:
		i = float64(t) // standardizes across systems
	case int32:
		i = float64(t) // standardizes across systems
	case int64:
		i = float64(t) // standardizes across systems
	case float32:
		i = float64(t) // standardizes across systems
	case float64:
		i = float64(t) // standardizes across systems
	case uint8:
		i = float64(t) // standardizes across systems
	case uint16:
		i = float64(t) // standardizes across systems
	case uint32:
		i = float64(t) // standardizes across systems
	case uint64:
		i = float64(t) // standardizes across systems
	default:
		const bitSize = 64
		i, _ = strconv.ParseFloat(fmt.Sprintf("%v", t), bitSize)
	}

	return i
}

// MinusInt32 取负数
func MinusInt32(mi interface{}) int32 {
	val := int32(ToInt(mi))
	if val <= 0 {
		return val
	}
	return -val
}

// MinusInt64 取负数
func MinusInt64(mi interface{}) int64 {
	val := ToInt64(mi)
	if val <= 0 {
		return val
	}
	return -val
}

// Max 取最大值
func Max(x, y interface{}) int {
	xVal := ToInt(x)
	yVal := ToInt(y)
	if xVal >= yVal {
		return xVal
	}
	return yVal
}

// Now 系统当前时间（毫秒）
func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// NowSecond 系统当前时间（秒）
func NowSecond() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

// GetDayTime 获取0点时间戳
func GetDayTime() int64 {
	t := time.Now()
	return (t.Unix() - int64(t.Hour()*3600+t.Minute()*60+t.Second())) * 1000
}

// GetDayTimeByTime 获取指定点对应的时间戳
func GetDayTimeByTime(t time.Time) int64 {
	return (t.Unix() - int64(t.Hour()*3600+t.Minute()*60+t.Second())) * 1000
}

// GetWeekTime 获取一周开始时间，逻辑上周一0点
func GetWeekTime() int64 {
	t := time.Now()
	// 逻辑周一，因此需要在周日的基础上再加1天时间
	return (t.Unix()-int64((int(t.Weekday())*24+t.Hour())*3600+t.Minute()*60+t.Second()))*1000 + DayMilliSecond
}

// GetWeeklyDiff 获取时间周间隔
func GetWeeklyDiff(oldTime int64) int {
	t1 := time.Unix(oldTime/1000, 0)
	t2 := time.Now()
	weekFun := func(t time.Time) int {
		weekDay := int(t.Weekday())
		// 星期天是一周最后一天
		if weekDay == 0 {
			weekDay = 7
		}
		return weekDay
	}
	weekDay1 := weekFun(t1)
	weekDay2 := weekFun(t2)

	diffWeek := int(t2.Sub(t1).Hours())/(7*24) + 1
	if weekDay1 > weekDay2 {
		diffWeek++
	}

	return diffWeek
}

// GetMonthTime 获取一月开始时间
func GetMonthTime() int64 {
	t := time.Now()
	return (t.Unix()-int64((int(t.Day())*24+t.Hour())*3600+t.Minute()*60+t.Second()))*1000 + DayMilliSecond
}

// GetNextMonthTime 获取下个月月开始时间
func GetNextMonthTime() int64 {
	t := time.Now().AddDate(0, 1, 0)
	return (t.Unix()-int64((int(t.Day())*24+t.Hour())*3600+t.Minute()*60+t.Second()))*1000 + DayMilliSecond
}

// DateInt 系统当前时间（毫秒）
func DateInt() int {
	t := time.Now()
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

// FormatTime 格式化时间
func FormatTime(t int64) string {
	info := ""
	bHasPrefix := false

	if usedDay := t / DayMilliSecond; usedDay > 0 {
		info += fmt.Sprintf("%dd ", usedDay)
		t -= usedDay * DayMilliSecond
		bHasPrefix = true
	}

	if hour := t / HourMilliSecond; hour > 0 || bHasPrefix {
		info += fmt.Sprintf("%dh ", hour)
		t -= hour * HourMilliSecond
	}
	t = t / MilliSecond
	minute := t / 60
	second := t % 60
	info += fmt.Sprintf("%02d:%02d", minute, second)
	return info
}

// 随机域名字符串
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mrand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandomToken 获取随机字符串序号
func RandomToken(prefixID uint64, length int) string {
	if length == 0 {
		length = 16
	}
	b := make([]byte, length)
	rand.Read(b)
	if prefixID != 0 {
		binary.BigEndian.PutUint64(b, prefixID)
	}
	return hex.EncodeToString(b)
}

// ToString 转换成字符串
func ToString(mi interface{}) string {
	switch mi.(type) {

	case string:
		return mi.(string)
	case int, int32, int64, uint32, uint64:
		return fmt.Sprintf("%d", mi)
	}
	return fmt.Sprintf("%v", mi)
}

// GetIP 获取本地IP地址
func GetIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}

			}
		}
	}
	return "127.0.0.1"
}

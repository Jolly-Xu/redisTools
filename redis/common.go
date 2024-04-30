package redis

import (
	"encoding"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func PrintDebug(redisContext *InterceptorContext) {
	// 设置颜色
	colorReset := "\033[0m"
	colorYellow := "\033[33m"
	colorCyan := "\033[36m"

	// 获取当前时间
	currentTime := time.Now().Format("2006-04-02 15:04:05")

	// 打印时间戳
	fmt.Printf("%s%s%s  ", colorYellow, currentTime, colorReset)

	// 打印Redis命令
	fmt.Printf("%sExecuting Redis Command: %s", colorCyan, colorReset)
	cmd := strings.TrimSuffix(redisContext.Cmd().Cmd(), "\r\n")
	fmt.Printf("%s  ", cmd)

	fmt.Printf("%sResponse: %s", colorCyan, colorReset)
	fmt.Printf("%s\n", redisContext.Cmd().Res())
}

func PrintDebug2(redisContext *InterceptorContext) {
	// 设置颜色
	colorReset := "\033[0m"
	colorYellow := "\033[33m"
	colorCyan := "\033[36m"

	// 获取当前时间
	currentTime := time.Now().Format("2006-04-02 15:04:05")

	// 打印时间戳
	fmt.Printf("%s%s%s \t", colorYellow, currentTime, colorReset)

	// 打印Redis命令
	cmd := strings.TrimSuffix(redisContext.Cmd().Cmd(), "\n")
	fmt.Printf("%sExecuting Redis Command: %s%s\t", colorCyan, cmd, colorReset)

	// 打印响应
	res := redisContext.Cmd().Res()
	fmt.Printf("%sResponse: %s%s\n", colorCyan, res, colorReset)
}

func convertInterfaceToString(v interface{}) string {
	switch v := v.(type) {
	case nil:
		return ""
	case string:
		return v
	case *string:
		return *v
	case []byte:
		return hex.EncodeToString(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case *int:
		return strconv.FormatInt(int64(*v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case *int8:
		return strconv.FormatInt(int64(*v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case *int16:
		return strconv.FormatInt(int64(*v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case *int32:
		return strconv.FormatInt(int64(*v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case *int64:
		return strconv.FormatInt(*v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case *uint:
		return strconv.FormatUint(uint64(*v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case *uint8:
		return strconv.FormatUint(uint64(*v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case *uint16:
		return strconv.FormatUint(uint64(*v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case *uint32:
		return strconv.FormatUint(uint64(*v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case *uint64:
		return strconv.FormatUint(*v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case *float32:
		return strconv.FormatFloat(float64(*v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case *float64:
		return strconv.FormatFloat(*v, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case *bool:
		if *v {
			return "true"
		}
		return "false"
	case time.Time:
		return v.Format(time.RFC3339)
	case time.Duration:
		return v.String()
	case encoding.BinaryMarshaler:
		binary, err := v.MarshalBinary()
		if err != nil {
			fmt.Println("can't marshal %T to string or byte string", v)
			return ""
		}
		return hex.EncodeToString(binary)
	case net.IP:
		return hex.EncodeToString(v)
	default:
		fmt.Println("redis: can't marshal %T (implement encoding.BinaryMarshaler)", v)
		return ""
	}
}

package redis

import (
	"fmt"
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
	fmt.Printf("%s%s%s ", colorYellow, currentTime, colorReset)

	// 打印Redis命令
	fmt.Printf("\t%sExecuting Redis Command: %s", colorCyan, colorReset)
	cmd := strings.TrimSuffix(redisContext.Cmd().Cmd(), "\n")
	fmt.Printf("%s", cmd)

	fmt.Printf("\t%sResponse: %s", colorCyan, colorReset)
	fmt.Printf("%s\n", redisContext.Cmd().Res())
}

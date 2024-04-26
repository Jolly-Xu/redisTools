package conn

import (
	"bufio"
	"fmt"
	"net"
)

// 字符串常量
const (
	AUTH = "auth"
	END  = "\n"
)

// NewConn
// 创建一个新的连接，使用Golang的net包
func NewConn(netType string, addr string) (conn net.Conn) {
	// TODO 可能需要处理一下连接失败的情况，因为连接失败有很多种，可能网络故障，可能地址填错等
	conn, err := net.Dial(netType, addr)
	if err != nil {
		fmt.Println("创建连接失败：", err)
		return nil
	}
	return conn
}

func Auth(password string, conn net.Conn) bool {
	authCommand := AUTH + password + END
	_, err := conn.Write([]byte(authCommand))
	if err != nil {
		fmt.Println("Redis权限认证错误:", err)
		return false
	}
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取Redis返回值出错:", err)
		return false
	}
	// 如果认证成功，可以继续发送其他命令进行操作
	if response != "+OK\r\n" {
		fmt.Println("Redis 认证失败.")
		return false
	}
	return true
}

func CommandNoResult(cmd string, conn net.Conn) (success bool) {
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		fmt.Println("执行命令失败", err)
		return false
	}
	return true
}

func CommandGetResult(cmd string, conn net.Conn) (result []byte, success bool) {
	if CommandNoResult(cmd, conn) {
		reader := bufio.NewReader(conn)
		read, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("")
			return nil, false
		}
		return read, true
	}
	return nil, false
}

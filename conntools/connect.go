package conntools

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

type RedisConn struct {
	// TODO 先做一个简答的包装，后续如果需要再添加，这样可以把这个结构体注入到连接池里面
	conn net.Conn
}

// NewRedisConn 创建一个新的连接，使用Golang的net包
func NewRedisConn(netType string, addr string) *RedisConn {
	conn, err := net.Dial(netType, addr)
	if err != nil {
		fmt.Println("创建连接失败")
		return nil
	}
	return &RedisConn{conn: conn}
}

func (r *RedisConn) Auth(password string) bool {
	authCommand := AUTH + password + END
	_, err := r.conn.Write([]byte(authCommand))
	if err != nil {
		fmt.Println("Redis权限认证错误:", err)
		return false
	}
	reader := bufio.NewReader(r.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("认证时读取Redis返回值出错:", err)
		return false
	}
	// 如果认证成功，可以继续发送其他命令进行操作
	if response != "+OK\r\n" {
		fmt.Println("Redis 认证失败.")
		return false
	}
	return true
}

func (r *RedisConn) CommandNoResult(cmd string) (success bool) {
	_, err := r.conn.Write([]byte(cmd))
	if err != nil {
		fmt.Println("执行命令失败", err)
		return false
	}
	return true
}

func (r *RedisConn) CommandGetResult(cmd string) (result []byte, success bool) {
	if r.CommandNoResult(cmd) {
		reader := bufio.NewReader(r.conn)
		read, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("")
			return nil, false
		}
		return read, true
	}
	return nil, false
}

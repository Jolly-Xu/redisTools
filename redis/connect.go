package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// 字符串常量
const (
	AUTH = "auth"
	END  = "\r\n"
)

type RedisConn struct {
	// TODO 先做一个简答的包装，后续如果需要再添加，这样可以把这个结构体注入到连接池里面
	conn net.Conn
}

// NewRedisConn 创建一个新的连接，使用Golang的net包
func NewRedisConn(netType string, addr string) *RedisConn {
	conn, err := net.Dial(netType, addr)
	if err != nil {
		fmt.Println("创建连接失败", err)
		return nil
	}
	return &RedisConn{conn: conn}
}

func (r *RedisConn) Auth(password string) bool {
	authCommand := AUTH + " " + password + END
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
	_, err := r.conn.Write([]byte(cmd))
	success = true
	if err != nil {
		fmt.Println("执行命令失败", err)
		success = false
		return
	}
	if success {
		reader := bufio.NewReader(r.conn)
		respLenBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("无法读取响应长度:", err)
			return
		}

		// 解析响应长度
		respLenStr := strings.TrimPrefix(string(respLenBytes), "$")
		respLen := 0
		_, err = fmt.Sscanf(respLenStr, "%d", &respLen)
		if err != nil {
			return nil, false
		}
		// 读取响应内容
		for i := 0; i < respLen; i++ {
			responseByte, err := reader.ReadByte()
			if err != nil {
				fmt.Println("无法读取响应内容:", err)
				return
			}
			result = append(result, responseByte)
		}

		// 读取响应结束符
		_, err = reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("无法读取响应结束符:", err)
			return
		}
	}
	return
}

func (r *RedisConn) SendCommand(command string, args ...interface{}) (interface{}, error) {
	commandStr := command
	for _, arg := range args {
		commandStr += fmt.Sprintf(" %v", arg)
	}
	commandStr += End

	_, err := r.conn.Write([]byte(commandStr))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(r.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return response, nil
}

// BeginTransaction 开启事务
func (r *RedisConn) BeginTransaction() error {
	_, err := r.SendCommand("MULTI")
	return err
}

// EndTransaction 结束事务
func (r *RedisConn) EndTransaction() error {
	_, err := r.SendCommand("EXEC")
	return err
}

// Rollback 添加一个 Rollback 方法来执行回滚操作
func (r *RedisConn) Rollback() error {
	_, err := r.SendCommand("DISCARD") // 发送 DISCARD 命令来取消事务
	return err

}

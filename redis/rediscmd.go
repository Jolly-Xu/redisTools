package redis

import (
	"context"
	"time"
)

//const (
//	get = iota
//	set
//	setex
//	sel
//	ping
//	del
//	flushdb
//	flushall
//	keys
//	expire
//)

// RedisCommand 表示一个Redis命令
type RedisCommand struct {
	// 命令名
	Name string
	// 命令参数
	Args []interface{}
	// 预期的回复类型
	ResponseType string
	// 超时时间
	Timeout time.Duration
	// 上下文信息
	Ctx context.Context
	// 错误信息
	Err error
}

// NewRedisCommand 创建一个新的Redis命令实例
func NewRedisCommand(name string, args []interface{}, responseType string, timeout time.Duration) *RedisCommand {
	return &RedisCommand{
		Name:         name,
		Args:         args,
		ResponseType: responseType,
		Timeout:      timeout,
		Ctx:          context.Background(), // 默认情况下使用Background作为上下文
	}
}

// WithContext 为Redis命令设置上下文
func (cmd *RedisCommand) WithContext(ctx context.Context) *RedisCommand {
	cmd.Ctx = ctx
	return cmd
}

// SetError 为Redis命令设置错误信息
func (cmd *RedisCommand) SetError(err error) {
	cmd.Err = err
}

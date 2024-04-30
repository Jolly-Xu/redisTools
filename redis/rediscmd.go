package redis

import (
	"time"
)

const (
	get = iota
	set
	setex
	sel
	ping
	del
	flushdb
	flushall
	keys
	expire
)

const (
	Get = "GET "
	End = "\r\n"
	Set = "SET "
)

// RedisCommand 表示一个Redis命令
type RedisCommand struct {
	// 命令名
	cmd string

	// 参数，依次为key,value……等
	args []interface{}

	// 超时时间
	timeout time.Duration

	// 返回结果
	res []byte

	// 命令种类
	cmdType int
}

func (r *RedisCommand) Res() []byte {
	return r.res
}

func (r *RedisCommand) SetRes(res []byte) {
	r.res = res
}

func (r *RedisCommand) Args() []interface{} {
	return r.args
}

func (r *RedisCommand) SetArgs(args []interface{}) {
	r.args = args
}

func newRedisCommand(cmd string, timeout time.Duration, res []byte, cmdType int) *RedisCommand {
	return &RedisCommand{cmd: cmd, timeout: timeout, res: res, cmdType: cmdType}
}

func newGetCommand(cmd string, timeout time.Duration) *RedisCommand {
	cmd = Get + cmd + End
	return newRedisCommand(cmd, timeout, nil, get)
}

func newSetCommand(cmd string, timeout time.Duration) *RedisCommand {
	cmd = Set + cmd + End
	return newRedisCommand(cmd, timeout, nil, set)
}

func (r *RedisCommand) CmdType() int {
	return r.cmdType
}

func (r *RedisCommand) AddArgs(item interface{}) {
	r.args = append(r.args, item)
}

func (r *RedisCommand) SetCmdType(cmdType int) {
	r.cmdType = cmdType
}

func (r *RedisCommand) Cmd() string {
	return r.cmd
}

func (r *RedisCommand) SetCmd(cmd string) {
	r.cmd = cmd
}

func (r *RedisCommand) Timeout() time.Duration {
	return r.timeout
}

func (r *RedisCommand) SetTimeout(timeout time.Duration) {
	r.timeout = timeout
}

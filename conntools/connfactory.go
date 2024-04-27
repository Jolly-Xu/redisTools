package conn

import (
	"redisTools/redis"
	"strconv"
)

// DefaultFactory 默认创建连接工厂
func DefaultFactory(config *redis.Config) *RedisConn {
	dst := config.IpAddr + ":" + strconv.Itoa(config.Port)
	conn := NewRedisConn(config.NetType, dst)
	return conn
}

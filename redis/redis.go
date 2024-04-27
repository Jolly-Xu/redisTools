package redis

import (
	"redisTools/conntools"
)

// redis客户端主要代码

type Redis struct {
	//  配置类
	Config *Config
	// 连接池
	connPool conntools.ConnPool
}

func CreateRedis(config *Config) (r *Redis) {
	//TODO 进行一些初始化操作
	r = &Redis{
		Config: config,
	}
	conntools.NewConnPool(config, nil, 5, 10, 8)
	return
}

func (r *Redis) GetToString(key string) (res string) {
	conn := r.connPool.Get()
	result, success := conn.CommandGetResult(key)
	if success {
		res = string(result)
	}
	return
}

package redis

// redis客户端主要代码

const (
	get = "GET "
	end = "\n"
)

type Redis struct {
	//  配置类
	Config *Config
	// 连接池
	connPool *ConnPool
}

func CreateRedis(config *Config) (r *Redis) {
	//TODO 进行一些初始化操作
	r = &Redis{
		Config: config,
	}
	pool := NewConnPool(config, nil, 5, 10, 8)
	r.connPool = pool
	return
}

func (r *Redis) GetToString(key string) (res string) {
	// 获取连接对象
	conn := r.connPool.Get()
	result, success := conn.CommandGetResult(get + key + end)
	if success {
		res = string(result)
	}
	// 归还连接对象
	r.connPool.Put(conn)
	return
}

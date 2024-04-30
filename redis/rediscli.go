package redis

type redisCli struct {
	redis Redis
}

// CreateRedisCli 创建连接对象
func CreateRedisCli(config *Config) *redisCli {
	// 初始化中间连接层
	redis := Redis{
		Config: config,
	}
	// 创建连接池
	pool := NewConnPool(config, nil, 5, 10, 8)
	redis.connPool = pool
	// 判断是否开启debug 模式
	if config.Debug {
		redis.AddPostInterceptor(PrintDebug)
	}
	// TODO 还可以继续干很多事情，后面再写
	r := redisCli{
		redis,
	}
	return &r
}

// Get 默认Get方法，返回string类型
func (r *redisCli) Get(key string) (res string) {
	arr := r.GetByte(key)
	res = string(arr)
	return
}

func (r *redisCli) xxx(key string) {

}

func (r *redisCli) GetByte(key string) (res []byte) {
	command := newGetCommand(key, -1)
	res = r.redis.GetByteArr(command)
	return
}

func (r *redisCli) Set(key string, value interface{}) (success bool) {
	cmd := key + " " + convertInterfaceToString(value)
	command := newSetCommand(cmd, -1)
	command.AddArgs(value)
	success = r.redis.Set(command)
	return
}

func (r *redisCli) AddPostInterceptor(interceptor func(redisContext *InterceptorContext)) {
	r.redis.postInterceptors = append(r.redis.postInterceptors, interceptor)
}

package redis

// 拦截器参数传递，拦截器上下文
type interceptorContext struct {
	// 连接对象
	connection *RedisConn
	// 命令信息
	//TODO 后期可能更新为结构体对象
	cmd string
	//其他信息
}

package redis

// InterceptorContext 拦截器参数传递，拦截器上下文
type InterceptorContext struct {
	// 连接对象
	connection *RedisConn
	// 命令信息
	cmd *RedisCommand
	//TODO 后续会可能会继续添加其他信息
}

func (i *InterceptorContext) Connection() *RedisConn {
	return i.connection
}

func (i *InterceptorContext) SetConnection(connection *RedisConn) {
	i.connection = connection
}

func (i *InterceptorContext) Cmd() *RedisCommand {
	return i.cmd
}

func (i *InterceptorContext) SetCmd(cmd *RedisCommand) {
	i.cmd = cmd
}

func newInterceptorContext(connection *RedisConn, cmd *RedisCommand) *InterceptorContext {
	return &InterceptorContext{connection: connection, cmd: cmd}
}

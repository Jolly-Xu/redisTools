package redis

import (
	"fmt"
	"time"
)

/**
redis 中间层主要代码，负责在redisClient 和 connect 中起到承上启下的作用
*/

type Redis struct {
	//  配置类
	Config *Config
	// 连接池
	connPool *ConnPool
	// 前置拦截器，查询前拦截器
	preInterceptors []func(redisContext *InterceptorContext)
	// 后置拦截器，查询后的拦截器
	postInterceptors []func(redisContext *InterceptorContext)
}

func (r *Redis) GetByteArr(cmd *RedisCommand, t *transaction) (res []byte) {
	// 获取连接对象
	var conn *RedisConn
	if t == nil {
		conn = r.connPool.Get()
		// 归还连接对象
		defer r.connPool.Put(conn)
	} else {
		if t.status > 0 {
			return nil
		}
		conn = t.conn
	}
	// 构建一个命令对象
	context := newInterceptorContext(conn, cmd)
	// 执行前置拦截器链
	r.doPreInterceptor(context)
	result, success := conn.CommandGetResult(cmd.Cmd())
	if success {
		res = result
		cmd.SetRes(res)
	} else if t != nil {
		t.status = 2
	}
	// 执行后置拦截链
	r.doPostInterceptor(context)
	// 置空，方便垃圾回收
	context = nil
	return
}

// Set 默认set方法，未携带过期时间
func (r *Redis) Set(cmd *RedisCommand, t *transaction) (success bool) {
	// 获取连接对象
	var conn *RedisConn
	if t == nil {
		conn = r.connPool.Get()
		// 归还连接对象
		defer r.connPool.Put(conn)
	} else {
		if t.status > 0 {
			return false
		}
		conn = t.conn
	}
	// 构建一个命令对象
	context := newInterceptorContext(conn, cmd)
	// 执行前置拦截器链
	r.doPreInterceptor(context)
	if conn.CommandNoResult(cmd.Cmd()) {
		success = true
	} else if t != nil {
		t.status = 2
	}
	// 执行后置拦截链
	r.doPostInterceptor(context)
	// 归还连接对象
	r.connPool.Put(conn)
	// 置空，方便垃圾回收
	context = nil
	return
}

// SetWithExp set顺便设置过期时间，-1表示用不过期，大于0的数表示过期时间，单位为毫秒，
func (r *Redis) SetWithExp(key string, value interface{}, expireTime time.Duration) (success bool) {
	return
}

// AddPreInterceptor 添加前置拦截器
func (r *Redis) AddPreInterceptor(interceptor func(redisContext *InterceptorContext)) {
	r.preInterceptors = append(r.preInterceptors, interceptor)
}

// AddPostInterceptor 添加后置拦截器
func (r *Redis) AddPostInterceptor(interceptor func(redisContext *InterceptorContext)) {
	r.postInterceptors = append(r.postInterceptors, interceptor)
}

func (r *Redis) doPreInterceptor(redisContext *InterceptorContext) {
	for _, interceptor := range r.preInterceptors {
		interceptor(redisContext)
	}
}

func (r *Redis) doPostInterceptor(redisContext *InterceptorContext) {
	for _, interceptor := range r.postInterceptors {
		interceptor(redisContext)
	}
}

func (r *Redis) getRedisConn() *RedisConn {
	return r.connPool.Get()
}

func (r *Redis) putRedisConn(c *RedisConn) {
	r.connPool.Put(c)
}

/*
事务
*/
type transaction struct {
	// 当前事务的连接对象
	conn *RedisConn
	// 当前事务的状态 0表示开启，1表示结束，2表示回滚
	status int
}

func (t *transaction) Open() {
	if t.conn.BeginTransaction() {
		fmt.Println("事务开启失败！")
	}
}

func (t *transaction) Close() {
	if t.conn.EndTransaction() {
		fmt.Println("关闭事务失败")
	}
}

func (t *transaction) Rollback() {
	if t.conn.Rollback() {
		fmt.Println("回滚失败")
	}
}

package redis

import "time"

// redis客户端主要代码

const (
	Get = "GET "
	End = "\n"
	Set = "SET "
)

type Redis struct {
	//  配置类
	Config *Config
	// 连接池
	connPool *ConnPool
	// 前置拦截器，查询前拦截器
	preInterceptors []func(args ...interface{})
	// 后置拦截器，查询后的拦截器
	postInterceptors []func(args ...interface{})
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
	r.doPreInterceptor()
	conn := r.connPool.Get()
	result, success := conn.CommandGetResult(Get + key + End)
	if success {
		res = string(result)
	}
	// 归还连接对象
	r.doPostInterceptor()
	r.connPool.Put(conn)
	return
}

// Set 默认set方法，未携带过期时间
func (r *Redis) Set(key string, value interface{}) (success bool) {
	conn := r.connPool.Get()
	// TODO 需要进行健壮性升级，结构体的类型的处理还是个问题，到底要怎么去存，是提前直接给他序列化了，还是等待用户自己序列化
	cmd := Set + key + " " + value.(string) + END
	if conn.CommandNoResult(cmd) {
		success = true
	}
	return
}

// SetWithExp set顺便设置过期时间，-1表示用不过期，大于0的数表示过期时间，单位为毫秒，
func (r *Redis) SetWithExp(key string, value interface{}, expireTime time.Duration) (success bool) {
	return
}

// AddPreInterceptor 添加前置拦截器
func (r *Redis) AddPreInterceptor(interceptor func(args ...interface{})) {
	r.preInterceptors = append(r.preInterceptors, interceptor)
}

// AddPostInterceptor 添加后置拦截器
func (r *Redis) AddPostInterceptor(interceptor func(args ...interface{})) {
	r.postInterceptors = append(r.postInterceptors, interceptor)
}

func (r *Redis) doPreInterceptor() {
	for _, interceptor := range r.preInterceptors {
		interceptor()
	}
}

func (r *Redis) doPostInterceptor() {
	for _, interceptor := range r.postInterceptors {
		interceptor()
	}
}

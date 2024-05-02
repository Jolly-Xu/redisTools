package redis

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
)

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

func (r *redisCli) xxx(key string) {

}

// Get 默认Get方法，直接转化为string类型返回
func (r *redisCli) Get(key string, arg ...*transaction) (res string) {
	arr := r.GetByte(key, arg...)
	res = string(arr)
	return
}

// GetByte 获取原始二进制结果
func (r *redisCli) GetByte(key string, arg ...*transaction) (res []byte) {
	command := newGetCommand(key, -1)
	var t *transaction = nil
	if len(arg) > 0 {
		t = arg[0]
	}
	res = r.redis.GetByteArr(command, t)
	return
}

// Set 设置任意值，结构体除外
func (r *redisCli) Set(key string, value interface{}, arg ...*transaction) (success bool) {
	cmd := key + " " + convertInterfaceToString(value)
	command := newSetCommand(cmd, -1)
	command.AddArgs(value)
	var t *transaction = nil
	if len(arg) > 0 {
		t = arg[0]
	}
	success = r.redis.Set(command, t)
	return
}

// SetWithStruct 直接设置一个结构体对象
func (r *redisCli) SetWithStruct(key string, value interface{}, arg ...*transaction) (success bool) {
	v, err := msgpack.Marshal(value)
	if err != nil {
		fmt.Printf("不能将[]byte转为%T\n", v)
		fmt.Println(err)
		return
	}
	cmd := key + " " + convertInterfaceToString(v)
	command := newSetCommand(cmd, -1)
	command.AddArgs(value)
	var t *transaction = nil
	if len(arg) > 0 {
		t = arg[0]
	}
	success = r.redis.Set(command, t)
	return
}

// GetWithStruct 获取结构体对象，获取的对象就在传进来的对象中
func (r *redisCli) GetWithStruct(key string, vt interface{}, arg ...*transaction) {
	arr := r.GetByte(key, arg...)
	err := msgpack.Unmarshal(arr, vt)
	if err != nil {
		fmt.Printf("不能将[]byte转为%T\n", vt)
	}
	return
}

// AddPreInterceptor 添加前置拦截器
func (r *redisCli) AddPreInterceptor(interceptor func(redisContext *InterceptorContext)) {
	r.redis.preInterceptors = append(r.redis.preInterceptors, interceptor)
}

// AddPostInterceptor 添加后置拦截器
func (r *redisCli) AddPostInterceptor(interceptor func(redisContext *InterceptorContext)) {
	r.redis.postInterceptors = append(r.redis.postInterceptors, interceptor)
}

// CreateTransaction 创建一个事务对象
func (r *redisCli) CreateTransaction() *transaction {
	t := transaction{conn: nil, status: 0}
	return &t
}

// DoTransaction 事务执行
func (r *redisCli) DoTransaction(t *transaction, f func() (err error)) {
	// 开启事务
	t.conn = r.redis.getRedisConn()
	defer r.redis.putRedisConn(t.conn)
	t.Open()
	err := f()
	if err != nil {
		// 回滚
		t.status = 2
		fmt.Println(err, "事务回滚")
		t.Rollback()
	}
	// 提交事务
	t.status = 1
	t.Close()
}

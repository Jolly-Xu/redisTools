package redis

import (
	"fmt"
	"reflect"
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
	// 先实现结构体和字符串
	//TODO 后续需要继续实现其他类型
	var cmd string
	switch v := value.(type) {
	case string:
		cmd = key + " " + v
	case []byte: // 示例中的其他类型，可以添加更多
		fmt.Printf("value is a byte slice: %s\n", string(v))
	default:
		// 检查是否为结构体类型
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Struct || val.Kind() == reflect.Pointer {
			fmt.Println("警告！存储的是结构体或者指针信息，这样的存储是没有意义的，请先进行序列化！")
			return
		}
	}
	command := newSetCommand(cmd, -1)
	success = r.redis.Set(command)
	return
}

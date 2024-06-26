package redis

import (
	"time"
)

type Config struct {
	// redis连接地址
	IpAddr string
	// 端口
	Port int
	// 网络连接类型，默认Tcp
	NetType string
	//redis-server 协议类型，默认协议为3
	Protocol int
	//用户名
	Username string
	//密码
	Password string
	//数据库的索引
	DataBase int
	// 最大失败尝试次数
	MaxRetries int
	//连接超时时间
	ConnectionTimeOut time.Duration
	//最大读取时间
	ReadTimeOut time.Duration
	// 最大写入时间
	WriteTimeout time.Duration
	//连接池大小
	PoolSize int
	// 最大空闲连接数量
	MaxIdleConn int
	// 最大活动连接数量
	MaxActiveConn int
	// 是否开启Debug模式，开启日志打印,默认为false
	Debug bool
}

func FastConfig(addr string, port int, password string) *Config {
	return &Config{
		IpAddr:   addr,
		Port:     port,
		Password: password,
		NetType:  "tcp",
		DataBase: 16,
		Debug:    true,
	}
}

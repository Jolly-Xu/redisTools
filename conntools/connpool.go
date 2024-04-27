package conntools

import (
	"redisTools/redis"
	"strconv"
	"sync"
)

type ConnPool struct {
	pool         chan *RedisConn
	maxOpenConns int
	maxIdleConns int
	mu           sync.Mutex
	factory      func() *RedisConn
}

func (p *ConnPool) SetFactory(factory func() *RedisConn) {
	p.factory = factory
}

func (p *ConnPool) MaxOpenConns() int {
	return p.maxOpenConns
}

func (p *ConnPool) SetMaxOpenConns(maxOpenConns int) {
	p.maxOpenConns = maxOpenConns
}

func (p *ConnPool) MaxIdleConns() int {
	return p.maxIdleConns
}

func (p *ConnPool) SetMaxIdleConns(maxIdleConns int) {
	p.maxIdleConns = maxIdleConns
}

// DefaultFactory 默认创建连接工厂
func DefaultFactory(config *redis.Config) *RedisConn {
	dst := config.IpAddr + ":" + strconv.Itoa(config.Port)
	conn := NewRedisConn(config.NetType, dst)
	return conn
}

func NewConnPool(config *redis.Config, factory func() *RedisConn, initialOpenConns, maxOpenConns, maxIdleConns int) *ConnPool {
	// 创建一个连接池
	pool := &ConnPool{
		pool:         make(chan *RedisConn, maxOpenConns),
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		factory:      factory,
	}
	df := factory == nil
	// 初始化若干个连接对象
	for i := 0; i < initialOpenConns; {
		// 使用工厂方法对连接进行创建，如果没有传入，使用默认实现的工厂对象
		var conn *RedisConn
		if df {
			conn = DefaultFactory(config)
		} else {
			//TODO 执行自定义的连接创建工厂
		}
		// 直到创建完成initialOpenConns个连接对象才停止循环
		if conn.Auth(config.Password) {
			pool.Put(conn)
			i++
		}
	}
	return pool
}

func (p *ConnPool) Get() *RedisConn {
	select {
	case item := <-p.pool:
		return item
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		if len(p.pool) >= p.maxOpenConns {
			return nil
		}
		conn := p.factory()
		return conn
	}
}

func (p *ConnPool) Put(item *RedisConn) {
	select {
	case p.pool <- item:
		if len(p.pool) > p.maxIdleConns {
			item.conn.Close()
		}
	default:
		item.conn.Close()
	}
}

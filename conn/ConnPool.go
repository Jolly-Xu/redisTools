package conn

import (
	"fmt"
	"net"
	"sync"
)


type ConnPool struct {
	pool         chan net.Conn
	maxOpenConns int
	maxIdleConns int
	mu           sync.Mutex
	factory      func() (net.Conn, error)
}

func NewConnPool(factory func() (net.Conn, error), initialOpenConns, maxOpenConns, maxIdleConns int) *ConnPool {
	pool := &ConnPool{
		pool:         make(chan net.Conn, maxOpenConns),
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		factory:      factory,
	}
	for i := 0; i < initialOpenConns; i++ {
		conn, err := factory()
		if err != nil {
			fmt.Println("Error creating initial connections:", err)
			return nil
		}
		pool.Put(conn)
	}
	return pool
}

func (p *ConnPool) Get() (net.Conn, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		if len(p.pool) >= p.maxOpenConns {
			return nil, fmt.Errorf("connection pool exhausted")
		}
		conn, err := p.factory()
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (p *ConnPool) Put(conn net.Conn) {
	select {
	case p.pool <- conn:
		if len(p.pool) > p.maxIdleConns {
			conn.Close()
		}
	default:
		conn.Close()
	}
}

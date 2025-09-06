package pool

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	ErrPoolClosed          = errors.New("connection pool is closed")
	ErrTooManyConnections  = errors.New("too many active connections")
	ErrConnectionUnhealthy = errors.New("connection is unhealthy")
)

// PooledConnection 池化的连接
type PooledConnection struct {
	conn      net.Conn
	lastUsed  time.Time // 最后使用时间
	createdAt time.Time // 创建时间
	inUse     bool      // 是否正在使用
}

// ConnectionPool 连接池
type ConnectionPool struct {
	mu              sync.RWMutex
	connections     map[string][]*PooledConnection // 地址 -> 连接列表
	maxIdle         int                            // 最大空闲连接数（每个地址）
	maxActive       int                            // 最大活跃连接数（每个地址）
	idleTimeout     time.Duration                  // 空闲连接超时时间
	maxConnLifetime time.Duration                  // 连接最大生存时间
	closed          bool                           // 池是否已关闭
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxIdle         int           // 每个地址的最大空闲连接数，默认5
	MaxActive       int           // 每个地址的最大活跃连接数，默认20
	IdleTimeout     time.Duration // 空闲连接超时，默认5分钟
	MaxConnLifetime time.Duration // 连接最大生存时间，默认30分钟
}

// DefaultPoolConfig 返回默认的连接池配置
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxIdle:         5,
		MaxActive:       20,
		IdleTimeout:     5 * time.Minute,
		MaxConnLifetime: 30 * time.Minute,
	}
}

// NewConnectionPool 创建新的连接池
func NewConnectionPool(config *PoolConfig) *ConnectionPool {
	if config == nil {
		config = DefaultPoolConfig()
	}

	// DONE 1: 完成连接池的初始化
	// 提示：
	// - 创建ConnectionPool结构体实例
	// - 初始化connections map
	// - 从config中设置各种参数
	// - 设置closed为false
	// - 返回连接池实例
	return &ConnectionPool{
		mu:              sync.RWMutex{},
		connections:     make(map[string][]*PooledConnection),
		maxIdle:         config.MaxIdle,
		maxActive:       config.MaxActive,
		idleTimeout:     config.IdleTimeout,
		maxConnLifetime: config.MaxConnLifetime,
		closed:          false,
	}
}

// Get 获取连接
func (p *ConnectionPool) Get(address string) (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, ErrPoolClosed
	}

	// TODO 2: 实现获取连接的逻辑
	// 步骤：
	// 1. 检查是否有可用的空闲连接（inUse == false）
	// 2. 如果找到空闲连接，检查是否健康（调用p.isHealthy）
	// 3. 如果连接健康，标记为使用中（inUse = true），更新lastUsed时间，返回连接
	// 4. 如果没有空闲连接，检查活跃连接数是否超过maxActive
	// 5. 如果未超过限制，创建新连接（调用p.createConnection）
	// 6. 将新连接包装成PooledConnection并添加到池中
	//
	// 提示：你可以先写一个简单版本，只处理创建新连接的情况

	// mu              sync.RWMutex
	// connections     map[string][]*PooledConnection // 地址 -> 连接列表
	// maxIdle         int                            // 最大空闲连接数（每个地址）
	// maxActive       int                            // 最大活跃连接数（每个地址）
	// idleTimeout     time.Duration                  // 空闲连接超时时间
	// maxConnLifetime time.Duration                  // 连接最大生存时间
	// closed          bool                           // 池是否已关闭

	pools, exists := p.connections[address]
	if !exists {
		// 完全没有链接， 创建一个 并且放到池子里

		// conn      net.Conn
		// lastUsed  time.Time // 最后使用时间
		// createdAt time.Time // 创建时间
		// inUse     bool      // 是否正在使用
		conn, err := p.createConnection(address)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to dial tcp address %s", address))
		}
		t := &PooledConnection{
			conn:      conn,
			lastUsed:  time.Now(),
			createdAt: time.Now(),
			inUse:     true,
		}
		p.connections[address] = []*PooledConnection{t}
		return t.conn, nil
	}
	// 尝试获取空闲链接
	for i, conn := range pools {
		if conn.inUse == false && p.isHealthy(conn) {
			pools[i].lastUsed = time.Now()
			pools[i].inUse = true
			return conn.conn, nil
		}
	}
	// 没有空闲链接, 尝试创建
	if len(p.connections[address]) < p.maxActive {
		// 没超过上限 还能继续创建新的链接
		conn, err := p.createConnection(address)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to dial tcp address %s", address))
		}
		t := &PooledConnection{
			conn:      conn,
			lastUsed:  time.Now(),
			createdAt: time.Now(),
			inUse:     true,
		}
		p.connections[address] = append(p.connections[address], t)
		return t.conn, nil
	} else {
		// 疑问1 没想好， 如果所有链接都在使用， 我们也碰到了 MaxActive 的上限， 这首业务侧来请求链接
		// 怎么处理呢？
	}
	return nil, errors.New("TODO: implement Get method")

}

// Put 归还连接
func (p *ConnectionPool) Put(address string, conn net.Conn) error {
	// TODO 3: 实现归还连接的逻辑
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		// 疑问2  如果链接池已经关闭了， 应该把 还回来的conn 也关闭掉？
		return conn.Close()
	}
	for _, c := range p.connections[address] {
		if c.conn == conn {
			c.inUse = false
			c.lastUsed = time.Now()
			break
		}
	}
	idleConns := 0

	for _, c := range p.connections[address] {
		if c.inUse == false {
			idleConns++
		}
	}
	connNeedsToBeClosed := p.maxIdle - idleConns
	for _, c := range p.connections[address] {
		if connNeedsToBeClosed <= 0 {
			break
		}
		if c.inUse == false {
			// 疑问3 这里只要关闭就行了吗？ 感觉这个array 让我剔除操作很困难 不能像map 一样直接delete
			c.conn.Close()
			connNeedsToBeClosed--
		}
	}
	return nil
	// 步骤：
	// 1. 加锁保护并发访问
	// 2. 检查连接池是否已关闭
	// 3. 在connections[address]中查找对应的PooledConnection
	// 4. 找到后将inUse设置为false，更新lastUsed时间
	// 5. 检查空闲连接数是否超过maxIdle，如果超过就关闭连接
	//
	// 先实现最简单的版本：找到连接并标记为空闲

	// return errors.New("TODO: implement Put method")
}

// Close 关闭连接池
func (p *ConnectionPool) Close() error {
	// TODO 4: 实现关闭连接池的逻辑
	// 步骤：
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	for _, conns := range p.connections {
		for _, conn := range conns {
			conn.conn.Close()
		}
	}
	p.connections = nil
	p.closed = true
	return nil
	// 1. 加锁
	// 2. 检查是否已经关闭，避免重复关闭
	// 3. 设置closed为true
	// 4. 遍历所有连接，调用conn.Close()关闭连接
	// 5. 清空connections map

	//return errors.New("TODO: implement Close method")
}

// 辅助方法 - 我已经为你实现了这些，你可以直接使用

// createConnection 创建到指定地址的新连接
func (p *ConnectionPool) createConnection(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

// isHealthy 检查连接是否健康
func (p *ConnectionPool) isHealthy(pooledConn *PooledConnection) bool {
	if pooledConn == nil || pooledConn.conn == nil {
		return false
	}

	// 检查连接是否超过最大生存时间
	if time.Since(pooledConn.createdAt) > p.maxConnLifetime {
		return false
	}

	// 检查空闲时间是否超时
	if !pooledConn.inUse && time.Since(pooledConn.lastUsed) > p.idleTimeout {
		return false
	}

	// TODO 5 (可选): 添加更高级的健康检查
	// 你可以尝试向连接写入/读取一些数据来检查连接是否真的可用
	// 但现在先用时间检查就够了

	//  疑问4 这里是不是需要一些剔除不健康链接的操作呢？ 什么是这里的best practice？
	return true
}

// getActiveCount 获取指定地址的活跃连接数
func (p *ConnectionPool) getActiveCount(address string) int {
	count := 0
	for _, pooledConn := range p.connections[address] {
		if pooledConn.inUse {
			count++
		}
	}
	return count
}

// getIdleCount 获取指定地址的空闲连接数
func (p *ConnectionPool) getIdleCount(address string) int {
	count := 0
	for _, pooledConn := range p.connections[address] {
		if !pooledConn.inUse {
			count++
		}
	}
	return count
}

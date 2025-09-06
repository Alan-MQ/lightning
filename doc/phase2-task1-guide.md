# Phase 2 Task 1: 连接池基础实现指南

## 🎯 任务目标
实现连接池的核心功能，预计用时：1-2小时

## 📋 需要完成的TODO

### TODO 1: NewConnectionPool 函数 (15分钟)
**位置**: `pool/connection_pool.go:52`

**任务**: 完成连接池的初始化
```go
func NewConnectionPool(config *PoolConfig) *ConnectionPool {
    if config == nil {
        config = DefaultPoolConfig()
    }

    // 你需要实现这里
    pool := &ConnectionPool{
        // 填充所有字段
    }
    
    return pool
}
```

**提示**:
- 使用 `make(map[string][]*PooledConnection)` 初始化connections
- 从config中复制所有配置参数
- closed字段设置为false

**验证**: 运行 `go test -v ./pool -run TestBasicConnectionPool`

### TODO 2: Get 方法 (30-40分钟)
**位置**: `pool/connection_pool.go:60`

**任务**: 实现获取连接的逻辑

**分步实现**:

**Step 2.1: 简单版本 (15分钟)**
```go
// 先实现最简单的版本：每次都创建新连接
conn, err := p.createConnection(address)
if err != nil {
    return nil, err
}

// 创建PooledConnection并添加到池中
pooledConn := &PooledConnection{
    conn:      conn,
    lastUsed:  time.Now(),
    createdAt: time.Now(), 
    inUse:     true,
}

// 添加到connections map
p.connections[address] = append(p.connections[address], pooledConn)
return conn, nil
```

**Step 2.2: 添加连接复用 (15-20分钟)**
```go
// 首先检查是否有空闲连接
connections := p.connections[address]
for _, pooledConn := range connections {
    if !pooledConn.inUse && p.isHealthy(pooledConn) {
        // 找到可用连接
        pooledConn.inUse = true
        pooledConn.lastUsed = time.Now()
        return pooledConn.conn, nil
    }
}

// 没有空闲连接，创建新连接 (使用Step 2.1的代码)
```

**Step 2.3: 添加连接数限制 (5分钟)**
```go
// 在创建新连接前检查限制
activeCount := p.getActiveCount(address)
if activeCount >= p.maxActive {
    return nil, ErrTooManyConnections
}
```

### TODO 3: Put 方法 (20-30分钟)
**位置**: `pool/connection_pool.go:77`

**任务**: 实现归还连接的逻辑

**分步实现**:
```go
func (p *ConnectionPool) Put(address string, conn net.Conn) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        conn.Close()
        return ErrPoolClosed
    }

    // 查找对应的PooledConnection
    connections := p.connections[address]
    for _, pooledConn := range connections {
        if pooledConn.conn == conn {
            // 找到了，标记为空闲
            pooledConn.inUse = false
            pooledConn.lastUsed = time.Now()
            return nil
        }
    }

    // 没找到，可能是外部创建的连接，直接关闭
    conn.Close()
    return nil
}
```

### TODO 4: Close 方法 (15分钟)
**位置**: `pool/connection_pool.go:84`

**任务**: 实现关闭连接池

```go
func (p *ConnectionPool) Close() error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        return nil // 已经关闭了
    }

    p.closed = true

    // 关闭所有连接
    for address, connections := range p.connections {
        for _, pooledConn := range connections {
            pooledConn.conn.Close()
        }
        delete(p.connections, address)
    }

    return nil
}
```

### TODO 5: 高级健康检查 (可选，20分钟)
**位置**: `pool/connection_pool.go:108`

**任务**: 添加更严格的连接健康检查

**提示**: 可以尝试读写少量数据来检测连接状态，但要小心不要影响正常通信。

## 🧪 测试步骤

### 第一步：基础测试
```bash
go test -v ./pool -run TestBasicConnectionPool
```

### 第二步：创建真实服务器测试
你可以修改我们的 `main/main.go`，启动服务器然后测试连接池：

```go
// 在另一个终端运行
go run main/main.go

// 然后修改测试中的地址为实际的服务器地址
```

### 第三步：完整测试
```bash
go test -v ./pool
```

## 💡 实现技巧

### 1. 调试技巧
在关键位置添加日志：
```go
fmt.Printf("Pool stats: %+v\n", p.Stats())
```

### 2. 错误处理
注意区分不同类型的错误：
- 网络连接错误
- 池已满错误  
- 池已关闭错误

### 3. 并发安全
- 所有操作都要加锁
- 注意锁的范围，避免死锁

## 🎯 完成标准

当你完成所有TODO后，应该能够：
- [x] 创建连接池
- [x] 获取连接（支持复用）
- [x] 归还连接
- [x] 关闭连接池
- [x] 所有测试通过

## 🚀 下一步

完成这个任务后，我们将继续实现：
- 连接清理机制
- 统计信息
- 与RPC客户端集成

准备好了吗？开始编码吧！有任何问题随时问我。
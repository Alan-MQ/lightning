package pool

import (
	"testing"
)

// TestBasicConnectionPool 基础连接池测试
func TestBasicConnectionPool(t *testing.T) {
	// 创建连接池
	config := DefaultPoolConfig()
	pool := NewConnectionPool(config)

	if pool == nil {
		t.Fatal("NewConnectionPool returned nil")
	}

	defer pool.Close()

	// TODO: 当你完成TODO 1后，这个测试应该通过
	t.Log("✅ Connection pool creation test passed")
}

// TestGetConnection 测试获取连接
func TestGetConnection(t *testing.T) {
	pool := NewConnectionPool(nil)
	defer pool.Close()

	// 测试获取连接 - 需要一个真实的地址来测试
	// 你可以启动我们之前的RPC服务器来测试，或者先跳过这个测试

	// TODO: 当你完成TODO 2后，取消下面的注释来测试
	conn, err := pool.Get("localhost:0") // 使用一个无效地址来测试错误处理
	if err == nil {
		t.Errorf("Expected error when connecting to invalid address")
		conn.Close()
	}
	t.Log("✅ Connection error handling test passed")

	t.Log("⏳ Get connection test skipped - implement TODO 2 first")
}

// TestPutConnection 测试归还连接
func TestPutConnection(t *testing.T) {
	// TODO: 当你完成TODO 2和TODO 3后，实现这个测试
	pool := NewConnectionPool(DefaultPoolConfig())
	addr := "localhost:9090"
	conn, err := pool.Get(addr)
	if err == nil {
		t.Errorf("Expected error when connecting to invalid address")
		conn.Close()
	}
	pool.Put(addr, conn)
	newConn, err := pool.Get(addr)
	if err == nil {
		t.Errorf("Expected error when connecting to invalid address")
		newConn.Close()
	}
	// 测试步骤：
	// 1. 创建连接池
	// 2. 获取一个连接
	// 3. 归还连接
	// 4. 再次获取连接，应该复用之前的连接

	t.Log("⏳ Put connection test - implement TODO 3 first")
}

// TestClosePool 测试关闭连接池
func TestClosePool(t *testing.T) {
	pool := NewConnectionPool(nil)

	err := pool.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// 关闭后获取连接应该失败
	_, err = pool.Get("localhost:8080")
	if err != ErrPoolClosed {
		t.Errorf("Expected ErrPoolClosed, got: %v", err)
	}

	pool.Close() // 确保清理
	t.Log("⏳ Close pool test - implement TODO 4 first")
}

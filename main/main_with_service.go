package main

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/Alan-MQ/lightning/codec"
	"github.com/Alan-MQ/lightning/examples"
	"github.com/Alan-MQ/lightning/server"
)

func startServerWithService(addr chan string) {
	// 创建服务器
	s := server.NewServer()

	// 注册Calculator服务
	calc := &examples.Calculator{}
	stringSvc := &examples.StringService{}
	err := s.Register(calc)
	if err != nil {
		log.Fatal("Failed to register Calculator service:", err)
	}
	err = s.Register(stringSvc)
	if err != nil {
		log.Fatal("Failed to register Calculator service:", err)
	}
	log.Println("Calculator service registered successfully")

	// 监听端口
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()

	// 接受连接
	s.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServerWithService(addr)

	// 等待服务器启动
	time.Sleep(time.Second)

	// 客户端连接
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	// 发送配置选项
	_ = json.NewEncoder(conn).Encode(server.DefaultOption)
	cc := codec.NewGobCodec(conn)

	// 测试Calculator.Add方法
	log.Println("=== Testing Calculator.Add ===")
	{
		h := &codec.Header{
			ServiceMethod: "Calculator.Add",
			Seq:           1,
		}

		args := &examples.AddArgs{A: 10, B: 20}
		err := cc.Write(h, args)
		if err != nil {
			log.Fatal("write error:", err)
		}

		err = cc.ReadHeader(h)
		if err != nil {
			log.Fatal("read header error:", err)
		}

		var reply examples.AddReply
		err = cc.ReadBody(&reply)
		if err != nil {
			log.Fatal("read body error:", err)
		}

		log.Printf("10 + 20 = %d", reply.Sum)

		if h.Error != "" {
			log.Printf("call error: %s", h.Error)
		}
	}

	// 测试Calculator.Multiply方法
	log.Println("=== Testing Calculator.Multiply ===")
	{
		h := &codec.Header{
			ServiceMethod: "Calculator.Multiply",
			Seq:           2,
		}

		args := &examples.MultiplyArgs{A: 6, B: 7}
		err := cc.Write(h, args)
		if err != nil {
			log.Fatal("write error:", err)
		}

		err = cc.ReadHeader(h)
		if err != nil {
			log.Fatal("read header error:", err)
		}

		var reply examples.MultiplyReply
		err = cc.ReadBody(&reply)
		if err != nil {
			log.Fatal("read body error:", err)
		}

		log.Printf("6 * 7 = %d", reply.Product)

		if h.Error != "" {
			log.Printf("call error: %s", h.Error)
		}
	}

	// 测试String Reverse方法
	log.Println("=== Testing String.Reversed ===")
	{
		h := &codec.Header{
			ServiceMethod: "StringService.Reverse",
			Seq:           3,
		}

		args := &examples.ReverseArgs{S: "Hello, World!"}
		err := cc.Write(h, args)
		if err != nil {
			log.Fatal("write error:", err)
		}

		err = cc.ReadHeader(h)
		if err != nil {
			log.Fatal("read header error:", err)
		}

		var reply examples.ReverseReply
		err = cc.ReadBody(&reply)
		if err != nil {
			log.Fatal("read body error:", err)
		}

		log.Printf("Hello, World! reversed  = %s", reply.Reversed)

		if h.Error != "" {
			log.Printf("call error: %s", h.Error)
		}
	}
	log.Println("=== All tests completed ===")
}

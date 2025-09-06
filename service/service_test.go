package service

import (
	"reflect"
	"testing"

	"github.com/Alan-MQ/lightning/examples"
)

func TestServiceRegistry(t *testing.T) {
	// 创建服务注册中心
	registry := NewServiceRegistry()

	// 注册Calculator服务
	calc := &examples.Calculator{}
	err := registry.Register(calc)
	if err != nil {
		t.Fatalf("Failed to register Calculator service: %v", err)
	}

	// 测试服务列表
	services := registry.ListServices()
	if len(services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(services))
	}

	calcMethods, exists := services["Calculator"]
	if !exists {
		t.Fatal("Calculator service not found in service list")
	}

	// 验证只有符合RPC规范的方法被注册（Add和Multiply）
	expectedMethods := map[string]bool{
		"Add":      false,
		"Multiply": false,
	}

	for _, method := range calcMethods {
		if _, ok := expectedMethods[method]; ok {
			expectedMethods[method] = true
		} else {
			t.Errorf("Unexpected method found: %s", method)
		}
	}

	for method, found := range expectedMethods {
		if !found {
			t.Errorf("Expected method %s not found", method)
		}
	}
}

func TestServiceCall(t *testing.T) {
	// 创建服务注册中心并注册服务
	registry := NewServiceRegistry()
	calc := &examples.Calculator{}
	err := registry.Register(calc)
	if err != nil {
		t.Fatalf("Failed to register Calculator service: %v", err)
	}

	// 测试Add方法调用
	args := &examples.AddArgs{A: 10, B: 20}
	reply := &examples.AddReply{}

	argsValue := reflect.ValueOf(args)
	replyValue := reflect.ValueOf(reply)

	err = registry.Call("Calculator.Add", argsValue, replyValue)
	if err != nil {
		t.Fatalf("Failed to call Calculator.Add: %v", err)
	}

	if reply.Sum != 30 {
		t.Errorf("Expected sum to be 30, got %d", reply.Sum)
	}

	// 测试Multiply方法调用
	multiplyArgs := &examples.MultiplyArgs{A: 6, B: 7}
	multiplyReply := &examples.MultiplyReply{}

	argsValue = reflect.ValueOf(multiplyArgs)
	replyValue = reflect.ValueOf(multiplyReply)

	err = registry.Call("Calculator.Multiply", argsValue, replyValue)
	if err != nil {
		t.Fatalf("Failed to call Calculator.Multiply: %v", err)
	}

	if multiplyReply.Product != 42 {
		t.Errorf("Expected product to be 42, got %d", multiplyReply.Product)
	}
}

func TestServiceCallErrors(t *testing.T) {
	registry := NewServiceRegistry()
	calc := &examples.Calculator{}
	registry.Register(calc)

	args := &examples.AddArgs{A: 1, B: 2}
	reply := &examples.AddReply{}
	argsValue := reflect.ValueOf(args)
	replyValue := reflect.ValueOf(reply)

	// 测试服务不存在
	err := registry.Call("NonExistentService.Add", argsValue, replyValue)
	if err == nil {
		t.Error("Expected error for non-existent service")
	}

	// 测试方法不存在
	err = registry.Call("Calculator.NonExistentMethod", argsValue, replyValue)
	if err == nil {
		t.Error("Expected error for non-existent method")
	}

	// 测试格式错误的服务方法
	err = registry.Call("InvalidFormat", argsValue, replyValue)
	if err == nil {
		t.Error("Expected error for invalid service/method format")
	}
}
package service

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// MethodType 保存方法的元信息
type MethodType struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

// Service 表示一个注册的服务
type Service struct {
	name   string                 // 服务名称
	rcvr   reflect.Value          // 服务实例
	typ    reflect.Type           // 服务类型
	method map[string]*MethodType // 服务的方法集合
}

// ServiceRegistry 服务注册中心
type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]*Service
}

// NewServiceRegistry 创建新的服务注册中心
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*Service),
	}
}

// isExported 检查名称是否为导出的（首字母大写）
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// isExportedOrBuiltin 检查类型是否为导出的或内置类型
func isExportedOrBuiltin(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// Register 注册服务（使用类型名作为服务名）
func (registry *ServiceRegistry) Register(rcvr interface{}) error {
	return registry.register(rcvr, "", false)
}

// RegisterName 注册服务（使用指定名称作为服务名）
func (registry *ServiceRegistry) RegisterName(name string, rcvr interface{}) error {
	return registry.register(rcvr, name, true)
}

// register 核心注册逻辑
func (registry *ServiceRegistry) register(rcvr interface{}, name string, useName bool) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	s := new(Service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()

	if useName {
		sname = name
	}

	if sname == "" {
		errMsg := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(errMsg)
		return errors.New(errMsg)
	}

	if !ast.IsExported(sname) && !useName {
		errMsg := "rpc.Register: type " + sname + " is not exported"
		log.Print(errMsg)
		return errors.New(errMsg)
	}

	s.name = sname
	s.method = make(map[string]*MethodType)

	// 扫描所有方法
	for m := 0; m < s.typ.NumMethod(); m++ {
		method := s.typ.Method(m)
		mtype := method.Type
		mname := method.Name

		// 只处理导出的方法
		if !isExported(mname) {
			continue
		}

		// 检查方法签名：必须有3个参数 (receiver, args, reply)
		if mtype.NumIn() != 3 {
			log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			continue
		}

		// 检查第二个参数（args）类型
		argType := mtype.In(1)
		if !isExportedOrBuiltin(argType) {
			log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			continue
		}

		// 检查第三个参数（reply）必须是指针
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			continue
		}

		// 检查reply指针指向的类型是否导出
		if !isExportedOrBuiltin(replyType) {
			log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			continue
		}

		// 检查返回值：必须只有一个返回值，且为error类型
		if mtype.NumOut() != 1 {
			log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			continue
		}

		if returnType := mtype.Out(0); returnType != reflect.TypeOf((*error)(nil)).Elem() {
			log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			continue
		}

		// 方法验证通过，添加到服务中
		s.method[mname] = &MethodType{Method: method, ArgType: argType, ReplyType: replyType}
	}

	if len(s.method) == 0 {
		errMsg := "rpc.Register: type " + sname + " has no exported methods of suitable type"
		log.Print(errMsg)
		return errors.New(errMsg)
	}

	registry.services[sname] = s
	return nil
}

// GetService 获取指定名称的服务
func (registry *ServiceRegistry) GetService(serviceName string) (*Service, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	service, ok := registry.services[serviceName]
	return service, ok
}

// Call 动态调用指定的服务方法
func (registry *ServiceRegistry) Call(serviceMethod string, argv, replyv reflect.Value) error {
	// 解析服务名和方法名
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		return fmt.Errorf("service/method request ill-formed: %s", serviceMethod)
	}

	serviceName := serviceMethod[:dot]
	methodName := serviceMethod[dot+1:]

	// 查找服务
	service, ok := registry.GetService(serviceName)
	if !ok {
		return fmt.Errorf("service not found: %s", serviceName)
	}

	// 查找方法
	method, ok := service.method[methodName]
	if !ok {
		return fmt.Errorf("method not found: %s.%s", serviceName, methodName)
	}

	// 动态调用方法
	function := method.Method.Func
	returnValues := function.Call([]reflect.Value{service.rcvr, argv, replyv})

	// 检查返回的错误
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}

// ListServices 列出所有已注册的服务及其方法
func (registry *ServiceRegistry) ListServices() map[string][]string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	result := make(map[string][]string)
	for serviceName, service := range registry.services {
		var methods []string
		for methodName := range service.method {
			methods = append(methods, methodName)
		}
		result[serviceName] = methods
	}
	return result
}

// Method 获取指定的方法信息
func (s *Service) Method(methodName string) (*MethodType, bool) {
	method, ok := s.method[methodName]
	return method, ok
}
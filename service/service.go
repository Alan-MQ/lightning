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

// methodType 保存方法的元信息
type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

// service 表示一个注册的服务
type service struct {
	name   string                 // 服务名称
	rcvr   reflect.Value          // 服务实例
	typ    reflect.Type           // 服务类型
	method map[string]*methodType // 服务的方法集合
}

// ServiceRegistry 服务注册中心
type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]*service
}

// NewServiceRegistry 创建新的服务注册中心
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*service),
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
package dbmanager

import (
	"context"
	"reflect"
	"sync"
)

// DatabaseConnection 表示数据库连接的接口
type DatabaseConnection interface {
}

// DbMap 是一个泛型结构，用于管理不同租户的数据库连接
type DbMap[T DatabaseConnection] struct {
	sync.RWMutex
	m map[any]T
}

// NewDbMap 返回一个新的 DbMap 实例
func NewDbMap[T DatabaseConnection]() *DbMap[T] {
	return &DbMap[T]{
		m: make(map[any]T),
	}
}

// AddMap 添加一个新的数据库连接到 DbMap
func (l *DbMap[T]) AddMap(tenantId any, db T) {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.m[tenantId]; ok {
		l.m[tenantId] = db
	}
}

// GetMap 从 DbMap 获取指定租户的数据库连接
func (l *DbMap[T]) GetMap(tenantId any) T {
	l.RLock()
	defer l.RUnlock()

	var zero T
	if tenantId == "" {
		return zero
	}

	db, ok := l.m[tenantId]
	if !ok {
		return zero
	}

	v := reflect.ValueOf(db)
	if !v.IsValid() {
		return zero
	}
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return zero
	}

	return db
}

// DelMap 从 DbMap 删除指定租户的数据库连接
func (l *DbMap[T]) DelMap(tenantId any) {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.m[tenantId]; ok {
		delete(l.m, tenantId)
	}
}

// WithOptionTenantId 用于获取或创建指定租户的数据库连接
func (l *DbMap[T]) WithOptionTenantId(
	tenantId any,
	onInit func(id any) (T, error),
	onClose func(T) error,
) (T, error) {

	db := l.GetMap(tenantId)
	if !isNil(db) {
		return db, nil
	}

	db, err := onInit(tenantId)
	if err != nil {
		return db, err
	}

	l.AddMap(tenantId, db)

	// 注册关闭回调
	if onClose != nil {
		go func() {
			<-context.Background().Done() // 假设 App 关闭时触发
			onClose(db)
		}()
	}

	return db, nil
}

// isNil 判断泛型值是否为 nil（仅适用于指针、接口等可为 nil 的类型）
func isNil[T any](v T) bool {
	if any(v) == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

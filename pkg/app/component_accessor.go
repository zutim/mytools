// pkg/app/component_accessor.go
package app

import (
	"fmt"
	"mytool3/pkg/dbmanager"
)

// ComponentAccessor 泛型封装，用于访问任意类型的组件
type ComponentAccessor[T any] struct {
	name string
}

// NewComponentAccessor 创建一个新的泛型组件访问器
func NewComponentAccessor[T any](name string) *ComponentAccessor[T] {
	return &ComponentAccessor[T]{name: name}
}

// Get 获取指定租户的组件实例（带租户隔离）
func (ca *ComponentAccessor[T]) Get(tenantId any) (T, error) {
	dbMap := dbmanager.NewDbMap[T]()

	return dbMap.WithOptionTenantId(
		tenantId,
		func(id any) (T, error) {
			obj, err := GetDefaultApp().GetComponent(ca.name).Init(id)
			if err != nil {
				var zero T
				return zero, fmt.Errorf("failed to init component %s: %w", ca.name, err)
			}
			return obj.(T), nil
		},
		func(obj T) error { // 注册 onClose 回调
			if closer, ok := GetDefaultApp().GetComponent(ca.name).(interface {
				Close(any) error
			}); ok {
				return closer.Close(tenantId)
			}
			return nil
		},
	)
}

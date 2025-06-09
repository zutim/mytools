// pkg/app/mongo_component.go
package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/zutim/mongo"
	"sync"
)

type MongoComponent struct {
	instances map[any]*mongo.MongoClient
	mu        sync.Mutex
}

func NewMongoComponent() Component {
	return &MongoComponent{
		instances: make(map[any]*mongo.MongoClient),
	}
}

func (m *MongoComponent) Name() string {
	return "mongo"
}

// pkg/app/mongo_component.go
func (m *MongoComponent) Init(tenantId any) (any, error) {
	var uri string
	app := GetDefaultApp()
	cfg := app.GetGlobalConfig().Mongo
	tcfg := app.TenantConfig(tenantId).Mongo

	if tenantId == 0 {
		if cfg.URI == "" {
			return nil, errors.New("no global config")
		}
		uri = cfg.URI
	} else {
		if tcfg.URI == "" {
			return nil, errors.New(fmt.Sprintf("no tenant config %v", tenantId))
		}
		uri = tcfg.URI
	}

	maxPool := int(cfg.MaxPoolSize)
	if maxPool == 0 {
		maxPool = 20
	}

	return mongo.NewMongo(uri, maxPool), nil

}

func (m *MongoComponent) Close(tenantId any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.instances[tenantId]
	if !exists {
		return fmt.Errorf("no mongodb client instance for tenant %v", tenantId)
	}

	// 调用底层关闭方法（根据你使用的驱动决定）
	if err := client.Client.Disconnect(context.Background()); err != nil {
		return err
	}

	// 从映射中删除
	delete(m.instances, tenantId)

	return nil
}

func (m *MongoComponent) HealthCheck() bool {
	return true
}

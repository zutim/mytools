// pkg/app/redis_init.go
package app

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
)

type RedisComponent struct {
	instances map[any]*redis.Client
	mu        sync.Mutex
}

func NewRedisComponent() Component {
	return &RedisComponent{
		instances: make(map[any]*redis.Client),
	}
}

func (r *RedisComponent) Name() string {
	return "redis"
}

func (r *RedisComponent) Init(tenantId any) (any, error) {

	var redisCfg struct {
		Addr     string
		Password string
		DB       int
	}

	app := GetDefaultApp()
	gloablCfg := app.GetGlobalConfig().Redis
	tenantCfg := app.TenantConfig(tenantId).Redis

	if tenantId == 0 {
		if gloablCfg.Addr == "" {
			return nil, errors.New("no global config")
		}
		redisCfg.Addr = gloablCfg.Addr
		redisCfg.Password = gloablCfg.Password
		redisCfg.DB = gloablCfg.DB
	} else {
		if tenantCfg.Addr == "" {
			return nil, errors.New(fmt.Sprintf("no tenant config %v", tenantId))
		}
		redisCfg.Addr = tenantCfg.Addr
		redisCfg.Password = tenantCfg.Password
		redisCfg.DB = tenantCfg.DB
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	return client, nil
}

func (r *RedisComponent) Close(tenantId any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, exists := r.instances[tenantId]
	if !exists {
		return fmt.Errorf("no redis client instance for tenant %v", tenantId)
	}

	// 调用底层 redis.Client 的 Close() 方法
	if err := client.Close(); err != nil {
		return err
	}

	// 从实例映射中删除
	delete(r.instances, tenantId)

	return nil
}

func (r *RedisComponent) HealthCheck() bool {
	return true
}

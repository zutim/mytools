// pkg/app/app.go
package app

import (
	"fmt"
	"mytool3/pkg/config"
	"sync"
)

type Component interface {
	Init(tenantId any) (any, error)
	Name() string
	Close(tenantId any) error
	HealthCheck() bool
}

type App struct {
	globalCfg  *config.AppConfig
	components map[string]Component
	mu         sync.RWMutex
}

var defaultManager = &App{
	components: make(map[string]Component),
}

func GetDefaultApp() *App {
	return defaultManager
}

func (a *App) SetConfig(cfg *config.AppConfig) {
	config.GetConfigManager().SetConfig(cfg)
}

func (a *App) GetGlobalConfig() *config.GlobalConfig {
	return config.GetConfigManager().GetGlobal()
}

func (a *App) TenantConfig(tenantId any) *config.TenantConfig {
	cfg := config.GetConfigManager().GetConfig(tenantId)
	if cfg == nil {
		return nil
	}
	return cfg.Tenant[tenantId]
}

func (a *App) RegisterComponent(c Component) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.components[c.Name()] = c
}

func (a *App) GetComponent(name string) Component {
	a.mu.RLock()
	defer a.mu.RUnlock()
	c, exists := a.components[name]
	if !exists {
		//panic("component not registered: " + name)
		return nil
	}
	return c
}

func (a *App) CloseComponent(name string, tenantId any) error {
	a.mu.RLock()
	c, exists := a.components[name]
	a.mu.RUnlock()

	if !exists {
		return fmt.Errorf("component not registered: %s", name)
	}

	if closer, ok := c.(interface {
		Close(any) error
	}); ok {
		return closer.Close(tenantId)
	}

	return c.Close(tenantId)
}

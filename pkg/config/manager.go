// pkg/config/manager.go
package config

import (
	"sync"
)

var defaultManager = &Manager{
	cfg: &AppConfig{
		Global: &GlobalConfig{},
		Tenant: make(map[any]*TenantConfig),
	},
}

type Manager struct {
	cfg *AppConfig
	mu  sync.RWMutex
}

func GetConfigManager() *Manager {
	return defaultManager
}

func (m *Manager) SetConfig(cfg *AppConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg = cfg
}

func (m *Manager) SetGlobal(cfg *GlobalConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg.Global = cfg
}

func (m *Manager) GetGlobal() *GlobalConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cfg.Global
}

func (m *Manager) AddTenant(id string, cfg *TenantConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg.Tenant[id] = cfg
}

func (m *Manager) GetConfig(tenantId any) *AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return defaultManager.cfg
}

func (m *Manager) SetTenant(id string, cfg *TenantConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg.Tenant[id] = cfg
}

func (m *Manager) GetTenant(id string) *TenantConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cfg.Tenant[id]
}

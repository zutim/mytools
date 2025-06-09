// pkg/app/mysql_component.go
package app

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type MysqlComponent struct {
	instances map[any]*gorm.DB // 租户 ID -> 数据库连接
	mu        sync.Mutex
}

func NewMysqlComponent() Component {
	return &MysqlComponent{
		instances: make(map[any]*gorm.DB),
	}
}

func (m *MysqlComponent) Name() string {
	return "mysql"
}

// mysql_component.go
func (m *MysqlComponent) Init(tenantId any) (any, error) {

	var dsn string
	app := GetDefaultApp()
	globalMysqlConfig := app.GetGlobalConfig().MySQL
	cfg := app.TenantConfig(tenantId).MySQL
	if tenantId == 0 {
		if globalMysqlConfig.Dsn == "" {
			return nil, errors.New(fmt.Sprintf("no mysql dsn config %v", tenantId))
		}

		dsn = globalMysqlConfig.Dsn
	} else {
		if cfg.Dsn == "" {
			return nil, errors.New(fmt.Sprintf("no mysql dsn config %v", tenantId))
		}

		dsn = cfg.Dsn
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 globalMysqlConfig.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(globalMysqlConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(globalMysqlConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(globalMysqlConfig.MaxLife)

	m.mu.Lock()
	defer m.mu.Unlock()
	m.instances[tenantId] = db

	return db, nil
}

func (m *MysqlComponent) Close(tenantId any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, exists := m.instances[tenantId]
	if !exists {
		return fmt.Errorf("no database instance for tenant %v", tenantId)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	delete(m.instances, tenantId)
	return nil
}

func (m *MysqlComponent) HealthCheck() bool {
	return true
}

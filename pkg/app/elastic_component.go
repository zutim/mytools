// pkg/app/elastic_component.go
package app

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"mytool3/pkg/config"
	"sync"
	"time"
)

type ElasticComponent struct {
	instances map[any]*elastic.Client // 租户 ID -> 数据库连接
	mu        sync.Mutex
}

func NewElasticComponent() Component {
	return &ElasticComponent{
		instances: make(map[any]*elastic.Client),
	}
}

func (e *ElasticComponent) Name() string {
	return "elastic"
}

func (e *ElasticComponent) Init(tenantId any) (any, error) {

	var elasticCfg config.ElasticConnConfig
	app := GetDefaultApp()
	global := app.GetGlobalConfig().Elastic
	tenantCfg := app.TenantConfig(tenantId).Elastic

	if tenantId == 0 {
		if global.Hosts == nil {
			return nil, fmt.Errorf("elasticsearch hosts not configured")
		}
		elasticCfg = config.ElasticConnConfig{
			Hosts:    global.Hosts,
			Username: global.Username,
			Password: global.Password,
		}
	} else {
		if tenantCfg.Hosts == nil {
			return nil, fmt.Errorf("elasticsearch hosts not configured for tenant %v", tenantId)
		}
		elasticCfg = config.ElasticConnConfig{
			Hosts:    tenantCfg.Hosts,
			Username: tenantCfg.Username,
			Password: tenantCfg.Password,
		}
	}

	client, err := elastic.NewClient(
		elastic.SetURL(elasticCfg.Hosts...),
		elastic.SetBasicAuth(elasticCfg.Username, elasticCfg.Password),
		elastic.SetSniff(false), // 如果部署在 Docker/K8s 内部网络中，建议关闭 Sniff
		elastic.SetHealthcheckTimeoutStartup(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	e.mu.Lock()
	e.instances[tenantId] = client
	e.mu.Unlock()

	return client, nil
}

func (e *ElasticComponent) Close(tenantId any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, exists := e.instances[tenantId]
	if !exists {
		return fmt.Errorf("no elasticsearch client instance for tenant %v", tenantId)
	}

	delete(e.instances, tenantId)
	return nil
}

func (e *ElasticComponent) HealthCheck() bool {
	// 实现健康检查
	return true
}

// pkg/config/config.go
package config

import (
	"gorm.io/gorm/logger"
	"time"
)

type AppConfig struct {
	Global *GlobalConfig
	Tenant map[any]*TenantConfig
}

// GlobalConfig 全局共享配置
type GlobalConfig struct {
	Log     LogConfig
	MySQL   MySQLGlobalConfig
	Redis   RedisGlobalConfig
	Mongo   MongoGlobalConfig
	Kafka   KafkaGlobalConfig
	Elastic ElasticGlobalConfig
}

// TenantConfig 租户级配置（可覆盖全局）
type TenantConfig struct {
	Log     LogConfig
	MySQL   MySQLTenantConfig
	Redis   RedisTenantConfig
	Mongo   MongoTenantConfig
	Kafka   KafkaTenantConfig
	Elastic ElasticTenantConfig
}

// ------------------ 各组件配置模块 ------------------

type LogConfig struct {
	Level      string
	Path       string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type MySQLConnConfig struct {
	Dsn string
}

type MySQLGlobalConfig struct {
	MySQLConnConfig
	MaxOpen int
	MaxIdle int
	MaxLife time.Duration
	Logger  logger.Interface
}

type MySQLTenantConfig MySQLConnConfig

// 可以提取公共部分作为 interface 或 struct
type RedisCommonConfig struct {
	Addr     string
	Password string
	DB       int
}

type RedisGlobalConfig struct {
	RedisCommonConfig
	MaxOpen int
	MaxIdle int
	MaxLife time.Duration
}

type RedisTenantConfig RedisCommonConfig

// pkg/config/config.go
type MongoConnConfig struct {
	URI string
	DB  string
}

type MongoGlobalConfig struct {
	MongoConnConfig
	MaxPoolSize uint64        // 最大连接池大小
	MinPoolSize uint64        // 最小连接池大小
	Timeout     time.Duration // 连接超时时间
}

type MongoTenantConfig MongoConnConfig

type KafkaConnConfig struct {
	Brokers []string
	GroupID string
	Topics  []string
}

type KafkaGlobalConfig struct {
	KafkaConnConfig
}

type KafkaTenantConfig KafkaConnConfig

// ------------------ Elasticsearch 配置模块 ------------------

type ElasticConnConfig struct {
	Hosts    []string
	Username string
	Password string
}

type ElasticGlobalConfig struct {
	ElasticConnConfig
	MaxRetries int
}

type ElasticTenantConfig ElasticConnConfig

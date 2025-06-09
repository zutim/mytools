// pkg/app/component.go
package app

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/zutim/mongo"
	"github.com/zutim/mytools/pkg/dbmanager"
	"github.com/zutim/mytools/pkg/log2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path/filepath"
)

var (
	GetMysql              = NewComponentAccessor[*gorm.DB]("mysql")
	GetRedis              = NewComponentAccessor[*redis.Client]("redis")
	GetMongo              = NewComponentAccessor[*mongo.MongoClient]("mongo")
	GetKafkaConsumerGroup = NewComponentAccessor[sarama.ConsumerGroup]("kafka-consumer-group")
	GetKafkaProducer      = NewComponentAccessor[sarama.SyncProducer]("kafka-producer")
	GetElastic            = NewComponentAccessor[*elastic.Client]("elastic")
	GetLog                = newLogAccessor()
)

func newLogAccessor() func(options ...func(*log2.LogPathOption)) (*zap.SugaredLogger, error) {
	return func(options ...func(*log2.LogPathOption)) (*zap.SugaredLogger, error) {
		defaultOptions := log2.LogPathOption{
			LogPre:     "./",
			TenantId:   "1",
			ModuleName: "default",
		}

		// 查询是否配置了pre
		app := GetDefaultApp()
		path := app.GetGlobalConfig().Log.Path
		if path != "" {
			defaultOptions.LogPre = path
		}

		for _, setter := range options {
			setter(&defaultOptions)
		}

		// 重新赋值，还是不对，应该根据进来的模块去初始化
		tenantStr := fmt.Sprintf("%v", defaultOptions.TenantId)
		defaultOptions.LogPre = filepath.Join(defaultOptions.LogPre, tenantStr, defaultOptions.ModuleName+".log")

		logger, err := dbmanager.NewDbMap[*zap.SugaredLogger]().WithOptionTenantId(defaultOptions.LogPre, func(path any) (*zap.SugaredLogger, error) {
			obj, err := GetDefaultApp().GetComponent("log").Init(path)
			if err != nil {
				return nil, err
			}
			return obj.(*zap.SugaredLogger), nil
		}, func(logger *zap.SugaredLogger) error {
			return nil
		})

		return logger, err
	}
}

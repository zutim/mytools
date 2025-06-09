// main.go
package main

import (
	"context"
	"gorm.io/gorm/logger"
	"log"
	"mytool3/pkg/app"
	"mytool3/pkg/config"
	"mytool3/pkg/log2"
	"os"
	"time"
)

func main() {
	globalCfg := &config.AppConfig{
		Global: &config.GlobalConfig{
			Log: config.LogConfig{
				Level:      "info",
				Path:       "./logs",
				MaxSize:    100,
				MaxBackups: 10,
				MaxAge:     30,
				Compress:   false,
			},
			MySQL: config.MySQLGlobalConfig{
				MaxOpen: 100,
				MaxIdle: 10,
				MaxLife: 5 * time.Minute,
				Logger: logger.New(
					log.New(os.Stdout, "\r\n", log.LstdFlags),
					logger.Config{
						SlowThreshold:             200 * time.Millisecond,
						IgnoreRecordNotFoundError: true,
						Colorful:                  false,
						LogLevel:                  logger.Warn,
					}),
			},
			Redis: config.RedisGlobalConfig{
				MaxOpen: 50,
				MaxIdle: 10,
				MaxLife: 5 * time.Minute,
			},
			Mongo: config.MongoGlobalConfig{
				MaxPoolSize: 5,
				MinPoolSize: 1,
				Timeout:     5 * time.Second,
			},
			Kafka: config.KafkaGlobalConfig{
				KafkaConnConfig: config.KafkaConnConfig{
					Brokers: []string{"localhost:9092"},
					GroupID: "group-tenant-1",
					Topics:  []string{"topic-tenant-1"},
				},
			},
			Elastic: config.ElasticGlobalConfig{
				MaxRetries: 3,
				ElasticConnConfig: config.ElasticConnConfig{
					Hosts:    []string{""},
					Username: "elastic",
					Password: "password",
				},
			},
		},
		Tenant: map[any]*config.TenantConfig{
			"tenant-1": {
				MySQL: config.MySQLTenantConfig{
					Dsn: "root:root@tcp(127.0.0.1:3306)/geroom?charset=utf8mb4&parseTime=True&loc=Local",
				},
				Redis: config.RedisTenantConfig{
					Addr: "127.0.0.1:6379",
					DB:   1,
				},
				Mongo: config.MongoTenantConfig{
					URI: "",
					DB:  "",
				},
				Kafka: config.KafkaTenantConfig{
					Brokers: []string{"localhost:9092"},
					GroupID: "group-tenant-1",
					Topics:  []string{"topic-tenant-1"},
				},
			},
		},
	}

	apps := app.GetDefaultApp()
	apps.SetConfig(globalCfg)

	// 注册组件
	apps.RegisterComponent(app.NewMysqlComponent())
	apps.RegisterComponent(app.NewLogComponent())
	apps.RegisterComponent(app.NewRedisComponent())
	apps.RegisterComponent(app.NewMongoComponent())
	apps.RegisterComponent(app.NewKafkaConsumerGroupComponent())
	apps.RegisterComponent(app.NewElasticComponent())

	mainLog, err := app.GetLog(log2.WithTenantId("tenant-1"), log2.WithLogPre("./logs"), log2.WithModuleName("main"))
	if err != nil {
		panic(err)
	}

	//type Customer struct {
	//	Id   int
	//	Name string
	//}
	//var customer Customer
	//db, err := app.GetMysql.Get("tenant-1")
	//if err != nil {
	//	mainLog.Error(err)
	//	return
	//}
	//if err := db.Table("tbl_customers").First(&customer).Error; err != nil {
	//	mainLog.Error(err)
	//	return
	//}
	//mainLog.Info(customer.Name)

	//mainLog.Info("hello1 world")
	//redis, err := app.GetRedis.Get("tenant-2")
	//if err != nil {
	//	mainLog.Error(err)
	//	return
	//}
	//redis.Set(context.Background(), "myKey1", "my value 2", 0)
	//
	//res := redis.Get(context.Background(), "myKey1")
	//mainLog.Info(res.Val())

	//mongoClient, err := app.GetMongo.Get(5)
	//if err != nil {
	//	mainLog.Error(err)
	//	return
	//}
	//
	//mongo1 := mongoClient.GetCollection(apps.GetGlobalConfig().Mongo.DB, "sms_logs")
	//cur, err := mongo1.Find(context.Background(), bson.M{
	//	"mobile": "237777222222",
	//})
	//if err != nil {
	//	mainLog.Error(err)
	//	return
	//}
	//
	//type SmsLog struct {
	//	Mobile  string `bson:"mobile"`
	//	Content string `bson:"content"`
	//}
	//var results []*SmsLog
	//_ = cur.All(context.TODO(), &results)
	//for _, result := range results {
	//	mainLog.Info(result.Mobile)
	//}

	client, err := app.GetElastic.Get("tenant-1")
	if err != nil {
		mainLog.Error(err)
		return
	}

	ctx := context.Background()
	exists, err := client.IndexExists("my-index").Do(ctx)
	if err != nil {
		mainLog.Error(err)
		return
	}

	if !exists {
		_, err := client.CreateIndex("my-index").BodyString(`{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"user":    { "type": "keyword" },
				"message": { "type": "text" }
			}
		}
	}`).Do(ctx)
		if err != nil {
			mainLog.Error(err)
			return
		}
	}

	doc := map[string]interface{}{
		"user":    "ztm",
		"message": "Hello from tenant-1",
	}

	_, err = client.Index().
		Index("my-index").
		Type("_doc").
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		mainLog.Error(err)
		return
	}

	mainLog.Info("Document indexed successfully")

}

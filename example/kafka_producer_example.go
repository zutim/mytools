package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"mytool3/pkg/app"
	"mytool3/pkg/config"
)

func main() {
	apps := app.GetDefaultApp()

	// 创建一个配置
	config := &config.AppConfig{
		Global: &config.GlobalConfig{
			Kafka: config.KafkaGlobalConfig{
				KafkaConnConfig: config.KafkaConnConfig{
					Brokers: []string{"localhost:9092"},
					GroupID: "test-group",
					Topics:  []string{"test-topic"},
				},
			},
		},
		Tenant: map[any]*config.TenantConfig{
			"tenant-1": {
				Kafka: config.KafkaTenantConfig{
					Brokers: []string{"localhost:9092"},
					GroupID: "test-group",
					Topics:  []string{"test-topic"},
				},
			},
		},
	}

	apps.SetConfig(config)
	apps.RegisterComponent(app.NewKafkaProducerComponent())

	fmt.Println("12323455")

	kafkaProducer, err := app.GetKafkaProducer.Get("tenant-1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("123")

	msg := &sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.StringEncoder("Hello from tenant-1"),
	}

	fmt.Println("234")
	partition, offset, err := kafkaProducer.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Sent message to partition with offset", partition, offset)
	}

}

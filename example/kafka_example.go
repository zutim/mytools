package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/zutim/mytools/pkg/app"
	"github.com/zutim/mytools/pkg/config"
)

func main() {
	// 创建一个App
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
	apps.RegisterComponent(app.NewKafkaConsumerGroupComponent())

	kafkaCG, err := app.GetKafkaConsumerGroup.Get("tenant-1")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 此时 kafkaCG 就是 sarama.ConsumerGroup 类型
	err = kafkaCG.Consume(context.Background(), []string{"topic-tenant-1"}, &MyConsumerGroupHandler{})
	if err != nil {
		fmt.Println(err)
	}

}

type MyConsumerGroupHandler struct{}

func (h *MyConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *MyConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *MyConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		switch msg.Topic {
		case "topic-tenant-1":

			fmt.Printf("Received message: %s\n", string(msg.Value))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

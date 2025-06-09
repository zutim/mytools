package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
)

type KafkaConsumerGroupComponent struct {
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topics  []string
}

func NewKafkaConsumerGroupComponent() Component {
	return &KafkaConsumerGroupComponent{}
}

func (k *KafkaConsumerGroupComponent) Name() string {
	return "kafka-consumer-group"
}

func (k *KafkaConsumerGroupComponent) Init(tenantId any) (any, error) {
	app := GetDefaultApp()
	gcfg := app.GetGlobalConfig().Kafka
	tcfg := app.TenantConfig(tenantId).Kafka

	var kafkaCfg KafkaConfig
	if tenantId == 0 {
		if gcfg.Brokers == nil {
			return nil, errors.New("no global config")
		}
		kafkaCfg = KafkaConfig{
			Brokers: gcfg.Brokers,
			GroupID: gcfg.GroupID,
			Topics:  gcfg.Topics,
		}
	} else {
		if tcfg.Brokers == nil {
			return nil, errors.New(fmt.Sprintf("no tenant config %v", tenantId))
		}
		kafkaCfg = KafkaConfig{
			Brokers: tcfg.Brokers,
			GroupID: tcfg.GroupID,
			Topics:  tcfg.Topics,
		}
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(kafkaCfg.Brokers, kafkaCfg.GroupID, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			err := consumerGroup.Consume(context.Background(), kafkaCfg.Topics, nil)
			if err != nil {
				panic(err)
			}
		}
	}()

	return consumerGroup, nil
}

func (k *KafkaConsumerGroupComponent) Close(tenantId any) error {
	// 实现关闭逻辑
	return nil
}

func (k *KafkaConsumerGroupComponent) HealthCheck() bool {
	// 实现健康检查逻辑
	return true
}

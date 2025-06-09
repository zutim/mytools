// pkg/app/kafka_producer_component.go
package app

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"time"
)

type KafkaProducerComponent struct{}

func NewKafkaProducerComponent() Component {
	return &KafkaProducerComponent{}
}

func (k *KafkaProducerComponent) Name() string {
	return "kafka-producer"
}

func (k *KafkaProducerComponent) Init(tenantId any) (any, error) {
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
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.DialTimeout = 5 * time.Second
	config.Net.ReadTimeout = 5 * time.Second
	config.Net.WriteTimeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, config)

	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (k *KafkaProducerComponent) Close(tenantId any) error {
	return nil
}

func (k *KafkaProducerComponent) HealthCheck() bool {
	return true
}

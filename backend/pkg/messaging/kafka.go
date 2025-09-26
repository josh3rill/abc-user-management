package messaging

import (
	"context"
    "encoding/json"
    "log"
    "github.com/IBM/sarama"

)

// KafkaProducer handles event publishing to Kafka
type KafkaProducer struct {
    producer sarama.SyncProducer
    topic    string
}

// InitKafka creates Kafka producer as alternative to RabbitMQ
func InitKafka(brokers []string, topic string) (*KafkaProducer, error) {
    // Configure Kafka producer
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    
    // Create producer instance
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }
    
    return &KafkaProducer{
        producer: producer,
        topic:    topic,
    }, nil
}

// PublishEvent sends event to Kafka topic
func (k *KafkaProducer) PublishEvent(eventType string, data interface{}) error {
    // Prepare event payload
    event := map[string]interface{}{
        "type": eventType,
        "data": data,
    }
    
    // Marshal to JSON
    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // Create Kafka message
    msg := &sarama.ProducerMessage{
        Topic: k.topic,
        Key:   sarama.StringEncoder(eventType),
        Value: sarama.ByteEncoder(payload),
    }
    
    // Send message to Kafka
    partition, offset, err := k.producer.SendMessage(msg)
    if err != nil {
        log.Printf("Failed to publish to Kafka: %v", err)
        return err
    }
    
    log.Printf("Event published to Kafka - Partition: %d, Offset: %d", partition, offset)
    
    return nil
}

// Close shuts down Kafka producer
func (k *KafkaProducer) Close() error {
    return k.producer.Close()
}

// KafkaConsumer handles event consumption from Kafka
type KafkaConsumer struct {
    consumer sarama.ConsumerGroup
    topics   []string
    handler  sarama.ConsumerGroupHandler
}

// InitKafkaConsumer creates Kafka consumer for event processing
func InitKafkaConsumer(brokers []string, group string, topics []string) (*KafkaConsumer, error) {
    // Configure consumer group
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    config.Consumer.Offsets.Initial = sarama.OffsetNewest
    
    // Create consumer group
    consumer, err := sarama.NewConsumerGroup(brokers, group, config)
    if err != nil {
        return nil, err
    }
    
    return &KafkaConsumer{
        consumer: consumer,
        topics:   topics,
    }, nil
}

// Start begins consuming messages from Kafka
func (k *KafkaConsumer) Start(handler sarama.ConsumerGroupHandler) error {
    ctx := context.Background()
    
    // Start consuming in goroutine
    go func() {
        for {
            // Handle consumer errors
            if err := k.consumer.Consume(ctx, k.topics, handler); err != nil {
                log.Printf("Kafka consumer error: %v", err)
            }
        }
    }()
    
    log.Println("Kafka consumer started")
    return nil
}
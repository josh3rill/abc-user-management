package events

import (
    "encoding/json"
    "github.com/streadway/amqp"
    "log"
)

type EventPublisher struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

// NewEventPublisher creates RabbitMQ publisher for event-driven architecture
func NewEventPublisher(conn *amqp.Connection) (*EventPublisher, error) {
    if conn == nil {
        return nil, nil // Return nil if no connection
    }
    
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    // Declare exchange for user events
    err = ch.ExchangeDeclare(
        "user_events", // exchange name
        "topic",       // exchange type
        true,          // durable
        false,         // auto-deleted
        false,         // internal
        false,         // no-wait
        nil,           // arguments
    )
    if err != nil {
        return nil, err
    }
    
    return &EventPublisher{
        conn:    conn,
        channel: ch,
    }, nil
}

// Publish sends event to message broker
func (p *EventPublisher) Publish(eventType string, data interface{}) error {
    // Check if publisher is initialized
    if p == nil || p.channel == nil {
        log.Printf("Event publisher not initialized, skipping event: %s", eventType)
        return nil // Don't fail if events are not configured
    }
    
    // Serialize event data to JSON
    body, err := json.Marshal(map[string]interface{}{
        "type": eventType,
        "data": data,
    })
    if err != nil {
        return err
    }
    
    // Publish message to exchange
    err = p.channel.Publish(
        "user_events", // exchange
        eventType,     // routing key
        false,         // mandatory
        false,         // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    
    if err != nil {
        log.Printf("Failed to publish event %s: %v", eventType, err)
        return err
    }
    
    log.Printf("Event published: %s", eventType)
    return nil
}

// Close closes the channel (connection closed separately)
func (p *EventPublisher) Close() error {
    if p != nil && p.channel != nil {
        return p.channel.Close()
    }
    return nil
}
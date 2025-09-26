package messaging

import (
    "fmt"
    "time"
    "github.com/streadway/amqp"
)

// RabbitMQConnection wraps RabbitMQ connection with reconnection logic
type RabbitMQConnection struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    url     string
}

// InitRabbitMQ establishes connection to RabbitMQ server
func InitRabbitMQ(url string) (*amqp.Connection, error) {
    // Retry connection with exponential backoff
    var conn *amqp.Connection
    var err error
    
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        conn, err = amqp.Dial(url)
        if err == nil {
            break
        }
        
        // Exponential backoff
        waitTime := time.Duration(1<<uint(i)) * time.Second
        fmt.Printf("Failed to connect to RabbitMQ (attempt %d/%d), retrying in %s: %v\n", 
            i+1, maxRetries, waitTime, err)
        time.Sleep(waitTime)
    }
    
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %v", maxRetries, err)
    }
    
    fmt.Println("Successfully connected to RabbitMQ")
    return conn, nil
}

// NewRabbitMQConnection creates a managed RabbitMQ connection
func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
    conn, err := InitRabbitMQ(url)
    if err != nil {
        return nil, err
    }
    
    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %v", err)
    }
    
    return &RabbitMQConnection{
        conn:    conn,
        channel: channel,
        url:     url,
    }, nil
}

// DeclareExchange creates an exchange if it doesn't exist
func (r *RabbitMQConnection) DeclareExchange(name, kind string) error {
    return r.channel.ExchangeDeclare(
        name,  // exchange name
        kind,  // exchange type (direct, topic, fanout, headers)
        true,  // durable
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
}

// DeclareQueue creates a queue if it doesn't exist
func (r *RabbitMQConnection) DeclareQueue(name string) (amqp.Queue, error) {
    return r.channel.QueueDeclare(
        name,  // queue name
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
}

// BindQueue binds a queue to an exchange
func (r *RabbitMQConnection) BindQueue(queueName, exchangeName, routingKey string) error {
    return r.channel.QueueBind(
        queueName,    // queue name
        routingKey,   // routing key
        exchangeName, // exchange
        false,        // no-wait
        nil,          // arguments
    )
}

// Publish sends a message to an exchange
func (r *RabbitMQConnection) Publish(exchange, routingKey string, body []byte) error {
    return r.channel.Publish(
        exchange,   // exchange
        routingKey, // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         body,
            DeliveryMode: amqp.Persistent, // Make message persistent
            Timestamp:    time.Now(),
        },
    )
}

// Consume starts consuming messages from a queue
func (r *RabbitMQConnection) Consume(queueName string) (<-chan amqp.Delivery, error) {
    return r.channel.Consume(
        queueName, // queue
        "",        // consumer tag
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // arguments
    )
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQConnection) Close() error {
    if r.channel != nil {
        if err := r.channel.Close(); err != nil {
            return err
        }
    }
    
    if r.conn != nil {
        if err := r.conn.Close(); err != nil {
            return err
        }
    }
    
    return nil
}

// Reconnect attempts to reconnect to RabbitMQ
func (r *RabbitMQConnection) Reconnect() error {
    // Close existing connection
    r.Close()
    
    // Establish new connection
    conn, err := InitRabbitMQ(r.url)
    if err != nil {
        return err
    }
    
    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return err
    }
    
    r.conn = conn
    r.channel = channel
    
    return nil
}

// HealthCheck verifies RabbitMQ connection is alive
func (r *RabbitMQConnection) HealthCheck() error {
    // Try to declare a temporary queue to test connection
    _, err := r.channel.QueueDeclare(
        "",    // random name
        false, // durable
        true,  // delete when unused
        true,  // exclusive
        false, // no-wait
        nil,   // arguments
    )
    
    return err
}
package events

import (
    "encoding/json"
    "log"
    
    "github.com/streadway/amqp"
)

// StartConsumers initializes event consumers for background processing
func StartConsumers(conn *amqp.Connection, db interface{}) error {
    // Create channel for communication
    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    
    // Declare the events exchange
    err = ch.ExchangeDeclare(
        "user_events", // exchange name
        "topic",       // exchange type
        true,          // durable
        false,         // delete when unused
        false,         // exclusive
        false,         // no-wait
        nil,           // arguments
    )
    if err != nil {
        return err
    }
    
    // Declare queue for user events
    q, err := ch.QueueDeclare(
        "user_events_queue", // queue name
        true,                // durable
        false,               // delete when unused
        false,               // exclusive
        false,               // no-wait
        nil,                 // arguments
    )
    if err != nil {
        return err
    }
    
    // Bind queue to exchange for all event types
    err = ch.QueueBind(
        q.Name,        // queue name
        "#",           // routing key (all events)
        "user_events", // exchange
        false,
        nil,
    )
    if err != nil {
        return err
    }
    
    // Start consuming messages
    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer tag
        false,  // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    if err != nil {
        return err
    }
    
    // Initialize event handler registry
    registry := NewEventHandlerRegistry()
    
    // Register handlers for different event types
    registry.Register("user.created", &UserCreatedHandler{})
    registry.Register("user.updated", &UserUpdatedHandler{})
    registry.Register("user.deleted", &UserDeletedHandler{})
    
    // Process messages in goroutine
    go func() {
        for msg := range msgs {
            // Parse event from message body
            var event map[string]interface{}
            if err := json.Unmarshal(msg.Body, &event); err != nil {
                log.Printf("Failed to parse event: %v", err)
                msg.Nack(false, false)
                continue
            }
            
            // Extract event type and data
            eventType, _ := event["type"].(string)
            eventData, _ := event["data"].(map[string]interface{})
            
            log.Printf("Processing event: %s", eventType)
            
            // Route to appropriate handler
            if err := registry.ProcessEvent(eventType, eventData); err != nil {
                log.Printf("Error processing event %s: %v", eventType, err)
                msg.Nack(false, true) // Requeue on error
            } else {
                msg.Ack(false) // Acknowledge successful processing
            }
        }
    }()
    
    log.Println("Event consumers started successfully")
    return nil
}

// ConsumeEvent handles individual event processing
func ConsumeEvent(eventType string, data []byte) error {
    log.Printf("Consuming event: %s with data: %s", eventType, string(data))
    
    // Process based on event type
    switch eventType {
    case "user.created":
        // Handle user creation event
        log.Println("Processing user creation event")
        // Send welcome email, update analytics, etc.
        
    case "user.updated":
        // Handle user update event
        log.Println("Processing user update event")
        // Invalidate caches, sync with other services
        
    case "user.deleted":
        // Handle user deletion event
        log.Println("Processing user deletion event")
        // Clean up related data, archive records
        
    default:
        log.Printf("Unknown event type: %s", eventType)
    }
    
    return nil
}
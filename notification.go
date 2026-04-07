package main

import (
	"log"
	"context"
	"time"
)

func initMS() (*Broker) {
	broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    FailOnError(err, "Failed to connect to RabbitMQ")

	err = broker.DeclareExchange("Exchange")
	FailOnError(err, "Failed to create Exchange")

	_, err = broker.DeclareQueue("Notification")
	FailOnError(err, "Failed to create Notification queue")
	
	broker.BindQueue("Notification", "promocao.publicada", "Exchange")
	return broker
}

func main() {
	broker := initMS()

    defer broker.Conn.Close()
    defer broker.Ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	msgs, err := broker.Consume("Notification")
    if err != nil {
        return 
    }

    go func() {
        for msg := range msgs {
            log.Printf("Received: %s", "promocao." + string(msg.Body) + " foi publicada")
            broker.Publish(ctx, "Exchange",  "promocao." + string(msg.Body), string(msg.Body))
        }
    }()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")

    select {} 
}

package main

import (
    "context"
    "fmt"
    "log"
    "time"
)

func main() {
    broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")

    defer broker.Conn.Close()
    defer broker.Ch.Close()

    q, err := broker.DeclareQueue("hello")
    failOnError(err, "Failed to declare queue")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    for i := 0; i < 10; i++ {
        msg := fmt.Sprintf("Message %d", i)

        err := broker.Publish(ctx, q.Name, msg)
        failOnError(err, "Failed to publish message")

        log.Printf(" [x] Sent %s\n", msg)
    }
}
package main

import (
	"context"
	"log"
    amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
    Conn *amqp.Connection
    Ch   *amqp.Channel
}

func NewBroker(url string) (*Broker, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    return &Broker{
        Conn: conn,
        Ch:   ch,
    }, nil
}

func (b *Broker) DeclareQueue(name string) (amqp.Queue, error) {
    return b.Ch.QueueDeclare(
        name,
        true,
        false,
        false,
        false,
        amqp.Table{
            amqp.QueueTypeArg: amqp.QueueTypeQuorum,
        },
    )
}


func (b *Broker) Publish(ctx context.Context, queueName string, msg string) error {
    return b.Ch.PublishWithContext(
        ctx,
        "",
        queueName,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(msg),
        },
    )
}

func (b *Broker) Consume(queueName string) (<-chan amqp.Delivery, error) {
    return b.Ch.Consume(
        queueName,
        "",    
        true,  
        false, 
        false, 
        false, 
        nil,   
    )
}

func (b *Broker) StartConsumer(queueName string) error {
    q, err := b.DeclareQueue(queueName)
    if err != nil {
        return err
    }

    msgs, err := b.Consume(q.Name)
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            log.Printf("Received: %s", msg.Body)
        }
    }()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")

    select {} 
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}
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

func (b *Broker) DeclareExchange(name string) error {
    err := b.Ch.ExchangeDeclare(
        name, 
        "direct",      
        true,          
        false,         
        false,         
        false,         
        nil,           
    )
    return err 
}

func (b *Broker) DeclareQueue(name string) (amqp.Queue, error) {
    q, err := b.Ch.QueueDeclare(
        name,  
        false,  
        true, 
        false, 
        false, 
        nil,   
    )
    return q, err
}


func (b *Broker) BindQueue(queueName string, routingKey string, exchangeName string) error {
    return b.Ch.QueueBind(
        queueName,    
        routingKey,     
        exchangeName, 
        false,           
        nil,             
    )
}

func (b *Broker) Publish(ctx context.Context, exchangeName string, routingKey string, msg string) error {
    return b.Ch.PublishWithContext(
        ctx,
        exchangeName, 
        routingKey,   
        false,        
        false,        
        amqp.Publishing{
            ContentType: "application/json",
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

func FailOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}

package main

import (
	"log"
	"context"
	"time"
    "strings"
    "fmt"
)

type Score struct {
    Values  map[string]int
}

func NewScore() *Score {
	return &Score{
		Values: make(map[string]int),
	}
}

func (m *Score) Add(key string, value int) (int){
    m.Values[key] = value
    return value
}

func initMS() (*Broker) {
	broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    FailOnError(err, "Failed to connect to RabbitMQ")

	err = broker.DeclareExchange("Exchange")
	FailOnError(err, "Failed to create Exchange")

	_, err = broker.DeclareQueue("Ranking")
	FailOnError(err, "Failed to create Notification queue")
	
	broker.BindQueue("Ranking", "promocao.voto", "Exchange")
	return broker
}

func processVote(msg string, broker *Broker, score *Score) (int){
    parts := strings.Split(msg, " ")
    
    fmt.Print(parts)

    var vote int 
    if (parts[1] == "1"){
        vote = 1
    }else{
        vote = -1
    }
    
    return score.Add(parts[0], score.Values[parts[0]] + vote)
}

func main() {
	broker := initMS()
	score := NewScore()

    defer broker.Conn.Close()
    defer broker.Ch.Close()

    signer, _ := NewSigner()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	msgs, err := broker.Consume("Ranking")
    if err != nil {
        return 
    }

    go func() {
        for msg := range msgs {
            content, err := signer.Open(string(msg.Body)) 
			if err == nil{
                if(processVote(content, broker, score) >= 3){ 
                    broker.Publish(ctx, "Exchange",  "promocao.destaque", signer.Sign(strings.Split(content, " ")[0] + " hot deal"))
                }   
            }
        }
    }()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")

    select {} 
}

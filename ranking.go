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
    LargestKey string
    MaxVal   int
}

func NewScore() *Score {
	return &Score{
		Values: make(map[string]int),
	}
}

func (m *Score) Add(key string, value int) {
    m.Values[key] = value
    if (value > m.MaxVal || key == m.LargestKey) {
        m.MaxVal = value
        m.LargestKey = key
    }
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

func processVote(msg string, broker *Broker, score *Score){
    parts := strings.Split(msg, " ")
    
    var vote int 
    if (parts[1] == "positivo"){
        vote = 1
    }else{
        vote = -1
    }
    
    score.Add(parts[0], score.Values[parts[0]] + vote)
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
                processVote(content, broker, score)
                fmt.Printf("Max Val: %d\n", score.MaxVal)
                if(score.MaxVal >= 3){ 
                    broker.Publish(ctx, "Exchange",  "promocao.destaque", signer.Sign(content + " hot deal"))
                }   
            }
        }
    }()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")

    select {} 
}

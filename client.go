package main

import (
	"log"
	"fmt"
	"github.com/google/uuid"
)

func initClient(categorias [] string) (*Broker, string) {
	broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    FailOnError(err, "Failed to connect to RabbitMQ")

	err = broker.DeclareExchange("Exchange")
	FailOnError(err, "Failed to create Exchange")

	id := uuid.New()
	_, err = broker.DeclareQueue(id.String())
	FailOnError(err, "Failed to create client" + id.String() + "queue")
	
	for i := 0; i < len(categorias); i++ {
		broker.BindQueue(id.String(), "promocao." + categorias[i], "Exchange")
	}

	return broker, id.String()
}

func menu() ([] string) {
	var categorias []string
	var input string

	fmt.Println("Digite as categorias de promoções uma por uma (ou digite 'fim' para encerrar):")
	fmt.Println("ex: livros")

	for {
		fmt.Printf("CATEGORIA %d: ", len(categorias)+1)
		fmt.Scan(&input)

		if input == "fim" {
			break
		}

		categorias = append(categorias, input)
	}

	fmt.Printf("\nExecutando com as categorias: %v\n", categorias)
	return categorias
}

func main() {
	categorias := menu()
	broker, id := initClient(categorias)

    defer broker.Conn.Close()
    defer broker.Ch.Close()

	msgs, err := broker.Consume(id)
    if err != nil {
        return 
    }

    go func() {
        for msg := range msgs {
            log.Printf("Received: %s", msg.Body)
        }
    }()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")

    select {} 
}

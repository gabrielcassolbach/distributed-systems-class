package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"context"
	"time"
)

func initMS() (*Broker) {
	broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    FailOnError(err, "Failed to connect to RabbitMQ")

	err = broker.DeclareExchange("Exchange")
	FailOnError(err, "Failed to create Exchange")

	_, err = broker.DeclareQueue("Gateway")
	FailOnError(err, "Failed to create Gateway queue")
	
	broker.BindQueue("Gateway", "promocao.publicada", "Exchange")
	return broker
}

func main() {
	broker := initMS()
	defer broker.Conn.Close()
    defer broker.Ch.Close()
	
	reader := bufio.NewReader(os.Stdin)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	broker.StartConsuming("Gateway")

	for {
		fmt.Println("\n--- MENU PRINCIPAL ---")
		fmt.Println("1. Cadastrar")
		fmt.Println("2. Votar")
		fmt.Println("3. Sair")
		fmt.Print("Escolha uma opção: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			broker.Publish(ctx, "Exchange", "promocao.recebida", "promocao")
		case "2":
			broker.Publish(ctx, "Exchange", "promocao.voto", "promocaovoto")
		case "3":
			return
		default:
			fmt.Println("Opção inválida, tente novamente.")
		}
	}

}

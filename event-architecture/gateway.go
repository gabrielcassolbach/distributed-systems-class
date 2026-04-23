package main

import (
	"bufio"
	"fmt"
	"os"
	"context"
	"time"
	"sync"
	"log"
	"runtime"
	"strings"
	"os/exec"
)

var (
    historico []string
    mu        sync.Mutex 
)

func inHistory(promotion string) (bool) {
	for i := 0; i < len(historico); i++ {
		if(historico[i] == promotion){
			return true
		}
	}
	return false 
}

func menu(reader *bufio.Reader) (string, string, string) {
	CallClear()
	fmt.Println("--- MENU PRINCIPAL ---")
	fmt.Println("1. Cadastrar")
	fmt.Println("2. Votar")
	fmt.Println("3. Listar")
	fmt.Println("4. Sair")
	fmt.Print("Escolha uma opção: ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	if option == "4" {
		return option, "", ""
	}

	var promotion string
	if (option == "1" || option == "2") {
		CallClear()
		fmt.Println("--- Digite o nome da promoção ---")
		promotion, _ = reader.ReadString('\n')
		promotion = strings.TrimSpace(promotion)
		if(option == "1"){
			fmt.Printf("Promoção '%s' cadastrada!", promotion)
		}
		if(option == "2" && !inHistory(promotion)) {
			option = "NULLPTR"
			fmt.Printf("Promoção não existe (Pressione Enter)\n")
			reader.ReadString('\n') 
		}
	}

	var vote string
	vote = "undefined"
	if (option == "2"){
		for vote != "positivo" && vote != "negativo" {
			CallClear()
			fmt.Println("Digite positivo ou negativo:")
			vote, _ = reader.ReadString('\n')
			vote = strings.TrimSpace(vote)
		}
	}

	return option, promotion, vote
}


func CallClear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}


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

func StartConsuming(b *Broker, s *Signer){
	msgs, err := b.Consume("Gateway")
    if err != nil {
        return  
    }

    go func() {
        for msg := range msgs {
			content, err := s.Open(string(msg.Body)) 
			if err == nil{
				mu.Lock()
				historico = append(historico, content)
				mu.Unlock()
				log.Printf("\npromoção %s publicada\n", content)
			}
        }
    }()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C\n")
}

func listPromotions(){
	fmt.Println("\n --- Histórico de Promoções ---")
	for i := 0; i < len(historico); i++ {
		fmt.Println(historico[i])
	}
}

func main() {
	broker := initMS()
	defer broker.Conn.Close()
    defer broker.Ch.Close()
	
	signer, _ := NewSigner()

	reader := bufio.NewReader(os.Stdin)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	StartConsuming(broker, signer)

	for {
		option, promotion, vote := menu(reader)
		switch option {
		case "1":
			broker.Publish(ctx, "Exchange", "promocao.recebida", signer.Sign(promotion))
		case "2":
			promotion = promotion + " " + vote
			broker.Publish(ctx, "Exchange", "promocao.voto", signer.Sign(promotion))
		case "3":
			listPromotions()
			fmt.Println("\n\nPress any key to go back to menu")
			reader.ReadString('\n') 
		case "4":
			return
		default:
			fmt.Println("Opção inválida, tente novamente.")
		}
	}

}

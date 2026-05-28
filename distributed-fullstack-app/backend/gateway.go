package main

import (
	"log"
	"net/http"
	"sync"
)

var (
	historico    []string
	mu           sync.Mutex
)

func inHistory(promotion string) bool {
	for i := 0; i < len(historico); i++ {
		if historico[i] == promotion {
			return true
		}
	}
	return false
}

func StartConsuming(b *Broker, s *Signer) {
	msgs, err := b.Consume("Gateway")
	if err != nil {
		return
	}

	go func() {
		for msg := range msgs {
			content, err := s.Open(string(msg.Body))
			if err == nil {
				mu.Lock()
				historico = append(historico, content)
				mu.Unlock()
				log.Printf("\npromoção %s publicada\n", content)
			} else {
                log.Printf("[ERRO] Falha ao abrir/validar assinatura: %v", err)
            }
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C\n")
}

func initMS() *Broker {
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

	signer, _ := NewSigner()
	StartConsuming(broker, signer)

	http.HandleFunc("/api/promotions/register", registerPromotions(broker, signer))
	http.HandleFunc("/api/promotions/list", listPromotions(&historico, &mu))
	http.HandleFunc("/api/promotions/vote", voteInPromotion(broker, signer))
	http.HandleFunc("/api/categories/subscribe", registerInterest(broker))
	http.HandleFunc("/api/categories/unsubscribe", cancelInterest(broker))

	http.HandleFunc("/api/sse", sseHandler(broker))

	log.Println("Gateway server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

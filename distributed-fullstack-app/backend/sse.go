package main

import (
	"log"
	"net/http"
	"fmt"
)

func sseHandler(broker *Broker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

        id := r.URL.Query().Get("client_id")
        if id == "" {
            http.Error(w, "Missing client_id", http.StatusBadRequest)
            return
        }

        _, err := broker.DeclareQueue(id)
        if err != nil {
            log.Printf("Failed to create client queue: %v", err)
            return
        }

        broker.BindQueue(id, "promocao.destaque", "Exchange")

        msgs, err := broker.Consume(id)
        if err != nil {
            log.Printf("Failed to consume from queue: %v", err)
            return
        }

        for msg := range msgs {
            fmt.Fprintf(w, "data: %s\n\n", string(msg.Body))
            
            if flusher, ok := w.(http.Flusher); ok {
                flusher.Flush()
            }
        }
    }
}
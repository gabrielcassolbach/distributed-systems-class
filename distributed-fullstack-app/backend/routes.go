package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sync"
	"log"
)

type RequestPayload struct {
    Payload  string `json:"payload"`   
    ClientID string `json:"client_id"` 
}

func verifyRoute(w http.ResponseWriter, r *http.Request, method string) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return false
	}
	
	if r.Method != method {
		http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return false
	}
	return true
}

func publishMessage(b *Broker, routingKey string, payload string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return b.Publish(ctx, "Exchange", routingKey, payload)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := fmt.Sprintf(`{"status":"%s"}`, message)
	w.Write([]byte(response))
}

func registerPromotions(b *Broker, s *Signer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !verifyRoute(w, r, http.MethodPost) {
			return
		}

		var req RequestPayload
		if !decodeJSONBody(w, r, &req) {
			return
		}

		signedPromotion := s.Sign(req.Payload)

		if err := publishMessage(b, "promocao.recebida", signedPromotion); err != nil {
			http.Error(w, "Failed to publish promotion to broker", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, http.StatusAccepted, "promotion_published_successfully")
	}
}

func voteInPromotion(b *Broker, s *Signer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !verifyRoute(w, r, http.MethodPost) {
			return
		}

		var req RequestPayload
		if !decodeJSONBody(w, r, &req) {
			return
		}

		signedVote := s.Sign(req.Payload)

		if err := publishMessage(b, "promocao.voto", signedVote); err != nil {
			http.Error(w, "Failed to publish vote to broker", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, http.StatusOK, "vote_submitted_successfully")
	}
}

func listPromotions(historico *[]string, mu *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !verifyRoute(w, r, http.MethodGet) {
			return
		}

		mu.Lock()
		responseData, err := json.Marshal(*historico)
		fmt.Println("Histórico atual enviado para o front:", *historico)
		mu.Unlock()

		if err != nil {
			http.Error(w, "Internal server error parsing history", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}
}

func registerInterest(b *Broker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !verifyRoute(w, r, http.MethodPost) {
            return
        }

        var req RequestPayload
        if !decodeJSONBody(w, r, &req) {
            return
        }

        if req.ClientID == "" {
            http.Error(w, "Missing client_id (UUID)", http.StatusBadRequest)
            return
        }

        routingKey := fmt.Sprintf("promocao.%s", req.Payload)

        err := b.Ch.QueueBind(
            req.ClientID, 
            routingKey,  
            "Exchange",   
            false,        
            nil,          
        )
        if err != nil {
            log.Printf("Erro ao vincular fila %s à chave %s: %v", req.ClientID, routingKey, err)
            http.Error(w, "Failed to subscribe to category", http.StatusInternalServerError)
            return
        }

        sendJSONResponse(w, http.StatusOK, "subscribed_successfully")
    }
}

func cancelInterest(b *Broker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !verifyRoute(w, r, http.MethodPost) {
            return
        }

        var req RequestPayload
        if !decodeJSONBody(w, r, &req) {
            return
        }

        if req.ClientID == "" {
            http.Error(w, "Missing client_id (UUID)", http.StatusBadRequest)
            return
        }

        routingKey := fmt.Sprintf("promocao.%s", req.Payload)

        err := b.Ch.QueueUnbind(
            req.ClientID, 
            routingKey,   
            "Exchange",   
            nil,          
        )

        if err != nil {
            log.Printf("Erro ao desvincular fila %s da chave %s: %v", req.ClientID, routingKey, err)
            http.Error(w, "Failed to unsubscribe from category", http.StatusInternalServerError)
            return
        }

        sendJSONResponse(w, http.StatusOK, "unsubscribed_successfully")
    }
}
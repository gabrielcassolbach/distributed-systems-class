package main

import (
    "log"
	"fmt"
    "net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

var UsersMap = make(map[string]*User)

func main() {
    broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")

    defer broker.Conn.Close()
    defer broker.Ch.Close()

	registerEndpoints(broker)
	startUsersPrinter()

	http.Handle("/", http.FileServer(http.Dir("../frontend")))

    log.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerEndpoints(broker *Broker){
	createUserQueue(broker)
}

func registerUser(broker *Broker) (string) {
	var userID = uuid.New().String()
    UsersMap[userID] = NewUser(userID, []string{}) 
	return userID
}

func createUserQueue(broker *Broker) {
    http.HandleFunc("/createQueue", func(w http.ResponseWriter, r *http.Request) {
        var userID = registerUser(broker)

        queue, err := broker.DeclareQueue("user-" + userID + "-queue")
        if err != nil {
            http.Error(w, "failed to declare queue", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "userID":    userID,
            "queueName": queue.Name,
        })
    })
}

func startUsersPrinter() {
    go func() {
        for {
            fmt.Println("Current UsersMap:")
            for uuid, user := range UsersMap {
                fmt.Printf("UUID: %s, Bindings: %v\n", uuid, user.Bindings)
            }
            fmt.Println("------")
            time.Sleep(5 * time.Second)
        }
    }()
}
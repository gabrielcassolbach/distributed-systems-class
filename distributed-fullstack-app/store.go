package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"os/exec"
)

type RequestPayload struct {
	Payload string `json:"payload"`
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		clearScreen()
		fmt.Println("=== Promotion Register CLI (STORE) ===")
		fmt.Println("Type the promotion name and press ENTER (or type 'end' to leave):")
		fmt.Print("\n> ")
		
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if strings.ToLower(text) == "end" {
			fmt.Println("leaving...")
			break
		}

		if strings.TrimSpace(text) == "" {
			continue
		}

		sendPromotion(text)
	}
}

func sendPromotion(text string) {
	jsonData, _ := json.Marshal(RequestPayload{Payload: text})

	resp, err := http.Post("http://localhost:8080/api/promotions/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error connecting to the server: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		fmt.Println("Promotion Registered!")
	} else {
		fmt.Printf("Failure. Status HTTP: %d\n", resp.StatusCode)
	}
}
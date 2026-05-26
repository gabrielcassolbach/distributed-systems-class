package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	log.Println("Frontend rodando em http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", fs))
}
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	clientpkg "github.com/paraizofelipe/fullcycle-client-server/internal/client"
)

const (
	serverURL     = "http://localhost:8080/exchange_rates"
	clientTimeout = 300 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	httpClient := &http.Client{}
	bid, err := clientpkg.FetchBid(ctx, httpClient, serverURL)
	if err != nil {
		log.Printf("request error: %v", err)
		log.Fatal("failed to fetch exchange rate")
	}
	log.Printf("exchange rate received: bid=%s", bid)

	content := []byte("Dólar: " + bid)
	if err := os.WriteFile("exchange_rates.txt", content, 0644); err != nil {
		log.Fatalf("write file: %v", err)
	}
}

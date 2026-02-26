package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/paraizofelipe/fullcycle/client-server/internal/exchange_rates"
	"github.com/paraizofelipe/fullcycle/client-server/internal/storage"
	_ "modernc.org/sqlite"
)

const (
	serverAddr = ":8080"
	apiTimeout = 200 * time.Millisecond
	dbTimeout  = 10 * time.Millisecond
)

type exchangeRateResponse struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite", "exchange_rates.db")
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("close database: %v", closeErr)
		}
	}()

	store := storage.New(db)
	if err := store.Init(context.Background(), dbTimeout); err != nil {
		log.Fatalf("init database: %v", err)
	}

	httpClient := &http.Client{}
	server := &http.Server{
		Addr:              serverAddr,
		Handler:           routes(store, httpClient),
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Printf("listening on %s", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func routes(store *storage.Store, httpClient *http.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request received: %s %s", r.Method, r.URL.Path)
		if r.Method != http.MethodGet {
			log.Printf("request error: %s %s method not allowed", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		bid, err := exchange_rates.FetchBid(r.Context(), httpClient, exchange_rates.DefaultQuoteURL, apiTimeout)
		if err != nil {
			log.Printf("request error: fetch bid failed: %v", err)
			http.Error(w, "failed to fetch exchange rate", http.StatusBadGateway)
			return
		}
		log.Printf("exchange rate fetched: bid=%s", bid)

		if err := store.SaveBid(r.Context(), bid, dbTimeout); err != nil {
			log.Printf("request error: save bid failed: %v", err)
			http.Error(w, "failed to persist exchange rate", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(exchangeRateResponse{Bid: bid}); err != nil {
			log.Printf("request error: encode response failed: %v", err)
			return
		}
		log.Printf("request success: %s %s", r.Method, r.URL.Path)
	})
	return mux
}

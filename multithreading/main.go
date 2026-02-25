package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type BrasilAPI struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCEP struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

type Result struct {
	API  string
	Data interface{}
}

func fetchBrasilAPI(ctx context.Context, cep string, ch chan<- Result) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data BrasilAPI
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}
	ch <- Result{API: "BrasilAPI", Data: data}
}

func fetchViaCEP(ctx context.Context, cep string, ch chan<- Result) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data ViaCEP
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}
	ch <- Result{API: "ViaCEP", Data: data}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: cep-race <cep>")
		os.Exit(1)
	}
	cep := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan Result, 2)

	go fetchBrasilAPI(ctx, cep, ch)
	go fetchViaCEP(ctx, cep, ch)

	select {
	case result := <-ch:
		fmt.Printf("API: %s\n", result.API)
		switch d := result.Data.(type) {
		case BrasilAPI:
			fmt.Printf("CEP:    %s\n", d.CEP)
			fmt.Printf("Rua:    %s\n", d.Street)
			fmt.Printf("Bairro: %s\n", d.Neighborhood)
			fmt.Printf("Cidade: %s\n", d.City)
			fmt.Printf("Estado: %s\n", d.State)
		case ViaCEP:
			fmt.Printf("CEP:    %s\n", d.CEP)
			fmt.Printf("Rua:    %s\n", d.Logradouro)
			fmt.Printf("Bairro: %s\n", d.Bairro)
			fmt.Printf("Cidade: %s\n", d.Localidade)
			fmt.Printf("Estado: %s\n", d.UF)
		}
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "error: timeout - no API responded within 1 second")
		os.Exit(1)
	}
}

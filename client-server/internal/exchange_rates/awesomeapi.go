package exchange_rates

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/paraizofelipe/fullcycle/client-server/internal/ctxlog"
)

const DefaultQuoteURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type awesomeAPIResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func FetchBid(parent context.Context, client *http.Client, url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		ctxlog.LogDeadline(ctx, err, "exchange rate api")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unexpected exchange rate api status")
	}

	var payload awesomeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	if payload.USDBRL.Bid == "" {
		return "", errors.New("empty bid")
	}

	return payload.USDBRL.Bid, nil
}

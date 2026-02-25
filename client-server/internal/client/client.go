package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/paraizofelipe/fullcycle/client-server/internal/ctxlog"
)

type exchangeRateResponse struct {
	Bid string `json:"bid"`
}

func FetchBid(ctx context.Context, httpClient *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		ctxlog.LogDeadline(ctx, err, "client request")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unexpected server status")
	}

	var payload exchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	if payload.Bid == "" {
		return "", errors.New("empty bid in response")
	}

	return payload.Bid, nil
}

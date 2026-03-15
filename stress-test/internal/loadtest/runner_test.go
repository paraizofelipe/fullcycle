package loadtest

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRun_InvalidInputs(t *testing.T) {
	cases := []struct {
		name        string
		requests    int
		concurrency int
	}{
		{"zero requests", 0, 1},
		{"negative requests", -1, 1},
		{"zero concurrency", 1, 0},
		{"negative concurrency", 1, -1},
		{"both zero", 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			report := Run("http://example.com", tc.requests, tc.concurrency)
			if report.TotalSent != 0 {
				t.Errorf("expected TotalSent=0, got %d", report.TotalSent)
			}
			if report.Errors != 0 {
				t.Errorf("expected Errors=0, got %d", report.Errors)
			}
			if len(report.StatusCounts) != 0 {
				t.Errorf("expected empty StatusCounts, got %v", report.StatusCounts)
			}
		})
	}
}

func TestRun_AllHTTP200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	report := Run(srv.URL, 10, 3)

	if report.TotalSent != 10 {
		t.Errorf("expected TotalSent=10, got %d", report.TotalSent)
	}
	if report.StatusCounts[200] != 10 {
		t.Errorf("expected 10 HTTP 200, got %d", report.StatusCounts[200])
	}
	if report.Errors != 0 {
		t.Errorf("expected no errors, got %d", report.Errors)
	}
}

func TestRun_MixedStatusCodes(t *testing.T) {
	var count atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if count.Add(1)%2 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	report := Run(srv.URL, 10, 2)

	if report.TotalSent != 10 {
		t.Errorf("expected TotalSent=10, got %d", report.TotalSent)
	}
	total := report.StatusCounts[200] + report.StatusCounts[500]
	if total != 10 {
		t.Errorf("expected 10 total responses, got %d", total)
	}
}

func TestRun_ConnectionError(t *testing.T) {
	report := Run("http://127.0.0.1:1", 5, 2)

	if report.TotalSent != 5 {
		t.Errorf("expected TotalSent=5, got %d", report.TotalSent)
	}
	if report.Errors != 5 {
		t.Errorf("expected 5 errors, got %d", report.Errors)
	}
}

func TestRun_LargeRequestCount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	const large = maxChannelBuffer + 100
	report := Run(srv.URL, large, 2)

	if report.TotalSent != large {
		t.Errorf("expected TotalSent=%d, got %d", large, report.TotalSent)
	}
	if report.StatusCounts[200] != large {
		t.Errorf("expected %d HTTP 200, got %d (errors=%d, other=%v)",
			large, report.StatusCounts[200], report.Errors, report.StatusCounts)
	}
}

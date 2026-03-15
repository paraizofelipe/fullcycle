package loadtest

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

const maxChannelBuffer = 10_000

type result struct {
	statusCode int
	err        error
}

type Report struct {
	TotalTime    time.Duration
	TotalSent    int
	StatusCounts map[int]int
	Errors       int
}

func (r Report) Print() {
	fmt.Println("========== Relatório ==========")
	fmt.Printf("Tempo total:            %s\n", r.TotalTime.Round(time.Millisecond))
	fmt.Printf("Total de requisições:   %d\n", r.TotalSent)
	fmt.Printf("Respostas HTTP 200:     %d\n", r.StatusCounts[200])

	hasOtherCodes := false
	for code := range r.StatusCounts {
		if code != 200 {
			hasOtherCodes = true
			break
		}
	}

	if hasOtherCodes || r.Errors > 0 {
		fmt.Println("\nDistribuição de outros status:")

		codes := make([]int, 0, len(r.StatusCounts))
		for code := range r.StatusCounts {
			if code != 200 {
				codes = append(codes, code)
			}
		}
		sort.Ints(codes)

		for _, code := range codes {
			fmt.Printf("  HTTP %d: %d\n", code, r.StatusCounts[code])
		}
		if r.Errors > 0 {
			fmt.Printf("  Erros de conexão: %d\n", r.Errors)
		}
	}
	fmt.Println("================================")
}

func Run(url string, requests, concurrency int) Report {
	if requests <= 0 || concurrency <= 0 {
		return Report{StatusCounts: make(map[int]int)}
	}

	bufSize := min(requests, maxChannelBuffer)
	jobs := make(chan struct{}, bufSize)
	results := make(chan result, bufSize)

	go func() {
		for i := 0; i < requests; i++ {
			jobs <- struct{}{}
		}
		close(jobs)
	}()

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		completed int
	)

	printProgress := func() {
		mu.Lock()
		completed++
		done := completed
		mu.Unlock()
		if done%100 == 0 || done == requests {
			fmt.Printf("\rProgresso: %d/%d", done, requests)
		}
	}

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 30 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return nil
				},
			}
			for range jobs {
				resp, err := client.Get(url)
				printProgress()
				if err != nil {
					results <- result{err: err}
					continue
				}
				_, _ = io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				results <- result{statusCode: resp.StatusCode}
			}
		}()
	}

	report := Report{
		TotalSent:    requests,
		StatusCounts: make(map[int]int),
	}

	var collectWg sync.WaitGroup
	collectWg.Add(1)
	go func() {
		defer collectWg.Done()
		for r := range results {
			if r.err != nil {
				report.Errors++
			} else {
				report.StatusCounts[r.statusCode]++
			}
		}
	}()

	wg.Wait()
	fmt.Println()
	report.TotalTime = time.Since(start)
	close(results)
	collectWg.Wait()

	return report
}

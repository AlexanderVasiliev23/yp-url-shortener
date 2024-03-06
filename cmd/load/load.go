// Генерация нагрузки для профилирования с использованием pprof
package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	url = "http://localhost:8080/"
)

// + GET /api/user/urls
// + POST /api/shorten
// + POST /api/shorten/batch
// + POST /
// + DELETE /api/user/urls
// + GET /{token}

func main() {
	rate := vegeta.Rate{Freq: 500, Per: time.Second}
	duration := 20 * time.Second

	targets := make([]vegeta.Target, 0, 10_000)

	header := map[string][]string{
		"Cookie": {"jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjE3MDk3MTExNjE3NjY0MzY1NTF9.NsC8Xj88NwyaOou6p1sbkzchdzbgaEZOYDrUDfFh7HE"},
	}

	// + DELETE /api/user/urls
	for i := 0; i < 500; i++ {
		target := vegeta.Target{
			Method: http.MethodDelete,
			URL:    url + "api/user/urls",
			Body:   []byte(`["iWQADUUloU","EuooaEIvNI","RlbUtAKMvS"]`),
			Header: header,
		}
		targets = append(targets, target)
	}

	// + POST /api/shorten/batch
	for i := 0; i < 500; i++ {
		target := vegeta.Target{
			Method: http.MethodPost,
			URL:    url + "api/shorten/batch",
			Body:   []byte(fmt.Sprintf(`[{"correlation_id":"%s","original_url":"%s"}]`, gofakeit.UUID(), gofakeit.URL())),
			Header: header,
		}
		targets = append(targets, target)
	}

	// + POST /api/shorten
	for i := 0; i < 500; i++ {
		target := vegeta.Target{
			Method: http.MethodPost,
			URL:    url + "api/shorten",
			Body:   []byte(fmt.Sprintf(`{"url":"%s"}`, gofakeit.URL())),
			Header: header,
		}
		targets = append(targets, target)
	}

	// + GET /api/user/urls
	for i := 0; i < 500; i++ {
		target := vegeta.Target{
			Method: http.MethodGet,
			URL:    url + "api/user/urls",
			Header: header,
		}
		targets = append(targets, target)
	}

	// + POST /
	for i := 0; i < 500; i++ {
		target := vegeta.Target{
			Method: http.MethodPost,
			URL:    url,
			Body:   []byte(gofakeit.URL()),
			Header: header,
		}
		targets = append(targets, target)
	}

	// + GET /{token}
	for i := 0; i < 500; i++ {
		target := vegeta.Target{
			Method: http.MethodGet,
			URL:    url + gofakeit.LoremIpsumSentence(1),
			Header: header,
		}
		targets = append(targets, target)
	}

	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

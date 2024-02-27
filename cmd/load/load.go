package main

import (
	"log"
	"net/http"
	"strings"
)

const (
	url   = "http://localhost:8080"
	times = 300
)

func main() {
	for i := 0; i < times; i++ {
		if err := generateLoad(); err != nil {
			log.Panicf("generate load: %v", err)
		}
	}
}

func generateLoad() error {
	client := http.Client{}
	for _, req := range buildRequests() {
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	return nil
}

func buildRequests() []*http.Request {
	var requests []*http.Request

	r, _ := http.NewRequest(http.MethodGet, url+"/token", nil)
	requests = append(requests, r)

	r, _ = http.NewRequest(http.MethodGet, url+"/ping", nil)
	requests = append(requests, r)

	r, _ = http.NewRequest(http.MethodPost, url, strings.NewReader("http://test.me"))
	requests = append(requests, r)

	r, _ = http.NewRequest(http.MethodPost, url+"/api/shorten", strings.NewReader(`{"url":"http://ip-api.com/json"}`))
	requests = append(requests, r)

	body := `[{"correlation_id": "a6c87844-3c0a-4854-980b-cf03414f8bac2","original_url": "http://i3nwibrmf.biz/jdr41md0/xh6ni3v3qii1"}]`
	r, _ = http.NewRequest(http.MethodPost, url+"/api/shorten/batch", strings.NewReader(body))
	requests = append(requests, r)

	r, _ = http.NewRequest(http.MethodGet, url+"/api/user/urls", nil)
	requests = append(requests, r)

	r, _ = http.NewRequest(http.MethodDelete, url+"/api/user/urls", strings.NewReader(`["token_to_delete"]`))
	requests = append(requests, r)

	return requests
}

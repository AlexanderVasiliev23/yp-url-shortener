package load

import (
	"net/http"
	"strings"
	"testing"
)

const url = "http://localhost:8080"

func BenchmarkLoad(b *testing.B) {
	for range b.N {
		if err := generateLoad(); err != nil {
			b.Errorf("generateLoad failed: %v", err)
		}
	}
}

func generateLoad() error {
	client := http.Client{}
	for _, req := range buildRequests() {
		if _, err := client.Do(req); err != nil {
			return err
		}
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

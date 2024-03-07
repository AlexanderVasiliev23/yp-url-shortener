package batch

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/google"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Example демонстрация работы обработчика сохранения ссылки и получения токена для нее
func Example() {
	tokenGenerator := func() tokenGenerator {
		tokens := []string{token1, token2}
		tokensChan := make(chan string, len(tokens))
		for _, token := range tokens {
			tokensChan <- token
		}
		return tokenGeneratorMock{tokensSeq: tokensChan}
	}()

	handler := NewShortener(
		batchSaverMock{},
		tokenGenerator,
		google.UUIDGenerator{},
		&mock.UserContextFetcherMock{},
		addr).Handle

	req := fmt.Sprintf(`
		[
			{"correlation_id": "%s", "original_url": "https://test_url1.com"},
			{"correlation_id": "%s", "original_url": "https://test_url2.com"}
		]`,
		correlationID1, correlationID2)
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(req))
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)

	_ = handler(c)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 201
	// [{"correlation_id":"correlationId1","short_url":"https://my_url_shortener/token1"},{"correlation_id":"correlationId2","short_url":"https://my_url_shortener/token2"}]
}

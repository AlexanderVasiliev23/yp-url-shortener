package shorten

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Example демонстрация работы обработчика массового сохранения ссылок и получения токенов для них
func Example() {
	handler := NewShortener(
		repositoryMock{},
		tokenGeneratorMock{token: defaultToken},
		&mock.UserContextFetcherMock{},
		addr).Handle

	r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://test.me"}`))
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)

	_ = handler(c)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 201
	// {"result":"https://my_url_shortener/test_token"}
}

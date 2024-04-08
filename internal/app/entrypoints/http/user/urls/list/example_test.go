package list

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

// Example демонстрация работы обработчика получения токенов и соответствущих ссылок, принадлежащих запрашивающему пользователю
func Example() {
	handler := NewHandler(
		storageMock{result: []*models.ShortLink{
			{
				Token:    defaultToken,
				Original: defaultOriginal,
			},
		}},
		&userContextFetcherMock{userID: defaultUserID},
		defaultAddr).List

	r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)

	_ = handler(c)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// [{"short_url":"https://my_url_shortener/test_token","original_url":"test_original"}]
}

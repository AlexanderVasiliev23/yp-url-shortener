package list

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/list"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

// Example демонстрация работы обработчика получения токенов и соответствущих ссылок, принадлежащих запрашивающему пользователю
func Example() {
	handler := NewHandler(&useCaseMock{
		outDTO: &list.OutDTO{
			Items: []list.OutDTOItem{
				{
					ShortURL:    "https://my_url_shortener/test_token",
					OriginalURL: "test_original",
				},
			},
		},
		err: nil,
	}).List

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

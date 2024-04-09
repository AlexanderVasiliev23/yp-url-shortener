package get

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

// Example демонстрация работы обработчика, который на вход принимает токен,
// а на выходе редиректит на соответствующий url, если для данного токена такой есть
func Example() {
	handler := NewHandler(&useCaseMock{
		originalURL: "https://github.com",
		err:         nil,
	}).Get

	r := httptest.NewRequest(http.MethodGet, "/"+defaultToken, nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)
	c.SetPath("/:token")
	c.SetParamNames("token")
	c.SetParamValues(defaultToken + "a")

	_ = handler(c)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Location"))

	// Output:
	// 307
	// https://github.com
}

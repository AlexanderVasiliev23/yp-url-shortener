package deleteurl

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Example демонстрация работы обработчика асинхронного удаления токенов
func Example() {
	//ch := make(chan deleter.DeleteTask, 1)
	handler := NewHandler(&useCaseMock{}).Delete

	r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["token1", "token2"]`))
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)

	_ = handler(c)

	fmt.Println(w.Code)

	// Output:
	// 202
}

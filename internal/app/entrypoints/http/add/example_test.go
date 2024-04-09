package add

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

// Example демонстрация работы обработчика добавления урла
func Example() {
	handler := NewHandler(&mockUseCase{}).Add

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("test_url"))
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)

	_ = handler(c)

	fmt.Println(w.Code)

	// Output:
	// 201
}

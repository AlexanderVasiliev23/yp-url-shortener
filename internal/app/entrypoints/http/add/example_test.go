package add

import (
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
)

type mockTokenGenerator struct {
	err   error
	token string
}

func (m mockTokenGenerator) Generate() (string, error) {
	return m.token, m.err
}

type mockRepo struct {
	addingErr   error
	getTokenErr error
}

func (m mockRepo) Add(ctx context.Context, shortLink *models.ShortLink) error {
	return m.addingErr
}

func (m mockRepo) GetTokenByURL(ctx context.Context, url string) (string, error) {
	return defaultToken, m.getTokenErr
}

// Example демонстрация работы обработчика добавления урла
func Example() {
	_useCase := add.NewUseCase(
		mockRepo{},
		mockTokenGenerator{token: defaultToken},
		&mock.UserContextFetcherMock{},
		addr,
	)

	handler := NewHandler(_useCase).Add

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("test_url"))
	w := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(r, w)

	_ = handler(c)

	fmt.Println(w.Code)

	// Output:
	// 201
}

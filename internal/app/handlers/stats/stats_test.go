package stats

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	defaultUrlsCount  = 45
	defaultUsersCount = 23
)

func TestHandle(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		name string
		want want
	}{
		{
			name: "success",
			want: want{
				err:  nil,
				body: fmt.Sprintf(`{"urls":%d,"users":%d}`, defaultUrlsCount, defaultUsersCount) + "\n",
				code: http.StatusOK,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//recorder := httptest.NewRecorder()
			//request := httptest.NewRequest(http.MethodGet, "/", nil)

			// todo configure it
			//h := NewHandler().Handle

			//e := echo.New()
			//c := e.NewContext(request, recorder)
			//
			//err := h(c)
			//
			//if tc.want.err == nil {
			//	assert.NoError(t, err)
			//} else {
			//	assert.ErrorIs(t, tc.want.err, err)
			//}
			//
			//assert.Equal(t, tc.want.code, recorder.Code)
			//assert.Equal(t, tc.want.body, recorder.Body.String())
		})
	}
}

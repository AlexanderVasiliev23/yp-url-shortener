package gzip

import (
	"bytes"
	"compress/gzip"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	type req struct {
		header http.Header
		body   string
	}

	type resp struct {
		header http.Header
		body   string
	}

	testCases := []struct {
		name string
		req  req
		resp resp
	}{
		{
			name: "uncompressed request body",
			req: req{
				body: "not compressed content",
			},
			resp: resp{
				body:   "not compressed content",
				header: make(http.Header),
			},
		},
		{
			name: "compressed json in request body",
			req: req{
				body: compressString(`{"hi":"there"}`),
				header: http.Header{
					"Content-Encoding": []string{"gzip"},
				},
			},
			resp: resp{
				body:   `{"hi":"there"}`,
				header: make(http.Header),
			},
		},
		{
			name: "compressed html in request body",
			req: req{
				body: compressString(`<h1>Test header</h1>`),
				header: http.Header{
					"Content-Encoding": []string{"gzip"},
				},
			},
			resp: resp{
				body:   `<h1>Test header</h1>`,
				header: make(http.Header),
			},
		},
		{
			name: "request accepts gzip, json",
			req: req{
				body: `{"hi":"there"}`,
				header: http.Header{
					"Accept-Encoding": []string{"gzip"},
				},
			},
			resp: resp{
				body: compressString(`{"hi":"there"}`),
				header: http.Header{
					"Content-Encoding": []string{"gzip"},
				},
			},
		},
		{
			name: "request accepts gzip, html",
			req: req{
				body: `<h1>Test header</h1>`,
				header: http.Header{
					"Accept-Encoding": []string{"gzip"},
				},
			},
			resp: resp{
				body: compressString(`<h1>Test header</h1>`),
				header: http.Header{
					"Content-Encoding": []string{"gzip"},
				},
			},
		},
		{
			name: "request: doesn't accept gzip",
			req: req{
				body: `<h1>Test header</h1>`,
			},
			resp: resp{
				body:   `<h1>Test header</h1>`,
				header: make(http.Header),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			e := echo.New()
			e.Use(Middleware())
			e.POST("/", handler)

			req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.req.body))
			req.Header = tc.req.header
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tc.resp.body, rec.Body.String())
			assert.Equal(t, tc.resp.header, rec.Header())
		})
	}
}

func handler(c echo.Context) error {
	content, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	c.Response().Write(content)

	return nil
}

func compressString(s string) string {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	gzWriter.Write([]byte(s))
	gzWriter.Close()

	return buf.String()
}

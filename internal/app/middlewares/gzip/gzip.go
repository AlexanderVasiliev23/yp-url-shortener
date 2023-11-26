package gzip

import (
	"compress/gzip"
	"github.com/labstack/echo/v4"
	"net/http"
)

type gzipWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (gw *gzipWriter) Write(data []byte) (int, error) {
	return gw.writer.Write(data)
}

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if requestBodyIsCompressed(c.Request().Header) {
				gzReader, err := gzip.NewReader(c.Request().Body)
				if err != nil {
					return err
				}

				c.Request().Body = gzReader
			}

			if clientAcceptsGzip(c.Request().Header) {
				c.Response().Header().Set("Content-Encoding", "gzip")

				writer := gzip.NewWriter(c.Response().Writer)
				defer writer.Close()
				gzipWriter := &gzipWriter{
					ResponseWriter: c.Response().Writer,
					writer:         writer,
				}

				c.Response().Writer = gzipWriter
			}

			return next(c)
		}
	}
}

func clientAcceptsGzip(header http.Header) bool {
	for _, val := range header.Values("Accept-Encoding") {
		if val == "gzip" {
			return true
		}
	}

	return false
}

func requestBodyIsCompressed(header http.Header) bool {
	for _, val := range header.Values("Content-Encoding") {
		if val == "gzip" {
			return true
		}
	}

	return false
}

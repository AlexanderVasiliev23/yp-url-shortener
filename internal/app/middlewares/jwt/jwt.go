package jwt

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	jwtTokenCookieName = "jwt_token"
	wrongUserId        = -1
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func Middleware(JWTSecretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId := getUserIdFromCookie(c, JWTSecretKey)
			if userId == wrongUserId {
				userId = generateUserId()
				if err := setCookie(c, userId, JWTSecretKey); err != nil {
					c.Response().WriteHeader(http.StatusInternalServerError)
					return err
				}
			}

			c.SetRequest(
				c.Request().WithContext(auth.WithUserId(c.Request().Context(), userId)),
			)

			err := next(c)

			return err
		}
	}
}

func getUserIdFromCookie(c echo.Context, JWTSecretKey string) int {
	jwtCookie, err := c.Cookie(jwtTokenCookieName)

	if err != nil {
		return wrongUserId
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(jwtCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return wrongUserId
	}

	if !token.Valid {
		return wrongUserId
	}

	return claims.UserID
}

func generateUserId() int {
	return int(time.Now().UnixNano())
}

func setCookie(c echo.Context, userId int, JWTSecretKey string) error {
	token, err := buildJWTString(userId, JWTSecretKey)
	if err != nil {
		return fmt.Errorf("build JWT string: %w", err)
	}
	c.SetCookie(&http.Cookie{
		Name:  jwtTokenCookieName,
		Value: token,
	})

	return nil
}

func buildJWTString(userId int, JWTSecretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userId,
	})

	tokenString, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("sign string using secret key: %w", err)
	}

	return tokenString, err
}

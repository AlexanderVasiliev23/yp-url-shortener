package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth"
)

const (
	jwtTokenCookieName = "jwt_token"
)

var (
	errCookieNotFound = errors.New("cookie not found")
	errTokenParsing   = errors.New("token parsing")
	errInvalidJWT     = errors.New("invalid jwt")
	errUserIDNotSet   = errors.New("user id is not set")
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func Auth(JWTSecretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, err := getUserIDFromCookie(c, JWTSecretKey); err != nil {
				c.Response().WriteHeader(http.StatusUnauthorized)
				return err
			}

			return next(c)
		}
	}
}

func Middleware(JWTSecretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := getUserIDFromCookie(c, JWTSecretKey)
			if err != nil {
				if errors.Is(err, errUserIDNotSet) {
					c.Response().WriteHeader(http.StatusUnauthorized)
					return err
				}

				if errors.Is(err, errCookieNotFound) || errors.Is(err, errTokenParsing) || errors.Is(err, errInvalidJWT) {
					userID = generateUserID()
					if err := setCookie(c, userID, JWTSecretKey); err != nil {
						c.Response().WriteHeader(http.StatusInternalServerError)
						return err
					}
				} else {
					c.Response().WriteHeader(http.StatusInternalServerError)
					return err
				}
			}

			c.SetRequest(
				c.Request().WithContext(auth.WithUserID(c.Request().Context(), userID)),
			)

			err = next(c)

			return err
		}
	}
}

func getUserIDFromCookie(c echo.Context, JWTSecretKey string) (int, error) {
	jwtCookie, err := c.Cookie(jwtTokenCookieName)

	if err != nil {
		return 0, errCookieNotFound
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(jwtCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return 0, errTokenParsing
	}

	if !token.Valid {
		return 0, errInvalidJWT
	}

	if claims.UserID == 0 {
		return 0, errUserIDNotSet
	}

	return claims.UserID, nil
}

func generateUserID() int {
	return int(time.Now().UnixNano())
}

func setCookie(c echo.Context, userID int, JWTSecretKey string) error {
	token, err := buildJWTString(userID, JWTSecretKey)
	if err != nil {
		return fmt.Errorf("build JWT string: %w", err)
	}
	c.SetCookie(&http.Cookie{
		Name:  jwtTokenCookieName,
		Value: token,
	})

	return nil
}

func buildJWTString(userID int, JWTSecretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("sign string using secret key: %w", err)
	}

	return tokenString, err
}

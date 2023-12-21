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
	wrongUserID        = -1
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func Middleware(JWTSecretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := getUserIDFromCookie(c, JWTSecretKey)
			if userID == wrongUserID {
				userID = generateUserID()
				if err := setCookie(c, userID, JWTSecretKey); err != nil {
					c.Response().WriteHeader(http.StatusInternalServerError)
					return err
				}
			}

			c.SetRequest(
				c.Request().WithContext(auth.WithUserID(c.Request().Context(), userID)),
			)

			err := next(c)

			return err
		}
	}
}

func getUserIDFromCookie(c echo.Context, JWTSecretKey string) int {
	jwtCookie, err := c.Cookie(jwtTokenCookieName)

	if err != nil {
		return wrongUserID
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(jwtCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return wrongUserID
	}

	if !token.Valid {
		return wrongUserID
	}

	return claims.UserID
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

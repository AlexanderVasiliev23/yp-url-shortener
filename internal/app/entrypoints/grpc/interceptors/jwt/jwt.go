package jwt

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth"
	"github.com/golang-jwt/jwt/v4"
)

const (
	jwtTokenFieldName = "jwt_token"
)

var (
	errTokenParsing = errors.New("token parsing")
	errInvalidJWT   = errors.New("invalid jwt")
	errUserIDNotSet = errors.New("user id is not set")
)

// Claims missing godoc.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func UnaryInterceptor(jwtSecretKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var token string

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get(jwtTokenFieldName)
			if len(values) > 0 {
				token = values[0]
			}
		}

		userID, err := userIdFromToken(token, jwtSecretKey)
		if err != nil {
			if errors.Is(err, errUserIDNotSet) {
				return nil, status.Errorf(codes.Unauthenticated, errUserIDNotSet.Error())
			}

			if errors.Is(err, errTokenParsing) || errors.Is(err, errInvalidJWT) {
				userID = generateUserID()
			} else {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}

		ctx = auth.WithUserID(ctx, userID)

		return handler(ctx, req)
	}
}

func userIdFromToken(rawToken string, jwtSecretKey string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
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

package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/pkg/jwthelper"
)

type Authenticator struct {
	JWTSigningKey string
}

func NewAuthenticator(jwtSigningKey string) *Authenticator {
	return &Authenticator{
		JWTSigningKey: jwtSigningKey,
	}
}

func (a *Authenticator) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := a.extractClaims(r)
		if err != nil {
			_ = render.Render(w, r, response.ErrJWTUnverified(err))

			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) extractClaims(r *http.Request) (*jwthelper.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("header Authorization not found")
	}

	// We expect the header like "Authorization: BEARER T",
	// only the portion after BEARER is needed.
	if len(authHeader) <= 7 || strings.ToUpper(authHeader[0:6]) != "BEARER" {
		return nil, errors.New("wrong JWT format")
	}

	bearer := authHeader[7:]
	claims := &jwthelper.Claims{} // must pass pointer into ParseWithClaims.
	token, err := jwthelper.ParseWithClaims(a.JWTSigningKey, bearer, claims)
	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, errors.New("token is invalid or already expired")
	}

	if claims.UserID == 0 {
		return nil, errors.New("userID not found in the claims")
	}

	if claims.UserAgent != r.Header.Get("User-Agent") {
		return nil, errors.New("user agent do not match")
	}

	return claims, nil
}

package authmiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKeyClaims struct{}

type RequireRoleMiddleware struct {
	Roles []string
	Next  http.Handler
}

type CustomClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get token string
			tokenStr := r.Header.Get("Authorization")

			// validate valid tokenStr
			if tokenStr == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			// parse token with claim
			token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected sighning method")
				}
				return secret, nil
			})

			// check validity
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
			}

			// access claims
			claims := token.Claims.(*CustomClaims)

			// run next function
			ctx := context.WithValue(r.Context(), ctxKeyClaims{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		mw := RequireRoleMiddleware{
			Roles: requiredRoles,
			Next:  next,
		}
		return http.HandlerFunc(mw.ServeHTTP)
	}

}

func (h *RequireRoleMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ctxKeyClaims{}).(*CustomClaims)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	for _, role := range h.Roles {
		for _, userRole := range claims.Roles {
			if userRole == role {
				h.Next.ServeHTTP(w, r)
				return
			}
		}
	}
	http.Error(w, "forbidden", http.StatusForbidden)
}

func GetClaims(r *http.Request) (*CustomClaims, bool) {
	claims, ok := r.Context().Value(ctxKeyClaims{}).(*CustomClaims)
	return claims, ok
}

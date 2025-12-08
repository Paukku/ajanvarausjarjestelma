package server

import (
	"net/http"
	"strings"

	role "github.com/Paukku/ajanvarausjarjestelma/backend/internal/auth"
	pbcommon "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	"github.com/golang-jwt/jwt/v5"
)

func RoleMiddleware(required pbcommon.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if required == pbcommon.UserRole_UNAUTHORIZED {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			rawRole, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "role must be a string", http.StatusForbidden)
				return
			}

			userRole, exists := role.RoleStringToEnum[role.RoleString(rawRole)]
			if !exists {
				http.Error(w, "unknown role", http.StatusForbidden)
				return
			}

			if userRole < required {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

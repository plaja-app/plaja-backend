package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plaja-app/back-end/models"
	"net/http"
	"time"
)

// RequireAuth is a middleware that checks for the presence and validity of a JWT in the request cookie.
func (m *BaseMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the cookie of request
		tokenCookie, err := r.Cookie("pja_user_jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// decode/validate it
		token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(m.App.Env.JWTSecret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// find the user with token sub
			var user models.User

			m.App.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

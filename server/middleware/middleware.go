package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func Authorization() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt_token")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Println(http.StatusUnauthorized, "Error retrieving cookie: ")
					http.Error(w, "no cookie found", http.StatusUnauthorized)
					return
				}
				log.Println("Error retrieving cookie: ", err.Error())
				http.Error(w, "Error retrieving cookie", http.StatusUnauthorized)
				return
			}
			jwt := cookie.Value
			claims, err := DecodeJWT(jwt)
			if err != nil {
				log.Println(http.StatusBadRequest, "invalid token")
				http.Error(w, "invalid token", http.StatusBadRequest)
				return
			}
			userID, ok := claims["user_id"].(float64)
			if !ok {
				log.Println(http.StatusUnauthorized, "invalid claims")
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
			}
			ctx := context.WithValue(r.Context(), "userID", int(userID))
			log.Println(http.StatusOK, "Authorization Successful")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

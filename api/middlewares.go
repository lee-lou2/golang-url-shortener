package api

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"log"
	"net/http"
)

// RequestIdMiddleware Request ID 미들웨어
func RequestIdMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 존재 여부 확인
		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			// UUID 생성
			requestId = uuid.New().String()
			r.Header.Set("X-Request-ID", requestId)
		}
		w.Header().Set("X-Request-ID", requestId)
		next.ServeHTTP(w, r)
	}
}

// RecoverMiddleware recover 미들웨어
func RecoverMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				sentry.CaptureException(errors.New(r.(string)))
				log.Println("panic:", r)
				http.Error(w, r.(string), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// LoggerMiddleware 로거 미들웨어
func LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

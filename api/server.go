package api

import (
	"net/http"
)

func Server() http.HandlerFunc {
	mux := http.NewServeMux()

	// 라우터
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/robots.txt")
	})
	mux.HandleFunc("POST /v1/urls", createShortUrlHandler)
	mux.HandleFunc("GET /{short_key}", redirectShortUrlHandler)

	// 미들웨어
	chain := http.HandlerFunc(mux.ServeHTTP)
	chain = RequestIdMiddleware(chain)
	chain = RecoverMiddleware(chain)
	chain = LoggerMiddleware(chain)
	return chain
}

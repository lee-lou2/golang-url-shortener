package main

import (
	"context"
	"errors"
	"fit/api"
	"fit/cmd"
	"fit/config"
	"fit/pkg/loggers"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 로그 설정
	writer := &loggers.LokiLoggerWriter{
		App: "fit",
		Env: "prod",
	}
	log.SetOutput(writer)
	log.SetFlags(0)

	// 센트리 설정
	sentryDsn := config.GetEnv("SENTRY_DSN")
	_ = sentry.Init(sentry.ClientOptions{Dsn: sentryDsn})

	// 서버 실행
	server := api.Server()
	serverPort := config.GetEnv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3000"
	}
	srv := &http.Server{Addr: ":" + serverPort, Handler: server}

	// OS 신호 채널 설정
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 서버 실행
	go func() {
		log.Println("Run server on port " + serverPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 종료 신호가 있을때까지 대기
	<-quit
	log.Println("Server is shutting down")

	// 종료 전에 실행할 함수 호출
	_ = cmd.Teardown()
	_ = sentry.Flush(2 * time.Second)

	// 서버 종료 절차
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server is gracefully stopped")
}

package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"backstock/internal/migrate"
	"backstock/internal/server"
	"backstock/internal/store"
)

func main() {
	dbPath := os.Getenv("BACKSTOCK_DB_PATH")
	if dbPath == "" {
		dbPath = "backstock.db"
	}

	addr := os.Getenv("BACKSTOCK_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	s, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer s.Close()

	if err := migrate.Run(s.DB()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	frontend, err := fs.Sub(frontendDist, "frontend/dist")
	if err != nil {
		log.Fatalf("failed to get frontend fs: %v", err)
	}

	handler := server.New(s, frontend)

	httpSrv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpSrv.Shutdown(shutdownCtx)
	}()

	log.Printf("starting server on %s", addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

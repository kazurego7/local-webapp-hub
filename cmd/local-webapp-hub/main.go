package main

import (
    "context"
    "errors"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/kazurego7/local-webapp-hub/internal/scan"
    "github.com/kazurego7/local-webapp-hub/internal/server"
)

var (
    listenAddr = flag.String("addr", ":8787", "listen address, e.g. :8787")
)

func main() {
    flag.Parse()

    s := scan.NewDefault()
    e := server.New(s, *listenAddr)

    log.Printf("local-webapp-hub listening on %s", *listenAddr)

    // Graceful shutdown
    go func() {
        if err := e.Start(*listenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("server error: %v", err)
        }
    }()

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    <-ctx.Done()
    log.Printf("shutting down...")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = e.Shutdown(shutdownCtx)
}


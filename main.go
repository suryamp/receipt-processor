package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "github.com/gorilla/mux"
    "github.com/suryamp/receipt-processor/handlers"
    "github.com/suryamp/receipt-processor/processor"
)

func main() {
    // Initialize processor and handler
    receiptProcessor := &processor.InMemoryProcessor{}
    handler := handlers.NewHandler(receiptProcessor)

    // Set up router
    r := mux.NewRouter()
    r.HandleFunc("/receipts/process", handler.ProcessReceiptHandler).Methods("POST")
    r.HandleFunc("/receipts/{id}/points", handler.GetPointsHandler).Methods("GET")

    // Configure server
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Start server
    go func() {
        log.Printf("Server starting on port 8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited gracefully")
}

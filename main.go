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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/suryamp/receipt-processor/handlers"
	"github.com/suryamp/receipt-processor/logger"
	"github.com/suryamp/receipt-processor/middleware"
	"github.com/suryamp/receipt-processor/processor"
)

var handler *handlers.Handler
var receiptProcessor processor.ReceiptProcessor

func init() {
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}

func main() {
	receiptProcessor = processor.NewInMemoryProcessor()
	handler = handlers.NewHandler(receiptProcessor)

	// Set up router
	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	// Add metrics middleware
	r.Use(middleware.MetricsMiddleware)

	var okResponse = []byte("OK")
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write(okResponse)
	})

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
		logger.InfoLogger.Printf("Server starting on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorLogger.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.InfoLogger.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.ErrorLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.InfoLogger.Printf("Server exited gracefully")
}

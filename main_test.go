package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/suryamp/receipt-processor/handlers"
	"github.com/suryamp/receipt-processor/models"
	"github.com/suryamp/receipt-processor/processor"
)

func TestIntegration(t *testing.T) {
	// Setup
	receiptProcessor := &processor.InMemoryProcessor{}
	handler := handlers.NewHandler(receiptProcessor)

	// Test case: Process receipt and get points
	t.Run("process receipt and get points", func(t *testing.T) {
		// Create test receipt
		receipt := models.Receipt{
			Retailer:     "Target",
			PurchaseDate: "2024-01-01",
			PurchaseTime: "13:01",
			Total:        "35.00", // Round dollar amount for predictable points
			Items: []models.Item{
				{ShortDescription: "Mountain Dew", Price: "1.25"},
				{ShortDescription: "Pepsi", Price: "2.00"},
			},
		}

		// Process receipt
		body, err := json.Marshal(receipt)
		if err != nil {
			t.Fatalf("Failed to marshal receipt: %v", err)
		}

		req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.ProcessReceiptHandler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("ProcessReceiptHandler() status = %v, want %v", w.Code, http.StatusOK)
		}

		var processResp models.ProcessResponse
		if err := json.NewDecoder(w.Body).Decode(&processResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Get points using the ID from previous response
		req = httptest.NewRequest("GET", "/receipts/"+processResp.ID+"/points", nil)
		w = httptest.NewRecorder()

		// Need to set up router to get URL parameters
		router := mux.NewRouter()
		router.HandleFunc("/receipts/{id}/points", handler.GetPointsHandler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GetPointsHandler() status = %v, want %v", w.Code, http.StatusOK)
		}

		var pointsResp models.PointsResponse
		if err := json.NewDecoder(w.Body).Decode(&pointsResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify points calculation (based on the test receipt)
		expectedPoints := int64(93) // 6 (retailer) + 50 (round dollar) 25 (multiple of 0.25) + 5 (2 items) + 1 (item desciption length) + 6 (odd day)
		if pointsResp.Points != expectedPoints {
			t.Errorf("GetPointsHandler() points = %v, want %v", pointsResp.Points, expectedPoints)
		}
	})

	// Test case: Get points for non-existent receipt
	t.Run("get points for non-existent receipt", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/receipts/non-existent-id/points", nil)
		w := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/receipts/{id}/points", handler.GetPointsHandler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GetPointsHandler() status = %v, want %v", w.Code, http.StatusNotFound)
		}

		expected := "No receipt found for that ID.\n"
		if got := w.Body.String(); strings.Compare(got, expected) != 0 {
			t.Errorf("GetPointsHandler() response = %v, want %v", got, expected)
		}
	})
}

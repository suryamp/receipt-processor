package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/suryamp/receipt-processor/logger"
	"github.com/suryamp/receipt-processor/models"
)

func init() {
	if err := logger.Init(); err != nil {
		panic(err)
	}
}

// MockProcessor implements processor.ReceiptProcessor for testing
type MockProcessor struct {
	shouldError bool
	points      int64
}

func (m *MockProcessor) ProcessReceipt(receipt models.Receipt) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("mock error")
	}
	return "test-id", nil
}

func (m *MockProcessor) GetPoints(id string) (int64, error) {
	if m.shouldError {
		return 0, fmt.Errorf("mock error")
	}
	return m.points, nil
}

func TestProcessReceiptHandler(t *testing.T) {
	tests := []struct {
		name         string
		receipt      models.Receipt
		shouldError  bool
		wantStatus   int
		wantResponse string
	}{
		{
			name: "valid receipt",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.25"},
				},
			},
			shouldError:  false,
			wantStatus:   http.StatusOK,
			wantResponse: `{"id":"test-id"}`,
		},
		{
			name:         "invalid json",
			receipt:      models.Receipt{}, // Will send invalid JSON in test
			shouldError:  false,
			wantStatus:   http.StatusBadRequest,
			wantResponse: "The receipt is invalid.\n",
		},
		{
			name: "processor error",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.25"},
				},
			},
			shouldError:  true,
			wantStatus:   http.StatusBadRequest,
			wantResponse: "The receipt is invalid.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock processor
			mockProc := &MockProcessor{shouldError: tt.shouldError}
			handler := NewHandler(mockProc)

			// Create request
			var body []byte
			var err error
			if tt.name == "invalid json" {
				body = []byte(`{invalid json}`)
			} else {
				body, err = json.Marshal(tt.receipt)
				if err != nil {
					t.Fatalf("Failed to marshal receipt: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Call handler
			handler.ProcessReceiptHandler(w, req)

			// Check status code
			if got := w.Code; got != tt.wantStatus {
				t.Errorf("ProcessReceiptHandler() status = %v, want %v", got, tt.wantStatus)
			}

			// Check response
			if tt.wantStatus == http.StatusOK {
				var got models.ProcessResponse
				if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if got.ID != "test-id" {
					t.Errorf("ProcessReceiptHandler() response = %v, want %v", got.ID, "test-id")
				}
			} else {
				if got := w.Body.String(); got != tt.wantResponse {
					t.Errorf("ProcessReceiptHandler() response = %v, want %v", got, tt.wantResponse)
				}
			}
		})
	}
}

func TestGetPointsHandler(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		points       int64
		shouldError  bool
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "valid id",
			id:           "test-id",
			points:       100,
			shouldError:  false,
			wantStatus:   http.StatusOK,
			wantResponse: `{"points":100}`,
		},
		{
			name:         "not found",
			id:           "invalid-id",
			shouldError:  true,
			wantStatus:   http.StatusNotFound,
			wantResponse: "No receipt found for that ID.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock processor
			mockProc := &MockProcessor{
				shouldError: tt.shouldError,
				points:      tt.points,
			}
			handler := NewHandler(mockProc)

			// Create request with mux vars
			req := httptest.NewRequest("GET", "/receipts/{id}/points", nil)
			w := httptest.NewRecorder()

			// Add URL parameters to request
			vars := map[string]string{
				"id": tt.id,
			}
			req = mux.SetURLVars(req, vars)

			// Call handler
			handler.GetPointsHandler(w, req)

			// Check status code
			if got := w.Code; got != tt.wantStatus {
				t.Errorf("GetPointsHandler() status = %v, want %v", got, tt.wantStatus)
			}

			// Check response
			if tt.wantStatus == http.StatusOK {
				var got models.PointsResponse
				if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if got.Points != tt.points {
					t.Errorf("GetPointsHandler() response = %v, want %v", got.Points, tt.points)
				}
			} else {
				if got := w.Body.String(); got != tt.wantResponse {
					t.Errorf("GetPointsHandler() response = %v, want %v", got, tt.wantResponse)
				}
			}
		})
	}
}

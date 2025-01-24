package validator

import (
	"testing"

	"github.com/suryamp/receipt-processor/models"
)

func TestValidateReceipt(t *testing.T) {
	tests := []struct {
		name    string
		receipt models.Receipt
		wantErr bool
		errMsg  string
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
			wantErr: false,
		},
		{
			name: "invalid retailer with special chars",
			receipt: models.Receipt{
				Retailer:     "Target@#$",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.25"},
				},
			},
			wantErr: true,
			errMsg:  "invalid retailer format",
		},
		{
			name: "invalid date format",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "01-01-2024", // wrong format
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.25"},
				},
			},
			wantErr: true,
			errMsg:  "invalid date format",
		},
		{
			name: "invalid time format",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "1:01 PM", // wrong format
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.25"},
				},
			},
			wantErr: true,
			errMsg:  "invalid time format",
		},
		{
			name: "invalid total format",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.5", // missing second decimal
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.25"},
				},
			},
			wantErr: true,
			errMsg:  "invalid total format",
		},
		{
			name: "no items",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items:        []models.Item{},
			},
			wantErr: true,
			errMsg:  "at least one item required",
		},
		{
			name: "invalid item description",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain @Dew", Price: "1.25"},
				},
			},
			wantErr: true,
			errMsg:  "invalid item description format",
		},
		{
			name: "invalid item price",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew", Price: "1.5"},
				},
			},
			wantErr: true,
			errMsg:  "invalid item price format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReceipt(tt.receipt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReceipt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ValidateReceipt() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

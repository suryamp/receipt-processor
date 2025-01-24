package processor

import (
	"testing"

	"github.com/suryamp/receipt-processor/logger"
	"github.com/suryamp/receipt-processor/models"
)

func init() {
	if err := logger.Init(); err != nil {
		panic(err)
	}
}

func TestCalculateRetailerNamePoints(t *testing.T) {
	tests := []struct {
		name     string
		retailer string
		want     int64
	}{
		{
			name:     "simple retailer name",
			retailer: "Target",
			want:     6,
		},
		{
			name:     "retailer name with spaces",
			retailer: "M&M Corner Market",
			want:     14,
		},
		{
			name:     "retailer with special chars",
			retailer: "Target!!!",
			want:     6,
		},
		{
			name:     "empty retailer",
			retailer: "",
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRetailerNamePoints(tt.retailer); got != tt.want {
				t.Errorf("calculateRetailerNamePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateRoundDollarPoints(t *testing.T) {
	tests := []struct {
		name  string
		total string
		want  int64
	}{
		{
			name:  "round dollar amount",
			total: "35.00",
			want:  50,
		},
		{
			name:  "not round dollar amount",
			total: "35.99",
			want:  0,
		},
		{
			name:  "zero amount",
			total: "0.00",
			want:  50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRoundDollarPoints(tt.total); got != tt.want {
				t.Errorf("calculateRoundDollarPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateQuarterPoints(t *testing.T) {
	tests := []struct {
		name  string
		total string
		want  int64
	}{
		{
			name:  "quarter dollar amount",
			total: "10.25",
			want:  25,
		},
		{
			name:  "multiple quarter dollar amount",
			total: "10.75",
			want:  25,
		},
		{
			name:  "not quarter dollar amount",
			total: "10.20",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateQuarterPoints(tt.total); got != tt.want {
				t.Errorf("calculateQuarterPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateItemCountPoints(t *testing.T) {
	tests := []struct {
		name  string
		items []models.Item
		want  int64
	}{
		{
			name: "two items",
			items: []models.Item{
				{ShortDescription: "Item 1", Price: "10.00"},
				{ShortDescription: "Item 2", Price: "20.00"},
			},
			want: 5,
		},
		{
			name: "three items",
			items: []models.Item{
				{ShortDescription: "Item 1", Price: "10.00"},
				{ShortDescription: "Item 2", Price: "20.00"},
				{ShortDescription: "Item 3", Price: "30.00"},
			},
			want: 5,
		},
		{
			name:  "no items",
			items: []models.Item{},
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateItemCountPoints(tt.items); got != tt.want {
				t.Errorf("calculateItemCountPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateItemDescriptionPoints(t *testing.T) {
	tests := []struct {
		name  string
		items []models.Item
		want  int64
	}{
		{
			name: "description length divisible by 3",
			items: []models.Item{
				{ShortDescription: "abc", Price: "10.00"}, // length 3
			},
			want: 2, // ceil(10.00 * 0.2)
		},
		{
			name: "description length not divisible by 3",
			items: []models.Item{
				{ShortDescription: "abcd", Price: "10.00"}, // length 4
			},
			want: 0,
		},
		{
			name: "multiple items with mixed lengths",
			items: []models.Item{
				{ShortDescription: "abc", Price: "10.00"},    // length 3
				{ShortDescription: "abcd", Price: "20.00"},   // length 4
				{ShortDescription: "abcdef", Price: "30.00"}, // length 6
			},
			want: 8, // ceil(10.00 * 0.2) + ceil(30.00 * 0.2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateItemDescriptionPoints(tt.items); got != tt.want {
				t.Errorf("calculateItemDescriptionPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateOddDayPoints(t *testing.T) {
	tests := []struct {
		name         string
		purchaseDate string
		want         int64
	}{
		{
			name:         "odd day",
			purchaseDate: "2024-01-01",
			want:         6,
		},
		{
			name:         "even day",
			purchaseDate: "2024-01-02",
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateOddDayPoints(tt.purchaseDate); got != tt.want {
				t.Errorf("calculateOddDayPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateHappyHourPoints(t *testing.T) {
	tests := []struct {
		name         string
		purchaseTime string
		want         int64
	}{
		{
			name:         "during happy hour",
			purchaseTime: "14:30",
			want:         10,
		},
		{
			name:         "before happy hour",
			purchaseTime: "13:59",
			want:         0,
		},
		{
			name:         "after happy hour",
			purchaseTime: "16:01",
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateHappyHourPoints(tt.purchaseTime); got != tt.want {
				t.Errorf("calculateHappyHourPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name    string
		receipt models.Receipt
		want    int64
	}{
		{
			name: "complete receipt example",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2024-01-01",
				PurchaseTime: "14:30",
				Items: []models.Item{
					{ShortDescription: "abc", Price: "10.00"},
					{ShortDescription: "def", Price: "20.00"},
				},
				Total: "30.00",
			},
			want: 108, // 6 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 5 (2 items) + 2 (desc points) + 0 (desc points) + 0 (quarter points) + 10 (happy hour)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePoints(tt.receipt); got != tt.want {
				t.Errorf("calculatePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

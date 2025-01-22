package validator

import (
    "fmt"
    "regexp"
    "time"
    "github.com/yourusername/receipt-processor/models"
)

// Modified regex from api.yml
var (
    retailerPattern = regexp.MustCompile(`^[\w\s\-&]+$`)
    pricePattern   = regexp.MustCompile(`^\d+\.\d{2}$`)
    descPattern    = regexp.MustCompile(`^[\w\s\-]+$`)
)

func ValidateReceipt(r models.Receipt) error {
    if !retailerPattern.MatchString(r.Retailer) {
        return fmt.Errorf("invalid retailer format")
    }
	
	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
        return fmt.Errorf("invalid date format")
    }

    if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
        return fmt.Errorf("invalid time format")
    }
    
    if !pricePattern.MatchString(r.Total) {
        return fmt.Errorf("invalid total format")
    }

    if len(r.Items) < 1 {
        return fmt.Errorf("at least one item required")
    }

    for _, item := range r.Items {
        if !descPattern.MatchString(item.ShortDescription) {
            return fmt.Errorf("invalid item description format")
        }
        if !pricePattern.MatchString(item.Price) {
            return fmt.Errorf("invalid item price format")
        }
    }

    return nil
}

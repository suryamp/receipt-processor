package processor

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suryamp/receipt-processor/logger"
	"github.com/suryamp/receipt-processor/models"
)

const (
	RetailerNameMultiplier          = 1
	RoundDollarPoints               = 50
	QuarterDollarPoints             = 25
	ItemDescriptionPointsModulus    = 3
	ItemDescriptionPointsMultiplier = 0.2
	OddDayPoints                    = 6
	HappyHourPoints                 = 10
	ItemPairPoints                  = 5
	HappyHourStart                  = 14
	HappyHourEnd                    = 16
)

var (
	alphanumericRegex = regexp.MustCompile(`[a-zA-Z0-9]`)
)

// Interface for business logic
type ReceiptProcessor interface {
	ProcessReceipt(receipt models.Receipt) (string, error)
	GetPoints(id string) (int64, error)
}

// InMemoryProcessor implements ReceiptProcessor with in-memory storage
type InMemoryProcessor struct {
	receipts sync.Map // thread-safe map for storing receipts
}

func NewInMemoryProcessor() ReceiptProcessor {
	logger.InfoLogger.Printf("Initializing receipt processor...")
	return &InMemoryProcessor{
		receipts: sync.Map{},
	}
}

func (p *InMemoryProcessor) ProcessReceipt(receipt models.Receipt) (string, error) {
	id := uuid.New().String()
	p.receipts.Store(id, receipt)
	logger.InfoLogger.Printf("Processed new receipt with ID: %s", id)
	return id, nil
}

func (p *InMemoryProcessor) GetPoints(id string) (int64, error) {
	value, ok := p.receipts.Load(id)
	if !ok {
		return 0, fmt.Errorf("No receipt found for that ID.")
	}

	receipt, ok := value.(models.Receipt)

	points := calculatePoints(receipt)
	return points, nil
}

// Points calculation rules are based on various aspects of the receipt
func calculatePoints(receipt models.Receipt) int64 {
	var points int64

	// Points from retailer name: points for every alphanumeric character
	points += calculateRetailerNamePoints(receipt.Retailer)

	// Points from total amount: points for round dollar amounts (no cents)
	points += calculateRoundDollarPoints(receipt.Total)

	// Points from total amount: points if total is a multiple of 0.25
	points += calculateQuarterPoints(receipt.Total)

	// Points from item count: points for every two items
	points += calculateItemCountPoints(receipt.Items)

	// Points from item descriptions: points for length of trimmed description
	points += calculateItemDescriptionPoints(receipt.Items)

	// Points from purchase date: points if the day is odd
	points += calculateOddDayPoints(receipt.PurchaseDate)

	// Points from purchase time: points if time is in the happy hour timeframe
	points += calculateHappyHourPoints(receipt.PurchaseTime)

	logger.InfoLogger.Printf("Total points calculated for receipt: %d", points)
	return points
}

// calculateRetailerNamePoints awards one point for every alphanumeric character in the retailer name.
// Example: "Target" = (6 * RetailerNameMultiplier) points, "M&M Corner Market" = (14 * RetailerNameMultiplier) points
func calculateRetailerNamePoints(retailer string) int64 {
	matches := alphanumericRegex.FindAllString(retailer, -1)
	points := int64(len(matches)) * RetailerNameMultiplier
	logger.InfoLogger.Printf("Retailer name '%s' earned %d points for %d alphanumeric characters",
		retailer, points, len(matches))
	return points
}

// calculateRoundDollarPoints awards RoundDollarPoints points if the total amount has no cents.
// Example: "35.00" = RoundDollarPoints points, "35.99" = 0 points
func calculateRoundDollarPoints(total string) int64 {
	if strings.HasSuffix(total, ".00") {
		logger.InfoLogger.Printf("Round dollar amount found: %s", total)
		return RoundDollarPoints
	}
	return 0
}

// calculateQuarterPoints awards QuarterDollarPoints points if the total is a multiple of 0.25.
// Example: "35.25" = QuarterDollarPoints points, "35.99" = 0 points
func calculateQuarterPoints(total string) int64 {
	if amount, err := strconv.ParseFloat(total, 64); err == nil {
		if math.Mod(amount*100, 25) == 0 {
			logger.InfoLogger.Printf("Quarter dollar amount found: %s", total)
			return QuarterDollarPoints
		}
	}
	return 0
}

// calculateItemCountPoints 5 points for every two items on the receipt.
// Example: 3 items = 5 points, 1 item = 0 points, 4 items = 10 points
func calculateItemCountPoints(items []models.Item) int64 {
	points := int64((len(items) / 2) * ItemPairPoints)
	logger.InfoLogger.Printf("Item count points: %d for %d items", points, len(items))
	return points
}

// calculateItemDescriptionPoints awards points based on item descriptions.
// For each item:
// 1. If the trimmed length of the item description is a multiple of ItemDescriptionPointsModulus
// 2. Multiply the price by ItemDescriptionPointsMultiplier and round up to nearest integer
// Example: if "Mountain Dew" was divible by ItemDescriptionPointsModulus and it had price "2.25" = ceil(2.25 * ItemDescriptionPointsMultiplier) points
func calculateItemDescriptionPoints(items []models.Item) int64 {
	var points int64
	for _, item := range items {

		trimLen := len(strings.TrimSpace(item.ShortDescription))

		if trimLen%ItemDescriptionPointsModulus == 0 {

			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				itemDescriptionPoints := int64(math.Ceil(price * ItemDescriptionPointsMultiplier))
				logger.InfoLogger.Printf("Item '%s' earned %d points (description length %d is divisible by %d)", item.ShortDescription, itemDescriptionPoints, trimLen, ItemDescriptionPointsModulus)
				points += itemDescriptionPoints
			}
		}
	}
	logger.InfoLogger.Printf("Total points from item descriptions: %d", points)
	return points
}

// calculateOddDayPoints awards OddDayPoints points if the day in the purchase date is odd.
// Example: 12/31/2025 = OddDayPoints points, 01/12/2024 = 0 points
func calculateOddDayPoints(purchaseDate string) int64 {
	if day, err := strconv.Atoi(purchaseDate[8:]); err == nil {
		if day%2 == 1 {
			logger.InfoLogger.Printf("Odd day points awarded for day: %d", day)
			return OddDayPoints
		}
	}
	return 0
}

// calculateHappyHourPoints awards HappyHourPoints points if time is in the happy hour timeframe
// Example: 3:33PM = 6 points, 7:45AM = 0 points (if happy hour started at 3PM and ended at 5PM)
func calculateHappyHourPoints(purchaseTimeString string) int64 {
	if purchaseTime, err := time.Parse("15:04", purchaseTimeString); err == nil {
		startTime := time.Date(0, 1, 1, HappyHourStart, 0, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, HappyHourEnd, 0, 0, 0, time.UTC)
		checkTime := time.Date(0, 1, 1, purchaseTime.Hour(), purchaseTime.Minute(), 0, 0, time.UTC)

		if checkTime.After(startTime) && checkTime.Before(endTime) {
			logger.InfoLogger.Printf("Happy hour points awarded for time: %s", purchaseTimeString)
			return HappyHourPoints
		}
	}
	return 0
}

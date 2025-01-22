package processor

import (
    "fmt"
	"log"
	"os"
    "regexp"
	"sync"
    "strings"
    "strconv"
    "math"
    "time"
    "github.com/google/uuid"
    "github.com/suryamp/receipt-processor/models"
)

const (
    RoundDollarPoints = 50
    QuarterDollarPoints = 25
    OddDayPoints = 6
    HappyHourPoints = 10
    ItemPairPoints = 5
    HappyHourStart = 14
    HappyHourEnd = 16
)

var (
    alphanumericRegex = regexp.MustCompile(`[a-zA-Z0-9]`)
    logger = log.New(os.Stdout, "[RECEIPT-PROCESSOR] ", log.LstdFlags)
)

func init() {
    logger.Println("Initializing receipt processor...")
}


// Interface for business logic
type ReceiptProcessor interface {
    ProcessReceipt(receipt models.Receipt) (string, error)
    GetPoints(id string) (int64, error)
}

// InMemoryProcessor implements ReceiptProcessor with in-memory storage
type InMemoryProcessor struct {
    receipts sync.Map // thread-safe map for storing receipts
}

func (p *InMemoryProcessor) ProcessReceipt(receipt models.Receipt) (string, error) {
    id := uuid.New().String()
    p.receipts.Store(id, receipt)
    logger.Printf("Processed new receipt with ID: %s", id)
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

func calculatePoints(receipt models.Receipt) int64 {
    var points int64

	points += calculateRetailerNamePoints(receipt.Retailer)
    points += calculateRoundDollarPoints(receipt.Total)
    points += calculateQuarterPoints(receipt.Total)
    points += calculateItemCountPoints(receipt.Items)
    points += calculateItemDescriptionPoints(receipt.Items)
    points += calculateOddDayPoints(receipt.PurchaseDate)
    points += calculateHappyHourPoints(receipt.PurchaseTime)
	
	logger.Printf("Total points calculated for receipt: %d", points)
    
	return points
}

func calculateRetailerNamePoints(retailer string) int64 {
    matches := alphanumericRegex.FindAllString(retailer, -1)
	points := int64(len(matches))
	logger.Printf("Retailer name '%s' earned %d points for %d alphanumeric characters", 
        retailer, points, len(matches))
    return points
}

func calculateRoundDollarPoints(total string) int64 {
    if strings.HasSuffix(total, ".00") {
        logger.Printf("Round dollar amount found: %s", total)
        return RoundDollarPoints
    }
    return 0
}

func calculateQuarterPoints(total string) int64 {
    if amount, err := strconv.ParseFloat(total, 64); err == nil {
        if math.Mod(amount*100, 25) == 0 {
            logger.Printf("Quarter dollar amount found: %s", total)
            return QuarterDollarPoints
        }
    }
    return 0
}

func calculateItemCountPoints(items []models.Item) int64 {
    points := (len(items) / 2) * ItemPairPoints
    logger.Printf("Item count points: %d for %d items", points, len(items))
    return points
}

func calculateItemDescriptionPoints(items []models.Item) int64 {
    var points int64
    for _, item := range items {

        trimLen := len(strings.TrimSpace(item.ShortDescription))
        
		if trimLen%3 == 0 {
            
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				itemDescriptionPoints := int64(math.Ceil(price * 0.2))
                logger.Printf("Item '%s' earned %d points (description length %d is divisible by 3)", item.ShortDescription, itemDescriptionPoints, trimLen)
				points += itemDescriptionPoints
            }
        }
    }
	logger.Printf("Total points from item descriptions: %d", points)
    return points
}

func calculateOddDayPoints(purchaseDate string) int64 {
	if day, err := strconv.Atoi(purchaseDate[8:]); err == nil {
        if day%2 == 1 {
            logger.Printf("Odd day points awarded for day: %d", day)
            return OddDayPoints
        }
    }
    return 0
}

func calculateHappyHourPoints(purchaseTimeString string) int64 {
	if purchaseTime, err := time.Parse("15:04", purchaseTimeString); err == nil {
        startTime := time.Date(0, 1, 1, HappyHourStart, 0, 0, 0, time.UTC)
        endTime := time.Date(0, 1, 1, HappyHourEnd, 0, 0, 0, time.UTC)
        checkTime := time.Date(0, 1, 1, purchaseTime.Hour(), purchaseTime.Minute(), 0, 0, time.UTC)

        if checkTime.After(startTime) && checkTime.Before(endTime) {
            logger.Printf("Happy hour points awarded for time: %s", purchaseTimeString)
            return HappyHourPoints
        }
    }
    return 0
}
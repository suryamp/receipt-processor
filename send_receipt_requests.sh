#!/bin/bash

# Array of retailers
retailers=("Target" "Walmart" "Best Buy" "Costco" "Whole Foods")

# Array of item names
items=(
  "Milk" "Bread" "Eggs" "Cheese" "Chicken" 
  "Rice" "Pasta" "Tomato Sauce" "Cereal" "Coffee"
  "Bananas" "Apples" "Orange Juice" "Yogurt" "Butter"
  "Chips" "Cookies" "Water Bottle" "Energy Drink" "Chocolate Bar"
)

# Function to generate random receipt data
generate_receipt() {
  local retailer="${retailers[$RANDOM % ${#retailers[@]}]}"
  local purchaseYear=$((RANDOM % 5 + 2020))
  local purchaseMonth=$((RANDOM % 12 + 1))
  local purchaseDay=$((RANDOM % 28 + 1))
  local purchaseDate="$purchaseYear-$(printf "%02d" $purchaseMonth)-$(printf "%02d" $purchaseDay)"
  local purchaseTime=$(printf "%02d:%02d" $((RANDOM % 24)) $((RANDOM % 60)))
  local itemCount=$((RANDOM % 5 + 1))
  local total=0

  echo "{" 
  echo "  \"retailer\": \"$retailer\","
  echo "  \"purchaseDate\": \"$purchaseDate\","
  echo "  \"purchaseTime\": \"$purchaseTime\","
  echo "  \"items\": ["

  for ((i=0; i<itemCount; i++)); do
    local shortDesc="${items[$RANDOM % ${#items[@]}]}"
    local price=$(echo "scale=2; $RANDOM / 100" | bc)
    total=$(echo "scale=2; $total + $price" | bc)

    echo "    {"
    echo "      \"shortDescription\": \"$shortDesc\","
    echo "      \"price\": \"$price\""
    
    # Add comma only if it's not the last item
    if ((i < itemCount - 1)); then
      echo "    },"
    else
      echo "    }"
    fi
  done

  echo "  ],"
  echo "  \"total\": \"$total\""
  echo "}"
}

# Send POST requests in a loop
while true; do
  # Generate receipt data
  receipt_data=$(generate_receipt)
  
  # Process receipt and extract UUID
  receipt_uuid=$(curl -s -X POST http://localhost:8080/receipts/process \
    -H "Content-Type: application/json" \
    -d "$receipt_data" | jq -r '.id')

  echo "Sent request with data: $receipt_data"

  echo "Receipt UUID: $receipt_uuid"
  
  # Get points for the receipt
  points_response=$(curl -s GET "http://localhost:8080/receipts/$receipt_uuid/points")
  
  echo "Points for receipt: $points_response"
  
  # Random delay between requests
  sleep $((RANDOM % 5 + 1))
done
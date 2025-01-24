# Receipt Processor

A REST API service that processes receipts and calculates points based on specific rules.

## Features
- Process receipts and calculate points based on various rules
- RESTful API endpoints
- Prometheus metrics
- Grafana dashboards
- Docker containerization
- Health monitoring
- Comprehensive logging

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or higher (for local development)
- cURL (for testing)
- jq (for testing script)

### Quick Start with Docker

1. Clone the repository:
```bash
git clone https://github.com/suryamp/receipt-processor.git
cd receipt-processor
```

2. Start all services using Docker Compose:
```bash
docker-compose up -d
```

This will start:
- Receipt Processor service (port 8080)
- Prometheus (port 9090)
- Grafana (port 3000)

### Local Development Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

## API Documentation

### Health Check
Check if the service is healthy.

**Endpoint:** `GET /health`

```bash
curl -X GET http://localhost:8080/health
```

**Success Response (200 OK):**
```
OK
```

### Process Receipt
Process a receipt and get back an ID.

**Endpoint:** `POST /receipts/process`

```bash
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "Target",
    "purchaseDate": "2024-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew",
        "price": "1.25"
      }
    ],
    "total": "1.25"
  }'
```

### Get Points
Get points for a receipt.

**Endpoint:** `GET /receipts/{id}/points`

```bash
curl -X GET http://localhost:8080/receipts/{id}/points
```

## Monitoring

### Prometheus Metrics
Available at `http://localhost:8080/metrics`

Key metrics:
- `http_requests_total`: Total number of HTTP requests
- `http_request_duration_seconds`: Duration of HTTP requests

### Grafana Dashboards
Access Grafana at `http://localhost:3000`

Default dashboards include:
- Request Rate & Durations

### Debug Endpoints (To Be Implemented)
Available in development:
- `/debug/pprof/`: Index of pprof endpoints
- `/debug/pprof/goroutine`: Current goroutines
- `/debug/pprof/heap`: Heap profile
- `/debug/pprof/profile`: CPU profile

## Testing

### Run Unit Tests
```bash
go test ./...
```

### Testing Script
The repository includes a script for testing:

```bash
chmod +x send_receipt_requests.sh
./send_receipt_requests.sh
```

## Points Calculation Rules

Points are awarded based on the following rules:

1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
6. 6 points if the day in the purchase date is odd
7. 10 points if the time of purchase is between 2:00pm and 4:00pm

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

# Future Features

## Business/Functional Features

- **Internationalization**
  - Support for receipts that are priced in non-USD currencies
  - Multi-language support for user interface

- **Receipt Processing**
  - Image processing with OCR for data extraction
  - Bulk receipt processing capabilities
  - Advanced receipt validation and verification

- **User Management**
  - User accounts and authentication
  - Receipt history tracking
  - Points expiration system with notifications
  - Receipt categories and tagging

- **Search and Analytics**
  - Advanced receipt search and filtering
  - Monthly/weekly points summaries
  - Custom reporting capabilities

- **Business Rules**
  - Configurable point calculation rules engine
  - Special promotions and bonus points
  - Retailer-specific rules and bonuses

- **Fraud Prevention**
  - Duplicate receipt detection
  - Suspicious pattern identification
  - Transaction verification system

## Technical Features

- **API Enhancements**
  - GraphQL API alongside REST
  - WebSocket endpoints for real-time updates
  - API versioning and deprecation management

- **Performance Optimization**
  - Caching layer (Redis/Memcached)
  - Request ID tracking
  - Performance monitoring and alerting

- **Infrastructure**
  - Persistent storage solution
  - Multi-region deployment support
  - Service mesh integration
  - Feature flag management

- **Scalability**
  - Horizontal scaling capabilities
  - Load balancing
  - Rate limiting implementation

## Data Analysis Features

- **Analytics Dashboard**
  - Shopping pattern analysis
  - Retailer insights and metrics
  - Performance analytics
  - Usage pattern tracking

- **User Insights**
  - Points optimization suggestions
  - Spending pattern analysis
  - Receipt submission trends

- **System Analytics**
  - Anomaly detection
  - System performance metrics
  - API usage patterns

## Security Features

- **Data Protection**
  - Input sanitization middleware
  - Data encryption at rest
  - Detailed access logging
  - GDPR/privacy compliance

- **Access Control**
  - API key management
  - Role-based access control
  - Rate limiting
  - Authentication/Authorization

## Developer Experience

- **Documentation**
  - Interactive API documentation (Swagger/OpenAPI)
  - Developer portal
  - SDK generation for multiple languages
  - Comprehensive API examples

- **Testing Tools**
  - Postman collection
  - Load testing framework

- **Development Environment**
  - Local development setup scripts
  - Docker development environment
  - CI/CD pipeline configuration
  - Development guidelines and best practices
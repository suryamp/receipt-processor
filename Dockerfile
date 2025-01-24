# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o receipt-processor

# Final stage
FROM alpine:3.18
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache curl

# Copy binary from builder
COPY --from=builder /build/receipt-processor .

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./receipt-processor"]
# Receipt Processor Service - SRE Playbook

## Service Overview

### Basic Information
- **Service Name**: Receipt Processor
- **Port**: 8080

### Endpoints
- POST `/receipts/process`
- GET `/receipts/{id}/points`

### Dependencies
- Prometheus (Port 9090)
- Grafana (Port 3000)

## Health Monitoring

### Key Metrics to Monitor

#### 1. HTTP Metrics
- **Request rate**
- **Error rate** (4xx, 5xx)
- **Response time**
  - p50
  - p90
  - p99
- **Request payload size**

#### 2. Application Metrics
- Memory usage
- CPU usage
- Goroutine count
- In-memory receipt count

#### 3. Business Metrics
- Points calculation rate
- Average points per receipt
- Receipt processing success rate

### Grafana Details
- **Access**: `http://localhost:3000`
- **Default Credentials**: Anonymous access enabled
- **Default Dashboards**: `/var/lib/grafana/dashboards`

### Alert Thresholds (To Be Implemented)

#### High Severity Alerts
- Error rate > 5% over 5 minutes
- Response time p99 > 1s over 5 minutes
- Memory usage > 90%

#### Medium Severity Alerts
- Error rate > 2% over 15 minutes
- Response time p90 > 500ms over 15 minutes
- Memory usage > 80%

### Container Health Commands

```bash
# Check container status
docker-compose ps

# View container logs
docker-compose exec receipt-processor tail -f /app/logs/receipt-processor.log
```

## Troubleshooting Guide

### High Error Rate

1. **Check logs for error patterns**
   ```bash
   docker-compose exec receipt-processor tail -f /app/logs/receipt-processor.log
   ```

2. **Common Error Scenarios**

   #### Invalid Receipt Format
   ```json
   {
     "error": "The receipt is invalid",
     "details": "Check request payload against API schema"
   }
   ```

   #### Receipt Not Found
   ```json
   {
     "error": "No receipt found for that ID",
     "details": "Verify ID exists in memory store"
   }
   ```

### High Response Time

1. **Check system resources**
   ```bash
   docker stats $(docker-compose ps -q receipt-processor)
   ```

2. **Profile the application** (To Be Implemented)
   ```bash
   curl http://localhost:8080/debug/pprof/heap > heap.pprof
   go tool pprof heap.pprof
   ```

## Incident Response

### Service is Down

1. **Check system logs**
   ```bash
   docker-compose exec receipt-processor tail -f /app/logs/receipt-processor.log
   ```

2. **Restart service**
   ```bash
   docker-compose restart receipt-processor
   ```

3. **Verify service is up**
   ```bash
   curl -f http://localhost:8080/health
   ```

### Container Issues

1. **Rebuild container**
   ```bash
   docker-compose build receipt-processor
   docker-compose up -d receipt-processor
   ```

2. **Check container health**
   ```bash
   docker-compose ps receipt-processor
   ```

### Data Recovery

- No direct data recovery is possible
- Service restart will clear all receipts
- Clients need to resubmit receipts if data is lost

## Deployment

### Pre-deployment Checklist

1. **Run tests**
   ```bash
   go test ./...
   ```

2. **Build and test Docker image**
   ```bash
   docker-compose build
   docker-compose up -d
   ```

3. Verify API contract hasn't changed

### Post-deployment Verification

1. **Check health endpoint**
   ```bash
   curl -f http://localhost:8080/health
   ```

2. **Submit test receipt**
   ```bash
   curl -X POST http://localhost:8080/receipts/process \
     -H "Content-Type: application/json" \
     -d '{
       "retailer": "Test",
       "purchaseDate": "2024-01-23",
       "purchaseTime": "13:01",
       "items": [
         {
           "shortDescription": "Test Item",
           "price": "10.00"
         }
       ],
       "total": "10.00"
     }'
   ```

3. **Verify points calculation**
   ```bash
   curl -X GET http://localhost:8080/receipts/{id}/points
   ```

4. Verify metrics in Grafana

## Scaling Considerations

### Memory Usage
- Each receipt consumes approximately:
  - Base receipt: ~200 bytes
  - Per item: ~100 bytes
- Consider this for capacity planning

### Resource Requirements
- **Minimum**: 256MB RAM, 0.5 CPU
- **Recommended**: 512MB RAM, 1 CPU
- **Storage**: Minimal (in-memory storage)

### Performance Limits
- Single instance recommended limits:
  - Max 1000 requests/second
  - Max 100,000 receipts in memory
  - Max 100 concurrent connections

### Scaling Strategies

1. **Vertical Scaling**
   - Increase container resources

2. **Horizontal Scaling**
   - Requires persistent storage implementation
   - Load balancer configuration

## Maintenance

### Routine Checks

1. **Daily**:
   - Monitor error rates
   - Check memory usage trend

2. **Weekly**:
   - Review performance metrics
   - Check for goroutine leaks

3. **Monthly**:
   - Load test service
   - Review scaling needs

### Regular Maintenance

1. **Log Rotation**
   ```bash
   /etc/logrotate.d/receipt-processor
   ```

2. **Memory Cleanup**
   - Service restart during low traffic
   - Implement TTL for old receipts

## Contact Information

### Escalation Path
1. Surya Manchikanti

### Important Links
- **Monitoring Dashboard**: `http://localhost:3000`
- **Metrics**: `http://localhost:9090`
- **Logs**: `docker-compose logs receipt-processor`
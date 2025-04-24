# E-commerce Search Service

A microservice for handling product search functionality in an e-commerce application. This service provides APIs for managing products and implements a search engine using Elasticsearch.

## Features

- RESTful API for product management (CRUD operations)
- Full-text search using Elasticsearch
- Event-driven architecture using Kafka
- MongoDB for data persistence
- Clean architecture implementation
- Docker support for easy deployment

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- MongoDB
- Elasticsearch
- Kafka

## Project Structure

```
.
├── cmd/
│   ├── api/         # API service
│   └── worker/      # Index update worker
├── config/          # Configuration files
├── internal/        # Internal packages
│   ├── config/      # Configuration
│   ├── domain/      # Domain models and interfaces
│   ├── repository/  # Repository implementations
│   ├── service/     # Business logic
│   └── delivery/    # Delivery layer (HTTP handlers)
├── pkg/            # Shared packages
├── test/           # Test files
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ecommerce-search-svc.git
   cd ecommerce-search-svc
   ```

2. Copy the example configuration:
   ```bash
   cp config/config.example.yaml config/config.yaml
   ```

3. Update the configuration in `config/config.yaml` with your settings.

4. Install dependencies:
   ```bash
   make setup
   ```

## Kafka Topics Setup

After starting the Kafka service, create the required topics:

```bash
# Create topic for product updates
docker exec -it ecommerce-search-svc-kafka-1 kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic product-updates \
  --partitions 3 \
  --replication-factor 1

# Create topic for search events
docker exec -it ecommerce-search-svc-kafka-1 kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic search-events \
  --partitions 3 \
  --replication-factor 1

# List all topics to verify
docker exec -it ecommerce-search-svc-kafka-1 kafka-topics --list \
  --bootstrap-server localhost:9092
```

## Running the Service

### Using Docker Compose

1. Build and start all services:
   ```bash
   make docker-build
   make docker-run
   ```

### Running Locally

1. Start the API service:
   ```bash
   make build
   make run
   ```

2. Start the worker:
   ```bash
   make run-worker
   ```

## API Endpoints

### Create Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Latest Apple smartphone with A17 Pro chip",
    "price": 999.99,
    "category": "Electronics",
    "tags": ["smartphone", "apple", "iphone"],
    "brand": "Apple",
    "views": 0,
    "buys": 0
  }'
```

### Update Product
```bash
curl -X PUT http://localhost:8080/products/123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Latest Apple smartphone with A17 Pro chip",
    "price": 899.99,
    "category": "Electronics",
    "tags": ["smartphone", "apple", "iphone"],
    "brand": "Apple",
    "views": 100,
    "buys": 50
  }'
```

### Delete Product
```bash
curl -X DELETE http://localhost:8080/products/123
```

### Get Product by ID
```bash
curl -X GET http://localhost:8080/products/123
```

### Search Products
```bash
curl -X GET "http://localhost:8080/products/search?q=iphone&category=Electronics&min_price=500&max_price=1000&page=1&limit=10"
```

### Increment Product Views
```bash
curl -X POST http://localhost:8080/products/123/views
```

### Increment Product Buys
```bash
curl -X POST http://localhost:8080/products/123/buys
```

Note: The search endpoint supports the following query parameters:
- `q`: Search query string
- `category`: Filter by category
- `min_price`: Minimum price filter
- `max_price`: Maximum price filter
- `page`: Page number for pagination (default: 1)
- `limit`: Number of items per page (default: 10)

## Testing

Run the tests:
```bash
make test
```

## Development

1. Install development dependencies:
   ```bash
   go mod download
   ```

2. Run the service in development mode:
   ```bash
   go run cmd/api/main.go
   ```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
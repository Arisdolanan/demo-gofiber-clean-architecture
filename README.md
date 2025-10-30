# GoFiber Clean Architecture with PostgreSQL, Redis, Kafka, and RabbitMQ

This project demo a clean architecture implementation using GoFiber framework with PostgreSQL, Redis, Kafka, and RabbitMQ integrations.


## Prerequisites

- Go 1.23.5
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL
- Redis
- RabbitMQ
- Kafka and Zookeeper

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd demo-gofiber-clean-architecture
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Copy the example configuration file:
   ```bash
   cp config-example.json config.json
   ```

4. Update `config.json` with your configuration values.

## Configuration

The application uses a `config.json` file for configuration. Copy `config-example.json` to `config.json` and update the values as needed.

## Running the Application

### Development Mode

To run the application in development mode with hot-reloading:

```bash
make dev
```

### Production Mode

To build and run the application in production mode:

```bash
make build
make run
```

Or run directly with Go:

```bash
go run cmd/web/main.go
```

### Running the Worker Service

To run the worker service for consuming messages:

```bash
make worker
```

### Running Tests

To run all tests:

```bash
make test
```

To run unit tests only:

```bash
make test-unit
```

To generate a coverage report:

```bash
make test-coverage-html
```

### Database Migrations

To apply database migrations:

```bash
make migrate-up
```

To rollback database migrations:

```bash
make migrate-down
```

To create a new migration:

```bash
make migrate-create name=migration_name
```

## API Documentation

The API documentation is available via Swagger UI at `http://localhost:3000/swagger/index.html`.

## Docker Setup

This project includes Docker configuration for easy deployment and development.

### Services Included

- **PostgreSQL**: Database for storing user and application data
- **Redis**: In-memory cache for improved performance
- **RabbitMQ**: Message broker for asynchronous communication
- **Kafka**: Distributed streaming platform for high-volume data processing
- **GoFiber Application**: Main application service

#### Common Commands

To stop all services:

```bash
docker-compose down
```

To view logs:

```bash
docker-compose logs -f
```

### Accessing Services

- **Application**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **Kafka**: localhost:9092

### Messaging Infrastructure

The project now includes messaging infrastructure for both RabbitMQ and Kafka:
- **Kafka**: Using IBM sarama library with consumer group support
- **RabbitMQ**: Using rabbitmq/amqp091-go library (updated from deprecated streadway/amqp)

Connection management for both messaging systems is handled in the `internal/infrastructure/messaging` directory.

#### Kafka Implementation

The Kafka implementation includes both producer and consumer functionality with support for consumer groups:

1. **Producer**: Send messages to Kafka topics
   - Located in `internal/infrastructure/messaging/kafka.go` (infrastructure layer)
   - Exposed via `internal/delivery/messaging/kafka/kafka_producer.go` (delivery layer)

2. **Consumer**: Receive and process messages from Kafka topics
   - Located in `internal/infrastructure/messaging/kafka.go` (infrastructure layer)
   - Exposed via `internal/delivery/messaging/kafka/kafka_consumer.go` (delivery layer)
   - Example consumer implementation in `cmd/worker/main.go`
   - Supports both partition-based consumption and consumer groups

To use the Kafka consumer, run the worker service:
```bash
go run cmd/worker/main.go
```

#### RabbitMQ Implementation

The RabbitMQ implementation includes both producer and consumer functionality:

1. **Producer**: Send messages to RabbitMQ queues
   - Located in `internal/infrastructure/messaging/rabbitmq.go` (infrastructure layer)
   - Exposed via `internal/delivery/messaging/rabbitmq/rabbitmq_producer.go` (delivery layer)

2. **Consumer**: Receive and process messages from RabbitMQ queues
   - Located in `internal/infrastructure/messaging/rabbitmq.go` (infrastructure layer)
   - Exposed via `internal/delivery/messaging/rabbitmq/rabbitmq_consumer.go` (delivery layer)
   - Example consumer implementation in `cmd/worker/main.go`

To use the RabbitMQ consumer, run the worker service:
```bash
go run cmd/worker/main.go
```

## File Storage

- **All PDFs**: Stored in `./storage/private/pdfs/` with unique filenames to prevent conflicts
- The URL to access the PDF is returned in the response, but files are not accessible via HTTP
- To access generated PDFs, you need to implement internal access mechanisms in your application

## Logging

The application supports multiple logging configurations:

1. **Console Only**
2. **File Only**
3. **Console + Single File**
4. **Console + Level-separated Files**

See [Logging Configuration](internal/config/README.md) for more details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request
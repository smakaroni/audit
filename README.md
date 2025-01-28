# Audit Log Service

This service acts as an audit log, picking up Kafka messages, anonymizing the data, and saving it to a PostgreSQL database.

## Features

- Consumes messages from Kafka
- Anonymizes sensitive data
- Stores anonymized data in PostgreSQL

## Getting Started

### Prerequisites

- Go 1.23+
- Kafka
- PostgreSQL

### Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Set up your Kafka and PostgreSQL connections in the configuration
4. Run the service: `go run cmd/main.go`

## Testing

Run tests with: `go test ./...`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
# Audit Log Service

This service acts as an audit log, picking up Kafka messages, and saving them to a MariaDB database using Protocol Buffers for serialization.

## Features

- Consumes messages from Kafka
- Uses Protocol Buffers for efficient data serialization
- Stores serialized data in MariaDB

## Getting Started

### Prerequisites

- Go 1.23+
- Kafka
- MariaDB
- Protocol Buffers compiler (protoc)

### Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Install protoc plugins: `make get-protoc-plugins`
4. Generate protos: `make proto-gen`
5. Set up your Kafka and MariaDB connections in the configuration
6. Run the service: `go run cmd/main.go`

## Protocol Buffers

This project uses Protocol Buffers for efficient serialization of AuditLog data. The AuditLog message is defined in [internal/protos/audit_log.proto](cci:7://file:///Users/jokke/GolandProjects/audit/internal/protos/audit_log.proto:0:0-0:0):

```protobuf
syntax = "proto3";
package models;

import "google/protobuf/timestamp.proto";

message AuditLog {
  string event_type = 1;
  string anonymized_user = 2;
  google.protobuf.Timestamp timestamp = 3;
  string data = 4;
}
```

## Testing

Run tests with: `go test ./...`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
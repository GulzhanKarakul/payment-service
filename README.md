# Payment Service

A production-ready payment processing microservice built in Go. Designed as the foundation 
for a real fintech platform with atomic transactions, comprehensive testing, and enterprise-grade practices.

[![Go Report Card](https://goreportcard.com/badge/github.com/GulzhanKarakul/payment-service)](https://goreportcard.com/report/github.com/GulzhanKarakul/payment-service)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

✅ **Atomic Transactions**
- Multi-step PostgreSQL transactions with ACID guarantees
- Automatic rollback on failure (no partial updates)
- Handles concurrent requests safely

✅ **Comprehensive Testing**
- 90%+ code coverage
- testcontainers for realistic database testing
- mockery for service layer unit tests
- go test -race for concurrency bug detection

✅ **Production Ready**
- Structured JSON logging with request tracing
- Graceful shutdown with signal handling
- Panic recovery middleware
- Proper error handling and mapping

✅ **REST API**
- Chi HTTP router with layered architecture
- 8 endpoints for transactions, clients, businesses
- Input validation and error responses
- Pagination support

✅ **Easy Deployment**
- Docker with multi-stage builds
- docker-compose for local development
- Environment-based configuration
- Health check endpoint

---

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & docker-compose
- PostgreSQL 15 (or use docker-compose)

### Installation

1. **Clone the repository**
```bash
   git clone https://github.com/GulzhanKarakul/payment-service.git
   cd payment-service
```

2. **Setup environment**
```bash
   cp .env.example .env
```

3. **Start with Docker**
```bash
   docker-compose up -d
```

4. **Build and run**
```bash
   go build -o payment-service ./cmd
   ./payment-service
```

Server runs on `http://localhost:8080`

---

## Architecture

### Layered Architecture
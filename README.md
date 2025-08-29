# anyompt - LLMs integration API

This project provides an API that integrates multiple large language models (LLM).

## Features

- Send prompts to OpenAI or Groq from a single endpoint.
- Dynamically switch the model used without modifying client code.
- Scale and extend to other LLMs in the future.
- Event-driven integration with Gateway: Optional, send responses to a Gateway for further processing.

Currently, it is integrated with OpenAI and Groq. Groq offers multiple free models with certain token limits; see
documentation at: [Groq](https://console.groq.com/docs/overview)

### Prerequisites

- Go 1.21 or higher
- OpenAI API key (optional, for OpenAI integration)
- Groq API key (optional, for Groq integration)
- Gateway

## 🚀 Installation

1. Install dependencies:

```bash
go mod tidy
```

2. Configure environment variables:

```bash
cp env.example .env
# Edit .env with the values described below.
```

3. Run the application:

```bash
go run main.go
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file based on `env.example`:

- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Log level (debug, info, warn, error, fatal, panic - default: info)
- `CHAT_MODEL`: Chat model to use. If "OpenAI" is selected, the OpenAI API is used; otherwise, Groq is used.
    - Example for Groq: llama-3.3-70b-versatile
    - Default: openai/gpt-oss-20b
- `OPENAI_API_KEY`: OpenAI API key (required for OpenAI)
- `GROQ_API_KEY`: Groq API key (required for Groq)
- `GROQ_URL`: Groq API URL (default: https://api.groq.com/openai/v1/responses)
- `GATEWAY_URL`: Gateway API URL (optional)

### OpenAI API Setup

1. **Get OpenAI API Access:**
  - Create an OpenAI account
  - Create an API Token

### Groq API Setup

1. **Get Groq API Access:**
   - Create a Groq account
   - Create an API Token

## 📡 Endpoints

### POST /api/v1/chat/ask

Sends a prompt to the selected LLM and receives a response.

**Request:**

```json
{
  "prompt": "What is the capital of France?"
}
```

**Response:**

```json
{
  "response": "The capital of France is Paris."
}
```

### GET /health

Checks the API status.

**Response:**

```json
{
  "status": "OK",
  "message": "anyompt API is running"
}
```

#### Using curl:

```bash
# Health check
curl http://localhost:8080/health

# Chat endpoint
curl -X POST http://localhost:8080/api/v1/chat/ask \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: f81d4fae-7dec-11d0-a765-00a0c91e6bf6" \
  -H "X-Routing-Key: telegram:12345" \
  -d '{"prompt": "What is the capital of France?"}'
```

## 🎗️ Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities, repository interfaces, and use cases
- **Application**: Implementation of use cases
- **Infrastructure**: OpenAI and Groq repository implementations
- **Interfaces**: HTTP controllers and routers

## 📁 Project Structure

```
anyompt/
├── cmd/                  # Application entry points
│   └── server/           # Main server
├── internal/             # Project-specific code
│   ├── config/           # Configuration
│   ├── infrastructure/   # Repository implementations
│   └── interfaces/       # HTTP controllers
│       ├── http/         # Handler controller
│       └── middleware/   # Middlewares
├── pkg/                  # Reusable and public code
│   ├── domain/           # Domain entities and interfaces
│   └── application/      # Use cases
├── main.go               # Main entry point
├── go.mod                # Go dependencies
├── README_ES.md          # README in spanish
└── README.md             # This file
```

## 🧪 Testing

### Running Tests

To run all tests:

```bash
go test ./...
```

### Test Coverage

To check test coverage (excluding mocks):

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage excluding mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

### Running Benchmarks

```bash
go test -bench=. ./...
```

## BackLog

- [x] Unit Tests
- [ ] Add others paid LLMs
- [ ] Integration tests
- [ ] API documentation with Swagger
- [ ] Add request_id in header and its middleware

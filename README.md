# AnyPrompt API

This project provides an API that integrates multiple large language models (LLM).

## Features

- Send prompts to OpenAI or Groq from a single endpoint.
- Dynamically switch the model used without modifying client code.
- Scale and extend to other LLMs in the future.

Currently, it is integrated with OpenAI and Groq. Groq offers multiple free models with certain token limits; see
documentation at: [Groq](https://console.groq.com/docs/overview)

### Prerequisites

- Go 1.21 or higher
- WhatsApp Business API access token
- WhatsApp Business phone number

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

- `CHAT_MODEL`: Chat model to use. If "OpenAI" is selected, the OpenAI API is used; otherwise, Groq is used.
    - Example for Groq: llama-3.3-70b-versatile
- `OPENAI_API_KEY`: OpenAI API key (required for OpenAI)
- `GROQ_API_KEY`: Groq API key (required for Groq)
- `PORT`: Server port (default: 8080)

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
  "message": "AnyPrompt API is running"
}
```

#### Using curl:

```bash
# Health check
curl http://localhost:8080/health

# Chat endpoint
curl -X POST http://localhost:8080/api/v1/chat/ask \
  -H "Content-Type: application/json" \
  -d '{"prompt": "What is the capital of France?"}'
```

## 🎗️ Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities, repository interfaces, and use cases
- **Application**: Implementation of use cases
- **Infrastructure**: OpenAI repository implementation
- **Interfaces**: HTTP controllers and routers

## 📁 Project Structure

```
anyprompt/
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
└── README.md             # This file
```

## BackLog

- [ ] Unit Tests
- [ ] Add others paid LLMs
# anyompt - API de integración con LLMs

Este proyecto provee una API que integra múltiples modelos de lenguaje grandes (LLM).

## Características

- Enviar prompts a OpenAI o Groq desde un mismo endpoint.
- Cambiar dinámicamente el modelo utilizado sin modificar el código cliente.
- Escalar y extender a otros LLMs en el futuro.
- Integración orientada a eventos con Kafka: Opcional, envía las respuestas a un tópico de Kafka para su posterior procesamiento.

Por el momento está integrada con OpenAI y con Groq, este último permite múltiples modelos gratuitos con cierto límite de token, ver documentación en: [Groq](https://console.groq.com/docs/overview)

### Prerequisitos

- Go 1.21 o mayor
- API key de OpenAI (opcional, para integración con OpenAI)
- API key de Groq (opcional, para integración con Groq)
- Kafka (opcional, para integración con Kafka)

## 🚀 Instalación

1. Instalar dependencias:

```bash
go mod tidy
```

2. Configurar las variables de entorno:

```bash
cp env.example .env
# Editar .env con los valores descriptos debajo.
```

3. Ejecutar la aplicación:

```bash
go run main.go
```

## 🔧 Configuración

### Variables de Entorno

Crear un archivo `.env` basado en `env.example`:

- `CHAT_MODEL`: Modelo de chat a usar. Si se selecciona "OpenAI", se usa la API de OpenAI; de lo contrario, se usa Groq.
    - Ejemplo para Groq: llama-3.3-70b-versatile
    - Por defecto: openai/gpt-oss-20b
- `OPENAI_API_KEY`: API key de OpenAI (requerida para OpenAI)
- `GROQ_API_KEY`: API key de Groq (requerida para Groq)
- `GROQ_URL`: URL de la API de Groq (por defecto: https://api.groq.com/openai/v1/responses)
- `PORT`: Puerto del servidor (por defecto: 8080)
- `LOG_LEVEL`: Nivel de log (debug, info, warn, error, fatal, panic - por defecto: info)
- `KAFKA_ENABLED`: Habilita la integración con Kafka (true o false)
- `KAFKA_BROKERS`: Lista de brokers de Kafka separados por comas
- `KAFKA_TOPIC`: Tópico de Kafka para enviar eventos

### Configuración de OpenAI API

1. **Obtener acceso a la API de OpenAI:**
  - Crear una cuenta en OpenAI
  - Crear un Token de API

### Configuración de Groq API

1. **Obtener acceso a la API de Groq:**
   - Crear una cuenta en Groq
   - Crear un Token de API

## 📡 Endpoints

### POST /api/v1/chat/ask

Envía un prompt al LLM elegido y recibe una respuesta.

**Request:**
```json
{
  "prompt": "¿Cuál es la capital de Francia?"
}
```

**Response:**
```json
{
  "response": "La capital de Francia es París."
}
```

### GET /health

Verifica el estado de la API.

**Response:**
```json
{
  "status": "OK",
  "message": "anyompt API is running"
}
```

#### Usando curl:

```bash
# Health check
curl http://localhost:8080/health

# Chat endpoint
curl -X POST http://localhost:8080/api/v1/chat/ask \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Cuál es la capital de Francia?"}'
```

## 🎗️ Arquitectura

Este proyecto sigue los principios de Clean Architecture:

- **Domain**: Entidades, interfaces de repositorio y casos de uso
- **Application**: Implementación de casos de uso
- **Infrastructure**: Implementaciones de repositorios de OpenAI y Groq
- **Interfaces**: Controladores HTTP y routers

## 📁 Estructura del Proyecto

```
anyompt/
├── cmd/                  # Puntos de entrada de la aplicación
│   └── server/           # Servidor principal
├── internal/             # Código específico del proyecto
│   ├── config/           # Configuración
│   ├── infrastructure/   # Implementaciones de repositorios
│   └── interfaces/       # Controladores HTTP
│       ├── http/         # Controlador handler
│       └── middleware/   # Middlewares
├── pkg/                  # Código reutilizable y público
│   ├── domain/           # Entidades e interfaces del dominio
│   └── application/      # Casos de uso
├── main.go               # Punto de entrada principal
├── go.mod                # Dependencias de Go
├── README_ES.md          # Este archivo
└── README.md             # README en inglés
```

## 🧪 Pruebas

### Ejecutar Pruebas

Para ejecutar todas las pruebas:

```bash
go test ./...
```

Para ejecutar pruebas con salida detallada:

```bash
go test -v ./...
```

Para ejecutar pruebas de un paquete específico:

```bash
go test ./internal/config/
go test ./cmd/server/
```

### Cobertura de Pruebas

Para verificar la cobertura de pruebas (excluyendo mocks):

```bash
# Generar reporte de cobertura
go test -coverprofile=coverage.out ./...

# Ver reporte de cobertura en terminal
go tool cover -func=coverage.out

# Generar reporte HTML de cobertura
go tool cover -html=coverage.out -o coverage.html

# Ver cobertura excluyendo mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

### Ejecutar Benchmarks

```bash
go test -bench=. ./...
```

## BackLog

- [x] Pruebas unitarias
- [ ] Agregar otros LLMs de pago
- [ ] Pruebas de integración
- [ ] Documentación de API con Swagger


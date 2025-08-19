# AnyPrompt API

Este proyecto provee una API que integra distintos modelos de lenguaje, permitiendo:

- Enviar prompts a OpenAI o Groq desde un mismo endpoint.
- Cambiar dinámicamente el modelo utilizado sin modificar el código cliente.
- Escalar y extender a otros LLMs en el futuro.

Por el momento está integrada con OpenAI y con Groq, este último permite múltiples modelos gratuitos con cierto límite de token, ver documentación en: [Groq](https://console.groq.com/docs/overview)

## 🚀 Instalación

1. Instalar dependencias:
```bash
go mod tidy
```

3. Configurar las variables de entorno:
```bash
cp env.example .env
# Editar .env con los valores descriptos debajo.
```

4. Ejecutar la aplicación:
```bash
go run main.go
```

## 🔧 Configuración

### Variables de Entorno

- `CHAT_MODEL`: El modelo de chat a usar, si se elige "OpenAI", utiliza la OpenAI API, sino, utiliza Groq.
  - ejemplo con Groq: llama-3.3-70b-versatile.
- `OPENAI_API_KEY`: API key de OpenAI (requerida para usar OpenAI)
- `GROQ_API_KEY`: API key de Groq (requerida para usar Groq) 
- `PORT`: Puerto del servidor (por defecto: 8080)
- `MOCK_MODE`: Modo mock (true) devolverá una respuesta simulada. 

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
  "message": "AnyPrompt API is running"
}
```

## 🧪 Pruebas

### Con curl:

```bash
# Health check
curl http://localhost:8080/health

# Chat endpoint
curl -X POST http://localhost:8080/api/v1/chat/ask \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Cuál es la capital de Francia?"}'
```

## 🏗️ Arquitectura

Este proyecto sigue los principios de Clean Architecture:

- **Domain**: Entidades, interfaces de repositorio y casos de uso
- **Application**: Implementación de casos de uso
- **Infrastructure**: Implementación del repositorio OpenAI
- **Interfaces**: Controladores HTTP y routers

## 📁 Estructura del Proyecto

```
anyprompt/
├── cmd/                  # Puntos de entrada de la aplicación
│   └── server/           # Servidor principal
├── pkg/                  # Código reutilizable y público
│   ├── domain/           # Entidades e interfaces del dominio
│   └── application/      # Casos de uso
├── internal/             # Código específico del proyecto
│   ├── config/           # Configuración
│   ├── infrastructure/   # Implementaciones de repositorios
│   └── interfaces/       # Controladores HTTP
├── main.go               # Punto de entrada principal
├── go.mod                # Dependencias de Go
└── README.md             # Este archivo
```

## 📝 Licencia

Este proyecto está bajo la licencia MIT.


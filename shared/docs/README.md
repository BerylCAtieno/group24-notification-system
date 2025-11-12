# API Documentation

This directory contains OpenAPI/Swagger documentation for the notification system APIs.

## Files

- **`orchestrator-api.yaml`** - Orchestrator Service API documentation
- **`api-gateway-api.yaml`** - API Gateway API documentation (proxies to orchestrator)

## Viewing the Documentation

### Option 1: Swagger UI (Online)

1. Go to [Swagger Editor](https://editor.swagger.io/)
2. Copy the contents of either YAML file
3. Paste into the editor to view interactive documentation

### Option 2: Swagger UI (Local)

#### Using Docker

```bash
# Run Swagger UI container
docker run -p 8080:8080 -e SWAGGER_JSON=/docs/orchestrator-api.yaml -v $(pwd)/shared/docs:/docs swaggerapi/swagger-ui

# Access at http://localhost:8080
```

#### Using Node.js

```bash
# Install swagger-ui-serve globally
npm install -g swagger-ui-serve

# Serve orchestrator API docs
swagger-ui-serve shared/docs/orchestrator-api.yaml

# Serve API Gateway docs
swagger-ui-serve shared/docs/api-gateway-api.yaml
```

### Option 3: Redoc (Alternative UI)

```bash
# Install redoc-cli
npm install -g redoc-cli

# Generate HTML documentation
redoc-cli bundle shared/docs/orchestrator-api.yaml -o orchestrator-api.html
redoc-cli bundle shared/docs/api-gateway-api.yaml -o api-gateway-api.html

# Or serve directly
redoc-cli serve shared/docs/orchestrator-api.yaml
```

### Option 4: VS Code Extension

1. Install the "OpenAPI (Swagger) Editor" extension in VS Code
2. Open either YAML file
3. Use the preview feature to view formatted documentation

## Quick Reference

### Orchestrator Service

**Base URL**: `http://localhost:8080` (or `http://orchestrator:8080` in Docker)

**Endpoints**:
- `POST /api/v1/notifications` - Create notification
- `POST /api/v1/email/status` - Update email notification status
- `POST /api/v1/push/status` - Update push notification status
- `POST /api/v1/users` - Create user
- `GET /health` - Health check
- `GET /health/live` - Liveness probe
- `GET /health/ready` - Readiness probe

**Authentication**: API Key via `X-API-Key` header or `api_key` query parameter

### API Gateway

**Base URL**: `http://localhost:8080` (or `http://api-gateway:80` in Docker)

**Endpoints**: Same as Orchestrator (all requests are proxied)

**Features**:
- Rate limiting: 100 req/s for API, 10 req/s for health
- CORS enabled
- Request ID propagation
- Security headers

## Example Requests

### Create Email Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -H "X-Request-ID: req-123" \
  -d '{
    "request_id": "req-123e4567-e89b-12d3-a456-426614174000",
    "notification_type": "email",
    "user_id": "usr_7x9k2p",
    "template_code": "welcome_email",
    "variables": {
      "name": "John Doe",
      "link": "https://example.com/verify"
    },
    "priority": 2
  }'
```

### Update Notification Status

```bash
curl -X POST http://localhost:8080/api/v1/email/status \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "notification_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "delivered",
    "timestamp": "2025-11-12T16:30:05Z"
  }'
```

## Response Format

All endpoints return a standardized response format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "error": null
}
```

Error responses:

```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description"
}
```

## Field Naming Convention

All field names use **snake_case** (e.g., `request_id`, `notification_type`, `user_id`).

## Idempotency

The `POST /api/v1/notifications` endpoint supports idempotency via the `request_id` field. If the same `request_id` is sent twice, the cached response will be returned with `X-Idempotent-Replay: true` header.

## Rate Limiting

**API Gateway**:
- API endpoints: 100 requests/second per IP
- Health endpoints: 10 requests/second per IP

Rate limit exceeded returns `429 Too Many Requests`.

## Status Codes

- `200 OK` - Success (cached response for notifications)
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request payload
- `401 Unauthorized` - Missing or invalid API key
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `502 Bad Gateway` - Backend service unavailable
- `503 Service Unavailable` - Service temporarily unavailable
- `504 Gateway Timeout` - Request timeout

## Integration with Swagger Tools

### Generate Client SDKs

```bash
# Using OpenAPI Generator
openapi-generator-cli generate \
  -i shared/docs/orchestrator-api.yaml \
  -g go \
  -o generated/go-client

# Other languages
openapi-generator-cli generate \
  -i shared/docs/orchestrator-api.yaml \
  -g python \
  -o generated/python-client
```

### Validate Documentation

```bash
# Using swagger-codegen
swagger-codegen validate -i shared/docs/orchestrator-api.yaml

# Using spectral (OpenAPI linter)
npm install -g @stoplight/spectral-cli
spectral lint shared/docs/orchestrator-api.yaml
```

## Notes

- The API Gateway proxies all requests to the Orchestrator Service
- Both services use the same API structure
- Authentication is required for all `/api/v1/*` endpoints
- Health endpoints do not require authentication
- All timestamps are in ISO 8601 format (e.g., `2025-11-12T16:30:00Z`)


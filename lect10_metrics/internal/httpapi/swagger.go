package httpapi

const swaggerDocJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "System Monitoring Demo API",
    "description": "Small demo API with Prometheus metrics.",
    "version": "1.0.0"
  },
  "host": "localhost:8081",
  "basePath": "/",
  "schemes": ["http"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "ok",
            "schema": {
              "$ref": "#/definitions/HealthResponse"
            }
          }
        }
      }
    },
    "/api/v1/ping": {
      "get": {
        "summary": "Ping endpoint",
        "description": "Returns pong payload and updates all monitoring metrics.",
        "parameters": [
          {
            "in": "query",
            "name": "name",
            "type": "string",
            "required": false,
            "description": "Name for greeting"
          },
          {
            "in": "query",
            "name": "fail",
            "type": "boolean",
            "required": false,
            "description": "Set true to force 500 for error metrics"
          }
        ],
        "responses": {
          "200": {
            "description": "pong response",
            "schema": {
              "$ref": "#/definitions/PingResponse"
            }
          },
          "500": {
            "description": "forced error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "HealthResponse": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string",
          "example": "ok"
        }
      }
    },
    "PingResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "pong, monitoring"
        },
        "server_time": {
          "type": "string",
          "example": "2026-05-08T10:00:00Z"
        },
        "duration_ms": {
          "type": "integer",
          "format": "int64",
          "example": 42
        }
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string",
          "example": "forced failure"
        }
      }
    }
  }
}`

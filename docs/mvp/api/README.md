# API FMgo

## Descripción General
API REST para la emisión y gestión de Documentos Tributarios Electrónicos (DTE).

## Base URL
```
https://api.fmgo.cl/v1
```

## Autenticación
```http
Authorization: Bearer <token>
```

## Endpoints

### DTE

#### Emitir DTE
```http
POST /dte
Content-Type: application/json

{
  "tipo_dte": "33",
  "emisor": {
    "rut": "76123456-7",
    "razon_social": "EMPRESA SPA",
    "giro": "SERVICIOS INFORMATICOS",
    "direccion": "CALLE EJEMPLO 123",
    "comuna": "SANTIAGO"
  },
  "receptor": {
    "rut": "77654321-8",
    "razon_social": "CLIENTE LTDA",
    "giro": "COMERCIO",
    "direccion": "AV CLIENTE 456",
    "comuna": "PROVIDENCIA"
  },
  "detalles": [
    {
      "cantidad": 1,
      "descripcion": "Servicio Profesional",
      "precio_unitario": 100000,
      "monto_total": 100000
    }
  ],
  "totales": {
    "monto_neto": 100000,
    "tasa_iva": 19,
    "iva": 19000,
    "total": 119000
  }
}

Response 200:
{
  "id": "dte_123456789",
  "folio": 1234,
  "estado": "PENDIENTE",
  "timestamp": "2024-03-15T10:30:00Z"
}
```

#### Consultar Estado
```http
GET /dte/{id}/estado

Response 200:
{
  "id": "dte_123456789",
  "estado": "ACEPTADO",
  "track_id": "12345678",
  "timestamp": "2024-03-15T10:35:00Z",
  "detalles": {
    "estado_sii": "ACEPTADO",
    "fecha_proceso": "2024-03-15T10:34:00Z",
    "errores": []
  }
}
```

#### Reenviar DTE
```http
POST /dte/{id}/reenviar

Response 200:
{
  "id": "dte_123456789",
  "estado": "PENDIENTE",
  "timestamp": "2024-03-15T11:30:00Z"
}
```

#### Listar DTEs
```http
GET /dte?estado=ACEPTADO&fecha_inicio=2024-03-01&fecha_fin=2024-03-15

Response 200:
{
  "total": 100,
  "pagina": 1,
  "por_pagina": 20,
  "dtes": [
    {
      "id": "dte_123456789",
      "folio": 1234,
      "estado": "ACEPTADO",
      "timestamp": "2024-03-15T10:30:00Z"
    },
    ...
  ]
}
```

### Administración

#### Estadísticas
```http
GET /admin/stats

Response 200:
{
  "dtes_emitidos": 1000,
  "tasa_aceptacion": 99.5,
  "tiempo_promedio": 150,
  "errores": 5
}
```

#### Estado del Sistema
```http
GET /admin/health

Response 200:
{
  "status": "healthy",
  "servicios": {
    "api": "up",
    "redis": "up",
    "postgres": "up",
    "sii": "up"
  },
  "metricas": {
    "cpu": 45.5,
    "memoria": 1.2,
    "latencia": 150
  }
}
```

## Errores

### Formato
```json
{
  "error": {
    "codigo": "ERROR_CODE",
    "mensaje": "Descripción del error",
    "detalles": {
      "campo": "descripción"
    }
  }
}
```

### Códigos HTTP
- `200`: Éxito
- `400`: Error de validación
- `401`: No autorizado
- `403`: Prohibido
- `404`: No encontrado
- `429`: Demasiadas peticiones
- `500`: Error interno
- `503`: Servicio no disponible

### Códigos de Error
```yaml
validacion:
  DTE001: "RUT inválido"
  DTE002: "Monto total no coincide"
  DTE003: "CAF no disponible"

sistema:
  SYS001: "Servicio SII no disponible"
  SYS002: "Error de base de datos"
  SYS003: "Error de caché"

autenticacion:
  AUTH001: "Token inválido"
  AUTH002: "Token expirado"
  AUTH003: "Permisos insuficientes"
```

## Rate Limiting
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1521554400
```

## Paginación
```http
Link: <https://api.fmgo.cl/v1/dte?pagina=2>; rel="next",
      <https://api.fmgo.cl/v1/dte?pagina=10>; rel="last"
```

## Versionado
- `v1`: Versión actual
- `v2`: En desarrollo (beta)
- `v0`: Deprecada

## Ejemplos

### curl
```bash
# Emitir DTE
curl -X POST https://api.fmgo.cl/v1/dte \
  -H "Authorization: Bearer token123" \
  -H "Content-Type: application/json" \
  -d @dte.json

# Consultar estado
curl https://api.fmgo.cl/v1/dte/123/estado \
  -H "Authorization: Bearer token123"
```

### Python
```python
import requests

# Configuración
api_url = "https://api.fmgo.cl/v1"
headers = {
    "Authorization": "Bearer token123",
    "Content-Type": "application/json"
}

# Emitir DTE
response = requests.post(
    f"{api_url}/dte",
    headers=headers,
    json=dte_data
)

# Consultar estado
response = requests.get(
    f"{api_url}/dte/123/estado",
    headers=headers
)
```

## Webhooks

### Configuración
```http
POST /webhooks
Content-Type: application/json

{
  "url": "https://mi-empresa.com/webhook",
  "eventos": ["dte.aceptado", "dte.rechazado"],
  "secret": "mi_secret_123"
}
```

### Formato de Eventos
```json
{
  "id": "evt_123456",
  "tipo": "dte.aceptado",
  "datos": {
    "dte_id": "dte_123456789",
    "estado": "ACEPTADO",
    "timestamp": "2024-03-15T10:35:00Z"
  }
}
```

## Seguridad
- TLS 1.2+
- CORS configurado
- Rate limiting por token
- Validación de IPs
- Logs de auditoría 
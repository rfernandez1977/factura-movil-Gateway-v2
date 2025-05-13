# API Reference

## Autenticación
Todas las peticiones a la API requieren autenticación mediante un API Key.

```http
X-API-Key: tu-api-key
```

## Endpoints

### Documentos Tributarios

#### Crear Factura
```http
POST /api/v1/facturas
```

**Request Body:**
```json
{
  "rut_emisor": "76.123.456-7",
  "rut_receptor": "77.890.123-4",
  "razon_social": "Empresa Cliente SPA",
  "direccion": "Av. Principal 123",
  "comuna": "Santiago",
  "ciudad": "Santiago",
  "fecha_emision": "2024-03-15T00:00:00Z",
  "detalles": [
    {
      "codigo": "PROD001",
      "descripcion": "Producto 1",
      "cantidad": 2,
      "precio": 10000,
      "unidad_medida": "UN",
      "descuento": 0
    }
  ]
}
```

**Response:**
```json
{
  "id": "123456789",
  "track_id": "ABC123",
  "estado": "ENVIADO",
  "fecha_creacion": "2024-03-15T10:30:00Z"
}
```

#### Crear Boleta
```http
POST /api/v1/boletas
```

**Request Body:**
```json
{
  "rut_emisor": "76.123.456-7",
  "rut_receptor": "77.890.123-4",
  "razon_social": "Cliente Final",
  "direccion": "Calle 123",
  "comuna": "Santiago",
  "ciudad": "Santiago",
  "fecha_emision": "2024-03-15T00:00:00Z",
  "detalles": [
    {
      "codigo": "PROD002",
      "descripcion": "Producto 2",
      "cantidad": 1,
      "precio": 5000,
      "unidad_medida": "UN"
    }
  ]
}
```

### Gestión de Clientes

#### Crear Cliente
```http
POST /api/v1/clientes
```

**Request Body:**
```json
{
  "rut": "77.890.123-4",
  "razon_social": "Empresa Cliente SPA",
  "direccion": "Av. Principal 123",
  "comuna": "Santiago",
  "ciudad": "Santiago",
  "email": "contacto@empresa.cl",
  "telefono": "+56 2 2123 4567"
}
```

### Gestión de Productos

#### Crear Producto
```http
POST /api/v1/productos
```

**Request Body:**
```json
{
  "codigo": "PROD001",
  "nombre": "Producto 1",
  "descripcion": "Descripción del producto",
  "precio": 10000,
  "unidad_medida": "UN",
  "stock": 100,
  "categoria": "Categoría 1"
}
```

### Consultas

#### Consultar Estado de Documento
```http
GET /api/v1/documentos/{id}/estado
```

**Response:**
```json
{
  "id": "123456789",
  "tipo": "FACTURA",
  "estado": "ACEPTADO",
  "fecha_emision": "2024-03-15T00:00:00Z",
  "fecha_aceptacion": "2024-03-15T01:00:00Z",
  "track_id": "ABC123"
}
```

#### Descargar PDF
```http
GET /api/v1/documentos/{id}/pdf
```

**Response:**
- Content-Type: application/pdf
- Body: Archivo PDF del documento

### Métricas

#### Obtener Métricas
```http
GET /api/v1/metrics
```

**Response:**
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="POST",endpoint="/facturas"} 100
http_requests_total{method="GET",endpoint="/documentos"} 200
```

## Códigos de Error

| Código | Descripción |
|--------|-------------|
| 400 | Bad Request - La petición es inválida |
| 401 | Unauthorized - API Key inválida o faltante |
| 403 | Forbidden - No tiene permisos para realizar la acción |
| 404 | Not Found - El recurso no existe |
| 429 | Too Many Requests - Se ha excedido el límite de peticiones |
| 500 | Internal Server Error - Error interno del servidor |

## Ejemplos de Uso

### Python
```python
import requests

headers = {
    'X-API-Key': 'tu-api-key',
    'Content-Type': 'application/json'
}

# Crear factura
factura = {
    "rut_emisor": "76.123.456-7",
    "rut_receptor": "77.890.123-4",
    "razon_social": "Empresa Cliente SPA",
    "fecha_emision": "2024-03-15T00:00:00Z",
    "detalles": [
        {
            "codigo": "PROD001",
            "descripcion": "Producto 1",
            "cantidad": 2,
            "precio": 10000
        }
    ]
}

response = requests.post(
    'https://api.tudominio.com/api/v1/facturas',
    headers=headers,
    json=factura
)
```

### JavaScript
```javascript
const axios = require('axios');

const api = axios.create({
    baseURL: 'https://api.tudominio.com/api/v1',
    headers: {
        'X-API-Key': 'tu-api-key'
    }
});

// Crear factura
const factura = {
    rut_emisor: "76.123.456-7",
    rut_receptor: "77.890.123-4",
    razon_social: "Empresa Cliente SPA",
    fecha_emision: "2024-03-15T00:00:00Z",
    detalles: [
        {
            codigo: "PROD001",
            descripcion: "Producto 1",
            cantidad: 2,
            precio: 10000
        }
    ]
};

api.post('/facturas', factura)
    .then(response => console.log(response.data))
    .catch(error => console.error(error));
``` 
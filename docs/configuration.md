# Documentación de Configuración

## Variables de Entorno

### Configuración General
```env
# Ambiente (CERTIFICACION/PRODUCCION)
SII_AMBIENTE=CERTIFICACION

# Puerto del servidor
PORT=8080

# Modo de simulación
SIMULACION=true
```

### Configuración SII
```env
# URLs del SII
SII_BASE_URL=https://palena.sii.cl
SII_SEMILLA_URL=https://palena.sii.cl/DTEWS/CrSeed.jws
SII_TOKEN_URL=https://palena.sii.cl/DTEWS/GetTokenFromSeed.jws
SII_ENVIO_URL=https://palena.sii.cl/cgi_dte/UPL/DTEUpload
SII_ESTADO_URL=https://palena.sii.cl/DTEWS/QueryEstDte.jws

# Certificados
CERT_PATH=./certs/cert.pem
KEY_PATH=./certs/key.pem
CERT_PASSWORD=tu_contraseña
```

### Configuración Base de Datos
```env
# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=fmgo
MONGODB_USER=usuario
MONGODB_PASSWORD=contraseña
```

### Configuración de Logs
```env
# Nivel de log (DEBUG/INFO/WARN/ERROR)
LOG_LEVEL=INFO

# Ruta de logs
LOG_PATH=./logs

# Rotación de logs
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=5
LOG_MAX_AGE=30
```

## Archivos de Configuración

### config.json
```json
{
  "sii": {
    "ambiente": "CERTIFICACION",
    "base_url": "https://palena.sii.cl",
    "cert_path": "./certs/cert.pem",
    "key_path": "./certs/key.pem"
  },
  "database": {
    "uri": "mongodb://localhost:27017",
    "database": "fmgo",
    "user": "usuario",
    "password": "contraseña"
  },
  "logging": {
    "level": "INFO",
    "path": "./logs",
    "max_size": 100,
    "max_backups": 5,
    "max_age": 30
  }
}
```

### CAF Config
```yaml
# config/caf_config.yaml
ruta_caf: ./cafs
tipos_documento:
  - id: 33
    nombre: Factura Electrónica
  - id: 39
    nombre: Boleta Electrónica
  - id: 56
    nombre: Nota Débito Electrónica
  - id: 61
    nombre: Nota Crédito Electrónica
```

## Configuración de Servicios

### SIIClient
```go
type SIIClientConfig struct {
    BaseURL      string
    CertPath     string
    KeyPath      string
    CertPassword string
    Ambiente     string
    Timeout      time.Duration
}
```

### CAFManager
```go
type CAFConfig struct {
    RutaCAF string
    TiposDocumento []TipoDocumento
}
```

### DTEGenerator
```go
type DTEGeneratorConfig struct {
    RutaTemplates string
    RutaPDFs      string
    RutaXMLs      string
}
```

## Configuración de Seguridad

### Certificados Digitales
- Formato: PKCS#12 (.p12)
- Algoritmo: RSA
- Tamaño de clave: 2048 bits
- Hash: SHA-256

### Autenticación
- JWT para API
- Certificados para SII
- API Keys para servicios externos

## Configuración de Monitoreo

### Prometheus
```yaml
# config/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'fmgo'
    static_configs:
      - targets: ['localhost:8080']
```

### Grafana
```json
// config/grafana-dashboard.json
{
  "dashboard": {
    "id": null,
    "title": "FMgo Dashboard",
    "tags": ["fmgo", "sii"],
    "timezone": "browser",
    "panels": [
      // ... configuración de paneles
    ]
  }
}
```

## Configuración de Pruebas

### Mock Server
```yaml
# config/mock-server.yaml
port: 8081
endpoints:
  - path: /DTEWS/CrSeed.jws
    response: ./testdata/semilla_response.xml
  - path: /DTEWS/GetTokenFromSeed.jws
    response: ./testdata/token_response.xml
  - path: /cgi_dte/UPL/DTEUpload
    response: ./testdata/envio_response.xml
```

### Test Data
```json
// config/test-data.json
{
  "empresa": {
    "rut": "76.123.456-7",
    "razon_social": "Empresa de Prueba",
    "giro": "Servicios de Prueba"
  },
  "documentos": [
    {
      "tipo": "FACTURA",
      "folio": 1,
      "monto": 1000
    }
  ]
}
```

## Configuración de Despliegue

### Docker
```dockerfile
# Dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
```

### Docker Compose
```yaml
# docker-compose.yml
version: '3'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SII_AMBIENTE=CERTIFICACION
    volumes:
      - ./certs:/app/certs
      - ./cafs:/app/cafs
      - ./logs:/app/logs
``` 